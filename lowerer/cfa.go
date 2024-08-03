package lowerer

import (
	"fmt"
	"os"
	"slices"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (l *lowerer) cfa(statements []ir.Statement, location text.Location) []ir.Statement {
	g := graph{
		diagnostics: l.diagnostics,
		blocks:      []*basicBlock{},
		currentBlock: &basicBlock{
			statements:  []ir.Statement{},
			entries:     []*connection{},
			exits:       []*connection{},
		},
	}
	g.analyse(statements, location)
	l.diagnostics = g.diagnostics
	result := []ir.Statement{}

	for i, block := range g.blocks {
		if !block.unreachable {
			if len(block.statements) == 0 && len(block.exits) == 1 {
				continue
			}

			if len(block.entries) != 0 {
				result = append(result, &ir.Label{Name: fmt.Sprintf("block%d", i)})
			}

			result = append(result, block.statements...)
			for _, exit := range block.exits {
				if exit.condition == nil {
					result = append(result, &ir.Goto{Label: fmt.Sprintf("block%d", exit.to)})
				} else {
					result = append(result, &ir.GotoIf{Condition: exit.condition, Label: fmt.Sprintf("block%d", exit.to)})
				}
			}
		}
	}

	return result
}

type basicBlock struct {
	label       string
	statements  []ir.Statement
	entries     []*connection
	exits       []*connection
	unreachable bool
}

type connection struct {
	from, to  int
	condition ir.Expression
}

type graph struct {
	diagnostics  diagnostics.Manager
	blocks       []*basicBlock
	currentBlock *basicBlock
}

func (g *graph) Print(node *printer.Node) {
	node.Text("%sGRAPH", node.Colour(colour.NodeName))
	for i, block := range g.blocks {
		node.FakeNode(
			"%sBLOCK %s%d",
			func(n *printer.Node) {
				n.TextIf(len(block.label) != 0, " %s%s", n.Colour(colour.Name), block.label)

				for _, entry := range block.entries {
					n.Text(
						" %s<- %s%d",
						n.Colour(colour.Symbol),
						n.Colour(colour.Literal),
						entry.from,
					)
				}

				printer.Nodes(n, block.statements)
				for _, exit := range block.exits {
					n.Text(
						" %s-> %s%d",
						n.Colour(colour.Symbol),
						n.Colour(colour.Literal),
						exit.to,
					)
				}
			},
			node.Colour(colour.NodeName),
			node.Colour(colour.Literal),
			i,
		)
	}
}

func (g *graph) analyse(statements []ir.Statement, location text.Location) {
	g.separateBlocks(statements)
	g.makeConnections()

	p := printer.New(os.Stdout, true)
	p.Node(g)
	p.Print()
	fmt.Println()

	g.removeUnreachable()
	if !g.checkPaths() {
		g.diagnostics = append(g.diagnostics, *diagnostics.NotAllPathsReturn(location))
	}
}

func (g *graph) separateBlocks(statements []ir.Statement) {
	for _, statement := range statements {
		switch stmt := statement.(type) {
		case *ir.Label:
			g.beginBlock(stmt.Name)
		case *ir.Goto:
			g.statement(stmt)
			g.endBlock()
		case *ir.GotoIf:
			g.statement(stmt)
			g.endBlock()
		case *ir.VariableDeclaration:
			g.statement(stmt)
		case *ir.ReturnStatement:
			g.statement(stmt)
			g.endBlock()
		case ir.Expression:
			g.statement(stmt)
		default:
			panic(fmt.Sprintf("Unexpected lowered statement %T", statement))
		}
	}
	g.endBlock()
}

func (g *graph) endBlock() {
	g.blocks = append(g.blocks, g.currentBlock)
	g.currentBlock = &basicBlock{
		label:      "",
		statements: []ir.Statement{},
		entries:    []*connection{},
		exits:      []*connection{},
	}
}

func (g *graph) beginBlock(label string) {
	g.endBlock()
	g.currentBlock.label = label
}

func (g *graph) statement(statement ir.Statement) {
	g.currentBlock.statements = append(g.currentBlock.statements, statement)
}

func (g *graph) makeConnections() {
	if len(g.blocks) == 0 {
		return
	}
	for i, block := range g.blocks[:len(g.blocks)-1] {
		var last ir.Statement
		if len(block.statements) != 0 {
			last = block.statements[len(block.statements)-1]
		}

		switch stmt := last.(type) {
		case *ir.Goto:
			g.connection(i, g.blockWithLabel(stmt.Label))
			block.statements = block.statements[:len(block.statements)-1]
		case *ir.GotoIf:
			g.conditionalConnection(i, g.blockWithLabel(stmt.Label), stmt.Condition)
			block.statements = block.statements[:len(block.statements)-1]
		case *ir.ReturnStatement:
		default:
			g.connection(i, i+1)
		}
	}
}

func (g *graph) connection(from, to int) {
	g.doConnection(from, to, nil)
}

func (g *graph) conditionalConnection(from, to int, condition ir.Expression) {
	g.doConnection(from, to, condition)
	g.doConnection(from, from+1, negate(condition))
}

func (g *graph) doConnection(from, to int, condition ir.Expression) {
	if condition != nil {
		if boolean, ok := condition.ConstValue().(values.BoolValue); ok {
			if boolean.Value {
				g.doConnection(from, to, nil)
			}
			return
		}
	}

	conn := &connection{
		from:      from,
		to:        to,
		condition: condition,
	}
	g.blocks[to].entries = append(g.blocks[to].entries, conn)
	g.blocks[from].exits = append(g.blocks[from].exits, conn)
}

func (g *graph) blockWithLabel(label string) int {
	for i, block := range g.blocks {
		if block.label == label {
			return i
		}
	}
	panic("Block should exist")
}

func (g *graph) removeUnreachable() {
	for i, block := range g.blocks {
		// The first block always has an entry point
		if i == 0 {
			continue
		}
		if len(block.entries) == 0 {
			// Don't re-traverse already unreachable blocks
			if block.unreachable {
				continue
			}

			// We can't actually remove the block, since that messes up
			// all the indices, so we just mark the block as unreachable.
			g.blocks[i].unreachable = true

			for _, exit := range block.exits {
				g.blocks[exit.to].entries = slices.DeleteFunc(
					g.blocks[exit.to].entries,
					func(entry *connection) bool { return entry.from == i },
				)
			}
			g.blocks[i].exits = block.exits[:0]
			// We have to re-check all blocks because a block whose only
			// entry point is another unreachable block is also unreachable
			g.removeUnreachable()
			return
		}
	}
}

func (g *graph) checkPaths() bool {
	if len(g.blocks) == 0 {
		return false
	}

	for _, block := range g.blocks {
		if len(block.entries) == 0 {
			continue
		}
		if len(block.exits) == 0 {
			switch block.statements[len(block.statements)-1].(type) {
			case *ir.ReturnStatement:
			default:
				return false
			}
		}
	}
	return true
}

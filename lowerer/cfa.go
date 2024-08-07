package lowerer

import (
	"fmt"
	"slices"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (l *lowerer) cfa(
	statements []ir.Statement,
	location *text.Location,
	shouldReturn bool,
) []ir.Statement {
	g := graph{
		diagnostics: l.diagnostics,
		blocks:      []*basicBlock{},
		currentBlock: &basicBlock{
			statements: []ir.Statement{},
			entries:    []*connection{},
			exit:       nil,
			isStart:    true,
		},
		connections: []*connection{},
	}
	g.analyse(statements, location, shouldReturn)
	l.diagnostics = g.diagnostics
	result := []ir.Statement{}

	for i, block := range g.blocks {
		if !block.unreachable {
			if len(block.entries) != 0 || block.isStart {
				result = append(result, &ir.Label{Name: fmt.Sprintf("block%d", i)})
			}

			result = append(result, block.statements...)
			if exit := block.exit; exit != nil {
				if exit.condition == nil {
					result = append(result, &ir.Goto{Label: fmt.Sprintf("block%d", exit.to)})
				} else {
					result = append(result, &ir.Branch{
						Condition: exit.condition,
						IfLabel:   fmt.Sprintf("block%d", exit.to),
						ElseLabel: fmt.Sprintf("block%d", exit.elseTo),
					})
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
	exit        *connection
	unreachable bool
	isStart     bool
}

type connection struct {
	from, to  int
	condition ir.Expression
	elseTo    int
}

type graph struct {
	diagnostics  diagnostics.Manager
	blocks       []*basicBlock
	currentBlock *basicBlock
	connections  []*connection
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
				if exit := block.exit; exit != nil {
					n.Text(
						" %s-> %s%d",
						n.Colour(colour.Symbol),
						n.Colour(colour.Literal),
						exit.to,
					)
					if exit.condition != nil {
						n.Text(
							" %s-> %s%d",
							n.Colour(colour.Symbol),
							n.Colour(colour.Literal),
							exit.elseTo,
						)
					}
				}
			},
			node.Colour(colour.NodeName),
			node.Colour(colour.Literal),
			i,
		)
	}
}

func (g *graph) analyse(
	statements []ir.Statement,
	location *text.Location,
	shouldReturn bool,
) {
	g.separateBlocks(statements)
	g.makeConnections()

	// p := printer.New(os.Stdout, true)
	// p.Node(g)
	// p.Print()
	// fmt.Println()

	g.removeUnreachable()
	g.remapIds()
	if shouldReturn && !g.checkPaths() {
		g.diagnostics = append(g.diagnostics, *diagnostics.NotAllPathsReturn(*location))
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
		exit:       nil,
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
	g.doConnection(from, to, nil, 0)
}

func (g *graph) conditionalConnection(from, to int, condition ir.Expression) {
	g.doConnection(from, to, condition, from+1)
}

func (g *graph) doConnection(from, to int, condition ir.Expression, elseTo int) {
	if condition != nil {
		if boolean, ok := condition.ConstValue().(values.BoolValue); ok {
			if boolean.Value {
				g.doConnection(from, to, nil, 0)
			} else {
				g.doConnection(from, elseTo, nil, 0)
			}
			return
		}
	}

	conn := &connection{
		from:      from,
		to:        to,
		condition: condition,
		elseTo:    elseTo,
	}
	g.blocks[to].entries = append(g.blocks[to].entries, conn)
	if condition != nil {
		g.blocks[elseTo].entries = append(g.blocks[elseTo].entries, conn)
	}
	g.blocks[from].exit = conn
	g.connections = append(g.connections, conn)
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
		// Don't re-traverse already unreachable blocks
		if block.unreachable {
			continue
		}

		// The first block always has an entry point
		if len(block.entries) == 0 && !block.isStart {
			// We can't actually remove the block, since that messes up
			// all the indices, so we just mark the block as unreachable.
			block.unreachable = true

			if exit := block.exit; exit != nil {
				exitBlock := g.blocks[exit.to]
				exitBlock.entries = slices.DeleteFunc(
					exitBlock.entries,
					func(entry *connection) bool { return entry.from == i },
				)

				if exit.condition != nil {
					exitBlock := g.blocks[exit.elseTo]
					exitBlock.entries = slices.DeleteFunc(
						exitBlock.entries,
						func(entry *connection) bool { return entry.from == i },
					)
				}
			}
			block.exit = nil
			// We have to re-check all blocks because a block whose only
			// entry point is another unreachable block is also unreachable
			g.removeUnreachable()
			return
		}

		if len(block.statements) == 0 && block.exit != nil && block.exit.condition == nil {
			exit := g.blocks[block.exit.to]
			exit.entries = slices.DeleteFunc(
				exit.entries,
				func(c *connection) bool { return c.from == i },
			)

			if block.exit.condition != nil {
				exit := g.blocks[block.exit.elseTo]
				exit.entries = slices.DeleteFunc(
					exit.entries,
					func(c *connection) bool { return c.from == i },
				)
			}

			for _, entry := range block.entries {
				entry.to = block.exit.to
				entry.elseTo = block.exit.elseTo
				exit.entries = append(exit.entries, entry)
			}
			if block.isStart {
				exit.isStart = true
			}

			block.unreachable = true
			g.removeUnreachable()
			return
		}

		if block.exit != nil && block.exit.condition == nil {
			nextBlock := g.blocks[block.exit.to]

			if len(nextBlock.entries) == 1 && !nextBlock.isStart {
				block.statements = append(block.statements, nextBlock.statements...)
				block.exit = nextBlock.exit
				nextBlock.unreachable = true
				g.removeUnreachable()
				return
			}
		}
	}
}

func (g *graph) remapIds() {
	ids := make([]int, 0, len(g.blocks))
	for i := range len(g.blocks) {
		ids = append(ids, i)
	}
	for i := range len(g.blocks) {
		if !g.blocks[i].unreachable {
			continue
		}

		for j := i + 1; j < len(g.blocks); j++ {
			ids[j]--
		}
	}
	g.blocks = slices.DeleteFunc(g.blocks, func(block *basicBlock) bool {
		return block.unreachable
	})
	for _, connection := range g.connections {
		connection.from = ids[connection.from]
		connection.to = ids[connection.to]
		connection.elseTo = ids[connection.elseTo]
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
		if block.exit == nil {
			if len(block.statements) == 0 {
				return false
			}

			switch block.statements[len(block.statements)-1].(type) {
			case *ir.ReturnStatement:
			default:
				return false
			}
		}
	}
	return true
}

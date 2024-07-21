package ast

import (
	"fmt"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
)

type expression struct{}

func (expression) expressionNode() {}

type IntegerLiteral struct {
	expression
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) GetLocation() text.Location {
	return il.Token.Location
}
func (il *IntegerLiteral) Print(node *printer.Node) {
	litValue := fmt.Sprint(il.Value)

	node.
		Text(
			"%sINT_LIT %s%d",
			node.Colour(colour.NodeName),
			node.Colour(colour.Literal),
			il.Value,
		).
		TextIf(
			litValue != il.Token.Value,
			" %s",
			il.Token.Value,
		).
		Location(il)
}

type FloatLiteral struct {
	expression
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) GetLocation() text.Location {
	return fl.Token.Location
}
func (fl *FloatLiteral) Print(node *printer.Node) {
	litValue := fmt.Sprint(fl.Value)

	node.
		Text(
			"%sFLOAT_LIT %s%v",
			node.Colour(colour.NodeName),
			node.Colour(colour.Literal),
			fl.Value,
		).
		TextIf(
			litValue != fl.Token.Value,
			" %s",
			fl.Token.Value,
		).
		Location(fl)
}

type BooleanLiteral struct {
	expression
	Location text.Location
	Value    bool
}

func (bl *BooleanLiteral) GetLocation() text.Location {
	return bl.Location
}
func (bl *BooleanLiteral) Print(node *printer.Node) {
	node.
		Text(
			"%sBOOL_LIT %s%t",
			node.Colour(colour.NodeName),
			node.Colour(colour.Literal),
			bl.Value,
		).
		Location(bl)
}

type StringLiteral struct {
	expression
	Token token.Token
	Value string
}

func (sl *StringLiteral) GetLocation() text.Location {
	return sl.Token.Location
}
func (sl *StringLiteral) Print(node *printer.Node) {
	node.
		Text(
			"%sSTRING_LIT %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Literal),
			sl.Token.Value,
		).
		Location(sl)
}

type Identifier struct {
	expression
	Location text.Location
	Name     string
}

func (i *Identifier) GetLocation() text.Location {
	return i.Location
}
func (i *Identifier) Print(node *printer.Node) {
	node.
		Text(
			"%sIDENT %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			i.Name,
		).
		Location(i)
}

type BinaryExpression struct {
	expression
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b *BinaryExpression) GetLocation() text.Location {
	return b.Left.GetLocation().To(b.Right.GetLocation())
}

func (b *BinaryExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sBIN_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Symbol),
			b.Operator.Value,
		).
		Location(b).
		Node(b.Left).
		Node(b.Right)
}

type ParenthesisedExpression struct {
	expression
	Location   text.Location
	Expression Expression
}

func (p *ParenthesisedExpression) GetLocation() text.Location {
	return p.Location
}

func (p *ParenthesisedExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sPAREN_EXPR",
			node.Colour(colour.NodeName),
		).
		Location(p).
		Node(p.Expression)
}

type PrefixExpression struct {
	expression
	Location text.Location
	Operator token.Kind
	Operand  Expression
}

func (p *PrefixExpression) GetLocation() text.Location {
	return p.Location
}

func (p *PrefixExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sPREFIX_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Symbol),
			p.Operator.String(),
		).
		Location(p).
		Node(p.Operand)
}

type PostfixExpression struct {
	expression
	OperatorLocation text.Location
	Operand          Expression
	Operator         token.Kind
}

func (p *PostfixExpression) GetLocation() text.Location {
	return p.Operand.GetLocation()
}

func (p *PostfixExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sPOSTFIX_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Symbol),
			p.Operator,
		).
		Location(p).
		Node(p.Operand)
}

type PointerType struct {
	expression
	Location text.Location
	Mutable  bool
	Operand  Expression
}

func (p *PointerType) GetLocation() text.Location {
	return p.Location
}

func (p *PointerType) Print(node *printer.Node) {
	node.
		Text("%sPTR_TYPE", node.Colour(colour.NodeName)).
		TextIf(p.Mutable, " %smut", node.Colour(colour.Attribute)).
		Location(p).
		Node(p.Operand)
}

type OptionType struct {
	expression
	Location text.Location
	Operand  Expression
}

func (o *OptionType) GetLocation() text.Location {
	return o.Location
}

func (o *OptionType) Print(node *printer.Node) {
	node.
		Text("%sOPTION_TYPE", node.Colour(colour.NodeName)).
		Location(o).
		Node(o.Operand)
}

type DerefExpression struct {
	expression
	Operand Expression
}

func (d *DerefExpression) GetLocation() text.Location {
	return d.Operand.GetLocation()
}

func (d *DerefExpression) Print(node *printer.Node) {
	node.
		Text("%sDEREF_EXPR", node.Colour(colour.NodeName)).
		Location(d).
		Node(d.Operand)
}

type RefExpression struct {
	expression
	Location text.Location
	Mutable  bool
	Operand  Expression
}

func (r *RefExpression) GetLocation() text.Location {
	return r.Location
}

func (r *RefExpression) Print(node *printer.Node) {
	node.
		Text("%sREF_EXPR", node.Colour(colour.NodeName)).
		TextIf(r.Mutable, " %smut", node.Colour(colour.Attribute)).
		Location(r).
		Node(r.Operand)
}

type ListLiteral struct {
	expression
	Location text.Location
	Values   []Expression
}

func (l *ListLiteral) GetLocation() text.Location {
	return l.Location
}

func (l *ListLiteral) Print(node *printer.Node) {
	node.
		Text("%sLIST_EXPR", node.Colour(colour.NodeName)).
		Location(l)
	printer.Nodes(node, l.Values)
}

type KeyValue struct {
	Key   Expression
	Value Expression
}

func (kv KeyValue) Print(node *printer.Node) {
	node.
		Text("%sKEY_VALUE", node.Colour(colour.NodeName)).
		Node(kv.Key).
		Node(kv.Value)
}

type MapLiteral struct {
	expression
	Location  text.Location
	KeyValues []KeyValue
}

func (m *MapLiteral) GetLocation() text.Location {
	return m.Location
}

func (m *MapLiteral) Print(node *printer.Node) {
	node.
		Text("%sMAP_EXPR", node.Colour(colour.NodeName)).
		Location(m)
	printer.Nodes(node, m.KeyValues)
}

type FunctionCall struct {
	expression
	Callee    Expression
	Arguments []Expression
}

func (call *FunctionCall) GetLocation() text.Location {
	return call.Callee.GetLocation()
}

func (call *FunctionCall) Print(node *printer.Node) {
	node.
		Text("%sFUNCTION_CALL", node.Colour(colour.NodeName)).
		Location(call).
		Node(call.Callee)
	printer.Nodes(node, call.Arguments)
}

type IndexExpression struct {
	expression
	Left  Expression
	Location text.Location
	Index Expression
}

func (index *IndexExpression) GetLocation() text.Location {
	return index.Location
}

func (index *IndexExpression) Print(node *printer.Node) {
	node.
		Text("%sINDEX_EXPR", node.Colour(colour.NodeName)).
		Location(index).
		Node(index.Left).
		OptionalNode(index.Index)
}

type AssignmentExpression struct {
	expression
	Assignee Expression
	Operator token.Token
	Value    Expression
}

func (a *AssignmentExpression) GetLocation() text.Location {
	return a.Assignee.GetLocation().To(a.Value.GetLocation())
}

func (a *AssignmentExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sASSIGNMENT_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Symbol),
			a.Operator.Value,
		).
		Location(a).
		Node(a.Assignee).
		Node(a.Value)
}

type TupleExpression struct {
	expression
	Location text.Location
	Values   []Expression
}

func (t *TupleExpression) GetLocation() text.Location {
	return t.Location
}

func (t *TupleExpression) Print(node *printer.Node) {
	node.
		Text("%sTUPLE_EXPR", node.Colour(colour.NodeName)).
		Location(t)
	printer.Nodes(node, t.Values)
}

type MemberExpression struct {
	expression
	Location       text.Location
	MemberLocation text.Location
	Left           Expression
	Member         string
}

func (m *MemberExpression) GetLocation() text.Location {
	return m.Location
}

func (m *MemberExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sMEMBER_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			m.Member,
		).
		Location(m).
		Node(m.Left)
}

type StructMember struct {
	Location text.Location
	Name     *string
	Value    Expression
}

func (sm StructMember) Print(node *printer.Node) {
	node.
		Text("%sSTRUCT_MEMBER", node.Colour(colour.NodeName))

	if sm.Name != nil {
		node.Text(" %s%s", node.Colour(colour.Name), *sm.Name)
	}

	node.OptionalNode(sm.Value)
}

type InferredExpression struct {
	expression
	Location text.Location
}

func (i *InferredExpression) GetLocation() text.Location {
	return i.Location
}

func (i *InferredExpression) Print(node *printer.Node) {
	node.
		Text("%sINFERRED_EXPR", node.Colour(colour.NodeName)).
		Location(i)
}

type StructExpression struct {
	expression
	Struct  Expression
	Members []StructMember
}

func (s *StructExpression) GetLocation() text.Location {
	return s.Struct.GetLocation()
}

func (s *StructExpression) Print(node *printer.Node) {
	node.
		Text("%sSTRUCT_EXPR", node.Colour(colour.NodeName)).
		Location(s).
		Node(s.Struct)
	printer.Nodes(node, s.Members)
}

type CastExpression struct {
	expression
	Location text.Location
	Left     Expression
	Type     Expression
}

func (ce *CastExpression) GetLocation() text.Location {
	return ce.Location
}

func (ce *CastExpression) Print(node *printer.Node) {
	node.
		Text("%sCAST_EXPR", node.Colour(colour.NodeName)).
		Location(ce).
		Node(ce.Left).
		Node(ce.Type)
}

type TypeCheckExpression struct {
	expression
	Location text.Location
	Left     Expression
	Type     Expression
}

func (ce *TypeCheckExpression) GetLocation() text.Location {
	return ce.Location
}

func (tc *TypeCheckExpression) Print(node *printer.Node) {
	node.
		Text("%sTYPE_CHECK_EXPR", node.Colour(colour.NodeName)).
		Location(tc).
		Node(tc.Left).
		Node(tc.Type)
}

type RangeExpression struct {
	expression
	Location text.Location
	Start    Expression
	End      Expression
}

func (r *RangeExpression) GetLocation() text.Location {
	return r.Location
}

func (r *RangeExpression) Print(node *printer.Node) {
	node.
		Text("%sRANGE_EXPR", node.Colour(colour.NodeName)).
		Location(r).
		Node(r.Start).
		Node(r.End)
}

type FunctionExpression struct {
	expression
	Location   text.Location
	Parameters []Parameter
	ReturnType Expression
	Body       *Block
}

func (f *FunctionExpression) GetLocation() text.Location {
	return f.Location
}

func (f *FunctionExpression) Print(node *printer.Node) {
	node.
		TextIf(f.Body != nil, "%sFUNC_EXPR", node.Colour(colour.NodeName)).
		TextIf(f.Body == nil, "%sFUNC_TYPE", node.Colour(colour.NodeName)).
		Location(f)

	printer.Nodes(node, f.Parameters)

	node.
		OptionalNode(f.ReturnType).
		OptionalNode(f.Body)
}

type Block struct {
	expression
	Location   text.Location
	Statements []Statement
}

func (b *Block) Print(node *printer.Node) {
	node.
		Text("%sBLOCK", node.Colour(colour.NodeName)).
		Location(b)
	printer.Nodes(node, b.Statements)
}

func (b *Block) GetLocation() text.Location {
	return b.Location
}

type IfExpression struct {
	expression
	Location   text.Location
	Condition  Expression
	Body       *Block
	ElseBranch Statement
}

func (i *IfExpression) Print(node *printer.Node) {
	node.
		Text("%sIF_EXPR", node.Colour(colour.NodeName)).
		Location(i).
		Node(i.Condition).
		Node(i.Body)

	if i.ElseBranch != nil {
		node.FakeNode("%sELSE_BRANCH", func(n *printer.Node) {
			n.Node(i.ElseBranch)
		}, node.Colour(colour.NodeName))
	}
}

func (i *IfExpression) GetLocation() text.Location {
	return i.Location
}

type WhileLoop struct {
	expression
	Location  text.Location
	Condition Expression
	Body      *Block
}

func (wl *WhileLoop) Print(node *printer.Node) {
	node.
		Text("%sWHILE_LOOP", node.Colour(colour.NodeName)).
		Location(wl).
		Node(wl.Condition).
		Node(wl.Body)
}

func (w *WhileLoop) GetLocation() text.Location {
	return w.Location
}

type ForLoop struct {
	expression
	LLocation text.Location
	Variable  string
	Iterator  Expression
	Body      *Block
}

func (fl *ForLoop) Print(node *printer.Node) {
	node.
		Text(
			"%sFOR_LOOP %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			fl.Variable,
		).
		Location(fl).
		Node(fl.Iterator).
		Node(fl.Body)
}

func (f *ForLoop) GetLocation() text.Location {
	return f.LLocation
}

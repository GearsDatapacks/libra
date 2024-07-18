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
func (il *IntegerLiteral) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sINT_LIT %s%d",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		il.Value,
	)

	litValue := fmt.Sprint(il.Value)
	if litValue != il.Token.Value {
		context.AddInfo(" %s", il.Token.Value)
	}

	context.AddLocation(il)
}

type FloatLiteral struct {
	expression
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) GetLocation() text.Location {
	return fl.Token.Location
}
func (fl *FloatLiteral) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sFLOAT_LIT %s%v",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		fl.Value,
	)

	litValue := fmt.Sprint(fl.Value)
	if litValue != fl.Token.Value {
		context.AddInfo(" %s", fl.Token.Value)
	}

	context.AddLocation(fl)
}

type BooleanLiteral struct {
	expression
	Location text.Location
	Value    bool
}

func (bl *BooleanLiteral) GetLocation() text.Location {
	return bl.Location
}
func (bl *BooleanLiteral) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sBOOL_LIT %s%t",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		bl.Value,
	)

	context.AddLocation(bl)
}

type StringLiteral struct {
	expression
	Token token.Token
	Value string
}

func (sl *StringLiteral) GetLocation() text.Location {
	return sl.Token.Location
}
func (sl *StringLiteral) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sSTRING_LIT %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		sl.Token.Value,
	)

	context.AddLocation(sl)
}

type Identifier struct {
	expression
	Location text.Location
	Name     string
}

func (i *Identifier) GetLocation() text.Location {
	return i.Location
}
func (i *Identifier) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sIDENT %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		i.Name,
	)

	context.AddLocation(i)
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

func (b *BinaryExpression) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sBIN_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		b.Operator.Value,
	)
	context.AddLocation(b)

	context.QueueNode(b.Left)
	context.QueueNode(b.Right)
}

type ParenthesisedExpression struct {
	expression
	Location   text.Location
	Expression Expression
}

func (p *ParenthesisedExpression) GetLocation() text.Location {
	return p.Location
}

func (p *ParenthesisedExpression) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sPAREN_EXPR",
		context.Colour(colour.NodeName),
	)
	context.AddLocation(p)

	context.QueueNode(p.Expression)
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

func (p *PrefixExpression) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sPREFIX_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		p.Operator.String(),
	)
	context.AddLocation(p)

	context.QueueNode(p.Operand)
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

func (p *PostfixExpression) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sPOSTFIX_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		p.Operator,
	)
	context.AddLocation(p)

	context.QueueNode(p.Operand)
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

func (p *PointerType) Print(context *printer.Printer) {
	context.QueueInfo("%sPTR_TYPE", context.Colour(colour.NodeName))
	if p.Mutable {
		context.AddInfo(" %smut", context.Colour(colour.Attribute))
	}
	context.AddLocation(p)

	context.QueueNode(p.Operand)
}

type OptionType struct {
	expression
	Location text.Location
	Operand  Expression
}

func (o *OptionType) GetLocation() text.Location {
	return o.Location
}

func (o *OptionType) Print(context *printer.Printer) {
	context.QueueInfo("%sOPTION_TYPE", context.Colour(colour.NodeName))
	context.AddLocation(o)

	context.QueueNode(o.Operand)
}

type DerefExpression struct {
	expression
	Operand Expression
}

func (d *DerefExpression) GetLocation() text.Location {
	return d.Operand.GetLocation()
}

func (d *DerefExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sDEREF_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(d)
	context.QueueNode(d.Operand)
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

func (r *RefExpression) Print(context *printer.Printer) {

	context.QueueInfo("%sREF_EXPR", context.Colour(colour.NodeName))
	if r.Mutable {
		context.AddInfo(" %smut", context.Colour(colour.Attribute))
	}
	context.AddLocation(r)
	context.QueueNode(r.Operand)
}

type ListLiteral struct {
	expression
	Location text.Location
	Values   []Expression
}

func (l *ListLiteral) GetLocation() text.Location {
	return l.Location
}

func (l *ListLiteral) Print(context *printer.Printer) {
	context.QueueInfo("%sLIST_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(l)
	printer.QueueNodeList(context, l.Values)
}

type KeyValue struct {
	Key   Expression
	Value Expression
}

func (kv KeyValue) Print(context *printer.Printer) {
	context.QueueInfo("%sKEY_VALUE", context.Colour(colour.NodeName))
	context.QueueNode(kv.Key)
	context.QueueNode(kv.Value)
}

type MapLiteral struct {
	expression
	Location  text.Location
	KeyValues []KeyValue
}

func (m *MapLiteral) GetLocation() text.Location {
	return m.Location
}

func (m *MapLiteral) Print(context *printer.Printer) {
	context.QueueInfo("%sMAP_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(m)
	printer.QueueNodeList(context, m.KeyValues)
}

type FunctionCall struct {
	expression
	Callee    Expression
	Arguments []Expression
}

func (call *FunctionCall) GetLocation() text.Location {
	return call.Callee.GetLocation()
}

func (call *FunctionCall) Print(context *printer.Printer) {
	context.QueueInfo("%sFUNCTION_CALL", context.Colour(colour.NodeName))
	context.AddLocation(call)
	context.QueueNode(call.Callee)
	printer.QueueNodeList(context, call.Arguments)
}

type IndexExpression struct {
	expression
	Left  Expression
	Index Expression
}

func (index *IndexExpression) GetLocation() text.Location {
	return index.Left.GetLocation()
}

func (index *IndexExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sINDEX_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(index)
	context.QueueNode(index.Left)
	if index.Index != nil {
		context.QueueNode(index.Index)
	}
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

func (a *AssignmentExpression) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sASSIGNMENT_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		a.Operator.Value,
	)
	context.AddLocation(a)
	context.QueueNode(a.Assignee)
	context.QueueNode(a.Value)
}

type TupleExpression struct {
	expression
	Location text.Location
	Values   []Expression
}

func (t *TupleExpression) GetLocation() text.Location {
	return t.Location
}

func (t *TupleExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sTUPLE_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(t)
	printer.QueueNodeList(context, t.Values)
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

func (m *MemberExpression) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sMEMBER_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		m.Member,
	)
	context.AddLocation(m)
	context.QueueNode(m.Left)
}

type StructMember struct {
	Location text.Location
	Name     *string
	Value    Expression
}

func (sm StructMember) Print(context *printer.Printer) {
	context.QueueInfo("%sSTRUCT_MEMBER", context.Colour(colour.NodeName))

	if sm.Name != nil {
		context.AddInfo(" %s%s", context.Colour(colour.Name), *sm.Name)
	}

	if sm.Value != nil {
		context.QueueNode(sm.Value)
	}
}

type InferredExpression struct {
	expression
	Location text.Location
}

func (i *InferredExpression) GetLocation() text.Location {
	return i.Location
}

func (i *InferredExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sINFERRED_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(i)
}

type StructExpression struct {
	expression
	Struct  Expression
	Members []StructMember
}

func (s *StructExpression) GetLocation() text.Location {
	return s.Struct.GetLocation()
}

func (s *StructExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sSTRUCT_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(s)
	context.QueueNode(s.Struct)
	printer.QueueNodeList(context, s.Members)
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

func (ce *CastExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sCAST_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(ce)
	context.QueueNode(ce.Left)
	context.QueueNode(ce.Type)
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

func (tc *TypeCheckExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sTYPE_CHECK_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(tc)
	context.QueueNode(tc.Left)
	context.QueueNode(tc.Type)
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

func (r *RangeExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sRANGE_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(r)
	context.QueueNode(r.Start)
	context.QueueNode(r.End)
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

func (f *FunctionExpression) Print(context *printer.Printer) {
	if f.Body != nil {
		context.QueueInfo("%sFUNC_EXPR", context.Colour(colour.NodeName))
	} else {
		context.QueueInfo("%sFUNC_TYPE", context.Colour(colour.NodeName))
	}
	context.AddLocation(f)

	printer.QueueNodeList(context, f.Parameters)

	if f.ReturnType != nil {
		context.QueueNode(f.ReturnType)
	}

	if f.Body != nil {
		context.QueueNode(f.Body)
	}
}

type Block struct {
	expression
	Location   text.Location
	Statements []Statement
}

func (b *Block) Print(context *printer.Printer) {
	context.QueueInfo("%sBLOCK", context.Colour(colour.NodeName))
	context.AddLocation(b)
	printer.QueueNodeList(context, b.Statements)
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

func (i *IfExpression) Print(context *printer.Printer) {
	context.QueueInfo("%sIF_EXPR", context.Colour(colour.NodeName))
	context.AddLocation(i)
	context.QueueNode(i.Condition)
	context.QueueNode(i.Body)

	if i.ElseBranch != nil {
		context.Nest()
		context.AddInfo("%sELSE_BRANCH", context.Colour(colour.NodeName))
		context.QueueNode(i.ElseBranch)
		context.UnNest()
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

func (wl *WhileLoop) Print(context *printer.Printer) {
	context.QueueInfo("%sWHILE_LOOP", context.Colour(colour.NodeName))
	context.AddLocation(wl)
	context.QueueNode(wl.Condition)
	context.QueueNode(wl.Body)
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

func (fl *ForLoop) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sFOR_LOOP %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		fl.Variable,
	)
	context.AddLocation(fl)
	context.QueueNode(fl.Iterator)
	context.QueueNode(fl.Body)
}

func (f *ForLoop) GetLocation() text.Location {
	return f.LLocation
}

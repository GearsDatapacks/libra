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
func (il *IntegerLiteral) String(context printer.Context) {
	context.Write(
		"%sINT_LIT %s%d",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		il.Value,
	)

	litValue := fmt.Sprint(il.Value)
	if litValue != il.Token.Value {
		context.Write(" %s", il.Token.Value)
	}
}

type FloatLiteral struct {
	expression
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) GetLocation() text.Location {
	return fl.Token.Location
}
func (fl *FloatLiteral) String(context printer.Context) {
	context.Write(
		"%sFLOAT_LIT %s%v",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		fl.Value,
	)

	litValue := fmt.Sprint(fl.Value)
	if litValue != fl.Token.Value {
		context.Write(" %s", fl.Token.Value)
	}
}

type BooleanLiteral struct {
	expression
	Location text.Location
	Value    bool
}

func (bl *BooleanLiteral) GetLocation() text.Location {
	return bl.Location
}
func (bl *BooleanLiteral) String(context printer.Context) {
	context.Write(
		"%sBOOL_LIT %s%t",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		bl.Value,
	)
}

type StringLiteral struct {
	expression
	Token token.Token
	Value string
}

func (sl *StringLiteral) GetLocation() text.Location {
	return sl.Token.Location
}
func (sl *StringLiteral) String(context printer.Context) {
	context.Write(
		"%sSTRING_LIT %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Literal),
		sl.Token.Value,
	)
}

type Identifier struct {
	expression
	Location text.Location
	Name     string
}

func (i *Identifier) GetLocation() text.Location {
	return i.Location
}
func (i *Identifier) String(context printer.Context) {
	context.Write(
		"%sIDENT %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		i.Name,
	)
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

func (b *BinaryExpression) String(context printer.Context) {
	context.Write(
		"%sBIN_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		b.Operator.Value,
	)
	context.WriteNode(b.Left)
	context.WriteNode(b.Right)
}

type ParenthesisedExpression struct {
	expression
	Location   text.Location
	Expression Expression
}

func (p *ParenthesisedExpression) GetLocation() text.Location {
	return p.Location
}

func (p *ParenthesisedExpression) String(context printer.Context) {
	context.Write(
		"%sPAREN_EXPR",
		context.Colour(colour.NodeName),
	)
	context.WriteNode(p.Expression)
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

func (p *PrefixExpression) String(context printer.Context) {
	context.Write(
		"%sPREFIX_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		p.Operator.String(),
	)
	context.WriteNode(p.Operand)
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

func (p *PostfixExpression) String(context printer.Context) {
	context.Write(
		"%sPOSTFIX_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		p.Operator,
	)
	context.WriteNode(p.Operand)
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

func (p *PointerType) String(context printer.Context) {
	context.Write("%sPTR_TYPE", context.Colour(colour.NodeName))
	if p.Mutable {
		context.Write(" %smut", context.Colour(colour.Attribute))
	}

	context.WriteNode(p.Operand)
}

type OptionType struct {
	expression
	Location text.Location
	Operand  Expression
}

func (p *OptionType) GetLocation() text.Location {
	return p.Location
}

func (p *OptionType) String(context printer.Context) {
	context.Write("%sOPTION_TYPE", context.Colour(colour.NodeName))
	context.WriteNode(p.Operand)
}

type DerefExpression struct {
	expression
	Operand Expression
}

func (d *DerefExpression) GetLocation() text.Location {
	return d.Operand.GetLocation()
}

func (d *DerefExpression) String(context printer.Context) {
	context.Write("%sDEREF_EXPR", context.Colour(colour.NodeName))
	context.WriteNode(d.Operand)
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

func (r *RefExpression) String(context printer.Context) {

	context.Write("%sREF_EXPR", context.Colour(colour.NodeName))
	if r.Mutable {
		context.Write(" %smut", context.Colour(colour.Attribute))
	}
	context.WriteNode(r.Operand)
}

type ListLiteral struct {
	expression
	Location text.Location
	Values   []Expression
}

func (l *ListLiteral) GetLocation() text.Location {
	return l.Location
}

func (l *ListLiteral) String(context printer.Context) {
	context.Write("%sLIST_EXPR", context.Colour(colour.NodeName))
	if len(l.Values) != 0 {
		printer.WriteNodeList(context.WithNest(), l.Values)
	}
}

type KeyValue struct {
	Key   Expression
	Value Expression
}

func (kv KeyValue) String(context printer.Context) {
	context.Write("%sKEY_VALUE", context.Colour(colour.NodeName))
	context.WriteNode(kv.Key)
	context.WriteNode(kv.Value)
}

type MapLiteral struct {
	expression
	Location  text.Location
	KeyValues []KeyValue
}

func (m *MapLiteral) GetLocation() text.Location {
	return m.Location
}

func (m *MapLiteral) String(context printer.Context) {
	context.Write("%sMAP_EXPR", context.Colour(colour.NodeName))
	if len(m.KeyValues) != 0 {
		printer.WriteNodeList(context.WithNest(), m.KeyValues)
	}
}

type FunctionCall struct {
	expression
	Callee    Expression
	Arguments []Expression
}

func (call *FunctionCall) GetLocation() text.Location {
	return call.Callee.GetLocation()
}

func (call *FunctionCall) String(context printer.Context) {
	context.Write("%sFUNCTION_CALL", context.Colour(colour.NodeName))
	context.WriteNode(call.Callee)
	if len(call.Arguments) != 0 {
		printer.WriteNodeList(context.WithNest(), call.Arguments)
	}
}

type IndexExpression struct {
	expression
	Left  Expression
	Index Expression
}

func (index *IndexExpression) GetLocation() text.Location {
	return index.Left.GetLocation()
}

func (index *IndexExpression) String(context printer.Context) {
	context.Write("%sINDEX_EXPR", context.Colour(colour.NodeName))
	context.WriteNode(index.Left)
	if index.Index != nil {
		context.WriteNode(index.Index)
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

func (a *AssignmentExpression) String(context printer.Context) {
	context.Write(
		"%sASSIGNMENT_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Symbol),
		a.Operator.Value,
	)
	context.WriteNode(a.Assignee)
	context.WriteNode(a.Value)
}

type TupleExpression struct {
	expression
	Location text.Location
	Values   []Expression
}

func (t *TupleExpression) GetLocation() text.Location {
	return t.Location
}

func (t *TupleExpression) String(context printer.Context) {
	context.Write("%sTUPLE_EXPR", context.Colour(colour.NodeName))
	if len(t.Values) != 0 {
		printer.WriteNodeList(context.WithNest(), t.Values)
	}
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

func (m *MemberExpression) String(context printer.Context) {
	context.Write(
		"%sMEMBER_EXPR %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		m.Member,
	)
	context.WriteNode(m.Left)
}

type StructMember struct {
	Location text.Location
	Name     *string
	Value    Expression
}

func (sm StructMember) String(context printer.Context) {
	context.Write("%sSTRUCT_MEMBER", context.Colour(colour.NodeName))

	if sm.Name != nil {
		context.Write(" %s%s", context.Colour(colour.Name), *sm.Name)
	}

	if sm.Value != nil {
		context.WriteNode(sm.Value)
	}
}

type InferredExpression struct {
	expression
	Location text.Location
}

func (i *InferredExpression) GetLocation() text.Location {
	return i.Location
}

func (i *InferredExpression) String(context printer.Context) {
	context.Write("%sINFERRED_EXPR", context.Colour(colour.NodeName))
}

type StructExpression struct {
	expression
	Struct  Expression
	Members []StructMember
}

func (s *StructExpression) GetLocation() text.Location {
	return s.Struct.GetLocation()
}

func (s *StructExpression) String(context printer.Context) {
	context.Write("%sSTRUCT_EXPR", context.Colour(colour.NodeName))
	context.WriteNode(s.Struct)
	if len(s.Members) != 0 {
		printer.WriteNodeList(context.WithNest(), s.Members)
	}
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

func (ce *CastExpression) String(context printer.Context) {
	context.Write("%sCAST_EXPR", context.Colour(colour.NodeName))
	context.WriteNode(ce.Left)
	context.WriteNode(ce.Type)
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

func (tc *TypeCheckExpression) String(context printer.Context) {
	context.Write("%sTYPE_CHECK_EXPR", context.Colour(colour.NodeName))
	context.WriteNode(tc.Left)
	context.WriteNode(tc.Type)
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

func (r *RangeExpression) String(context printer.Context) {
	context.Write("%sRANGE_EXPR", context.Colour(colour.NodeName))
	context.WriteNode(r.Start)
	context.WriteNode(r.End)
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

func (f *FunctionExpression) String(context printer.Context) {
	if f.Body != nil {
		context.Write("%sFUNC_EXPR", context.Colour(colour.NodeName))
	} else {
		context.Write("%sFUNC_TYPE", context.Colour(colour.NodeName))
	}

	if len(f.Parameters) != 0 {
		printer.WriteNodeList(context.WithNest(), f.Parameters)
	}

	if f.ReturnType != nil {
		context.WriteNode(f.ReturnType)
	}

	if f.Body != nil {
		context.WriteNode(f.Body)
	}
}

type Block struct {
	expression
	Location   text.Location
	Statements []Statement
}

func (b *Block) String(context printer.Context) {
	context.Write("%sBLOCK", context.Colour(colour.NodeName))
	if len(b.Statements) != 0 {
		printer.WriteNodeList(context.WithNest(), b.Statements)
	}
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

func (i *IfExpression) String(context printer.Context) {
	context.Write("%sIF_EXPR", context.Colour(colour.NodeName))
	context.WriteNode(i.Condition)
	context.WriteNode(i.Body)

	if i.ElseBranch != nil {
		context = context.WithNest()
		context.Write("%sELSE_BRANCH", context.Colour(colour.NodeName))
		context.WriteNode(i.ElseBranch)
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

func (wl *WhileLoop) String(context printer.Context) {
	context.Write("%sWHILE_LOOP", context.Colour(colour.NodeName))
	context.WriteNode(wl.Condition)
	context.WriteNode(wl.Body)
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

func (fl *ForLoop) String(context printer.Context) {
	context.Write(
		"%sFOR_LOOP %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		fl.Variable,
	)
	context.WriteNode(fl.Iterator)
	context.WriteNode(fl.Body)
}

func (f *ForLoop) GetLocation() text.Location {
	return f.LLocation
}

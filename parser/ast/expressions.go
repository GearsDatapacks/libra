package ast

import (
	"fmt"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/lexer/token"
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
func (il *IntegerLiteral) String(context printContext) {
	context.write(
		"%sINT_LIT %s%d",
		context.colour(colour.NodeName),
		context.colour(colour.Literal),
		il.Value,
	)

	litValue := fmt.Sprint(il.Value)
	if litValue != il.Token.Value {
		context.write(" %s", il.Token.Value)
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
func (fl *FloatLiteral) String(context printContext) {
	context.write(
		"%sFLOAT_LIT %s%v",
		context.colour(colour.NodeName),
		context.colour(colour.Literal),
		fl.Value,
	)

	litValue := fmt.Sprint(fl.Value)
	if litValue != fl.Token.Value {
		context.write(" %s", fl.Token.Value)
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
func (bl *BooleanLiteral) String(context printContext) {
	context.write(
		"%sBOOL_LIT %s%t",
		context.colour(colour.NodeName),
		context.colour(colour.Literal),
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
func (sl *StringLiteral) String(context printContext) {
	context.write(
		"%sSTRING_LIT %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Literal),
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
func (i *Identifier) String(context printContext) {
	context.write(
		"%sIDENT %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
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

func (b *BinaryExpression) String(context printContext) {
	context.write(
		"%sBIN_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Symbol),
		b.Operator.Value,
	)
	context.writeNode(b.Left)
	context.writeNode(b.Right)
}

type ParenthesisedExpression struct {
	expression
	Location   text.Location
	Expression Expression
}

func (p *ParenthesisedExpression) GetLocation() text.Location {
	return p.Location
}

func (p *ParenthesisedExpression) String(context printContext) {
	context.write(
		"%sPAREN_EXPR",
		context.colour(colour.NodeName),
	)
	context.writeNode(p.Expression)
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

func (p *PrefixExpression) String(context printContext) {
	context.write(
		"%sPREFIX_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Symbol),
		p.Operator.String(),
	)
	context.writeNode(p.Operand)
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

func (p *PostfixExpression) String(context printContext) {
	context.write(
		"%sPOSTFIX_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Symbol),
		p.Operator,
	)
	context.writeNode(p.Operand)
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

func (p *PointerType) String(context printContext) {
	context.write("%sPTR_TYPE", context.colour(colour.NodeName))
	if p.Mutable {
		context.write(" %smut", context.colour(colour.Attribute))
	}

	context.writeNode(p.Operand)
}

type OptionType struct {
	expression
	Location text.Location
	Operand  Expression
}

func (p *OptionType) GetLocation() text.Location {
	return p.Location
}

func (p *OptionType) String(context printContext) {
	context.write("%sOPTION_TYPE", context.colour(colour.NodeName))
	context.writeNode(p.Operand)
}

type DerefExpression struct {
	expression
	Operand Expression
}

func (d *DerefExpression) GetLocation() text.Location {
	return d.Operand.GetLocation()
}

func (d *DerefExpression) String(context printContext) {
	context.write("%sDEREF_EXPR", context.colour(colour.NodeName))
	context.writeNode(d.Operand)
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

func (r *RefExpression) String(context printContext) {

	context.write("%sREF_EXPR", context.colour(colour.NodeName))
	if r.Mutable {
		context.write(" %smut", context.colour(colour.Attribute))
	}
	context.writeNode(r.Operand)
}

type ListLiteral struct {
	expression
	Location text.Location
	Values   []Expression
}

func (l *ListLiteral) GetLocation() text.Location {
	return l.Location
}

func (l *ListLiteral) String(context printContext) {
	context.write("%sLIST_EXPR", context.colour(colour.NodeName))
	if len(l.Values) != 0 {
		writeNodeList(context.withNest(), l.Values)
	}
}

type KeyValue struct {
	Key   Expression
	Value Expression
}

func (kv KeyValue) String(context printContext) {
	context.write("%sKEY_VALUE", context.colour(colour.NodeName))
	context.writeNode(kv.Key)
	context.writeNode(kv.Value)
}

type MapLiteral struct {
	expression
	Location  text.Location
	KeyValues []KeyValue
}

func (m *MapLiteral) GetLocation() text.Location {
	return m.Location
}

func (m *MapLiteral) String(context printContext) {
	context.write("%sMAP_EXPR", context.colour(colour.NodeName))
	if len(m.KeyValues) != 0 {
		writeNodeList(context.withNest(), m.KeyValues)
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

func (call *FunctionCall) String(context printContext) {
	context.write("%sFUNCTION_CALL", context.colour(colour.NodeName))
	context.writeNode(call.Callee)
	if len(call.Arguments) != 0 {
		writeNodeList(context.withNest(), call.Arguments)
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

func (index *IndexExpression) String(context printContext) {
	context.write("%sINDEX_EXPR", context.colour(colour.NodeName))
	context.writeNode(index.Left)
	if index.Index != nil {
		context.writeNode(index.Index)
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

func (a *AssignmentExpression) String(context printContext) {
	context.write(
		"%sASSIGNMENT_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Symbol),
		a.Operator.Value,
	)
	context.writeNode(a.Assignee)
	context.writeNode(a.Value)
}

type TupleExpression struct {
	expression
	Location text.Location
	Values   []Expression
}

func (t *TupleExpression) GetLocation() text.Location {
	return t.Location
}

func (t *TupleExpression) String(context printContext) {
	context.write("%sTUPLE_EXPR", context.colour(colour.NodeName))
	if len(t.Values) != 0 {
		writeNodeList(context.withNest(), t.Values)
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

func (m *MemberExpression) String(context printContext) {
	context.write(
		"%sMEMBER_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		m.Member,
	)
	context.writeNode(m.Left)
}

type StructMember struct {
	Location text.Location
	Name     *string
	Value    Expression
}

func (sm StructMember) String(context printContext) {
	context.write("%sSTRUCT_MEMBER", context.colour(colour.NodeName))

	if sm.Name != nil {
		context.write(" %s%s", context.colour(colour.Name), *sm.Name)
	}

	if sm.Value != nil {
		context.writeNode(sm.Value)
	}
}

type InferredExpression struct {
	expression
	Location text.Location
}

func (i *InferredExpression) GetLocation() text.Location {
	return i.Location
}

func (i *InferredExpression) String(context printContext) {
	context.write("%sINFERRED_EXPR", context.colour(colour.NodeName))
}

type StructExpression struct {
	expression
	Struct  Expression
	Members []StructMember
}

func (s *StructExpression) GetLocation() text.Location {
	return s.Struct.GetLocation()
}

func (s *StructExpression) String(context printContext) {
	context.write("%sSTRUCT_EXPR", context.colour(colour.NodeName))
	context.writeNode(s.Struct)
	if len(s.Members) != 0 {
		writeNodeList(context.withNest(), s.Members)
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

func (ce *CastExpression) String(context printContext) {
	context.write("%sCAST_EXPR", context.colour(colour.NodeName))
	context.writeNode(ce.Left)
	context.writeNode(ce.Type)
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

func (tc *TypeCheckExpression) String(context printContext) {
	context.write("%sTYPE_CHECK_EXPR", context.colour(colour.NodeName))
	context.writeNode(tc.Left)
	context.writeNode(tc.Type)
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

func (r *RangeExpression) String(context printContext) {
	context.write("%sRANGE_EXPR", context.colour(colour.NodeName))
	context.writeNode(r.Start)
	context.writeNode(r.End)
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

func (f *FunctionExpression) String(context printContext) {
	if f.Body != nil {
		context.write("%sFUNC_EXPR", context.colour(colour.NodeName))
	} else {
		context.write("%sFUNC_TYPE", context.colour(colour.NodeName))
	}

	if len(f.Parameters) != 0 {
		writeNodeList(context.withNest(), f.Parameters)
	}

	if f.ReturnType != nil {
		context.writeNode(f.ReturnType)
	}

	if f.Body != nil {
		context.writeNode(f.Body)
	}
}

type Block struct {
	expression
	Location   text.Location
	Statements []Statement
}

func (b *Block) String(context printContext) {
	context.write("%sBLOCK", context.colour(colour.NodeName))
	if len(b.Statements) != 0 {
		writeNodeList(context.withNest(), b.Statements)
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

func (i *IfExpression) String(context printContext) {
	context.write("%sIF_EXPR", context.colour(colour.NodeName))
	context.writeNode(i.Condition)
	context.writeNode(i.Body)

	if i.ElseBranch != nil {
		context = context.withNest()
		context.write("%sELSE_BRANCH", context.colour(colour.NodeName))
		context.writeNode(i.ElseBranch)
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

func (wl *WhileLoop) String(context printContext) {
	context.write("%sWHILE_LOOP", context.colour(colour.NodeName))
	context.writeNode(wl.Condition)
	context.writeNode(wl.Body)
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

func (fl *ForLoop) String(context printContext) {
	context.write(
		"%sFOR_LOOP %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		fl.Variable,
	)
	context.writeNode(fl.Iterator)
	context.writeNode(fl.Body)
}

func (f *ForLoop) GetLocation() text.Location {
	return f.LLocation
}

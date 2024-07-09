package ast

import (
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

func (il *IntegerLiteral) Tokens() []token.Token {
	return []token.Token{il.Token}
}
func (il *IntegerLiteral) Location() text.Location {
	return il.Token.Location
}
func (il *IntegerLiteral) String(context printContext) {
	context.write(
		"%sINT_LIT %s%d",
		context.colour(colour.NodeName),
		context.colour(colour.Literal),
		il.Value,
	)
}

type FloatLiteral struct {
	expression
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) Tokens() []token.Token {
	return []token.Token{fl.Token}
}
func (fl *FloatLiteral) Location() text.Location {
	return fl.Token.Location
}
func (fl *FloatLiteral) String(context printContext) {
	context.write(
		"%sFLOAT_LIT %s%f",
		context.colour(colour.NodeName),
		context.colour(colour.Literal),
		fl.Value,
	)
}

type BooleanLiteral struct {
	expression
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) Tokens() []token.Token {
	return []token.Token{bl.Token}
}
func (bl *BooleanLiteral) Location() text.Location {
	return bl.Token.Location
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

func (sl *StringLiteral) Tokens() []token.Token {
	return []token.Token{sl.Token}
}
func (sl *StringLiteral) Location() text.Location {
	return sl.Token.Location
}
func (sl *StringLiteral) String(context printContext) {
	context.write(
		"%sSTRING_LIT %s%q",
		context.colour(colour.NodeName),
		context.colour(colour.Literal),
		sl.Value,
	)
}

type Identifier struct {
	expression
	Token token.Token
	Name  string
}

func (i *Identifier) Tokens() []token.Token {
	return []token.Token{i.Token}
}
func (i *Identifier) Location() text.Location {
	return i.Token.Location
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

func (b *BinaryExpression) Tokens() []token.Token {
	tokens := []token.Token{}
	tokens = append(tokens, b.Left.Tokens()...)
	tokens = append(tokens, b.Operator)
	tokens = append(tokens, b.Right.Tokens()...)
	return tokens
}

func (b *BinaryExpression) Location() text.Location {
	return b.Left.Location().To(b.Right.Location())
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
	LeftParen  token.Token
	Expression Expression
	RightParen token.Token
}

func (p *ParenthesisedExpression) Tokens() []token.Token {
	tokens := []token.Token{p.LeftParen}
	tokens = append(tokens, p.Expression.Tokens()...)
	tokens = append(tokens, p.RightParen)
	return tokens
}

func (p *ParenthesisedExpression) Location() text.Location {
	return p.LeftParen.Location.To(p.RightParen.Location)
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
	Operator token.Token
	Operand  Expression
}

func (p *PrefixExpression) Tokens() []token.Token {
	return append([]token.Token{p.Operator}, p.Operand.Tokens()...)
}

func (p *PrefixExpression) Location() text.Location {
	return p.Operator.Location.To(p.Operand.Location())
}

func (p *PrefixExpression) String(context printContext) {
	context.write(
		"%sPREFIX_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Symbol),
		p.Operator.Value,
	)
	context.writeNode(p.Operand)
}

type PostfixExpression struct {
	expression
	Operand  Expression
	Operator token.Token
}

func (p *PostfixExpression) Tokens() []token.Token {
	return append(p.Operand.Tokens(), p.Operator)
}

func (p *PostfixExpression) Location() text.Location {
	return p.Operand.Location().To(p.Operator.Location)
}

func (p *PostfixExpression) String(context printContext) {
	context.write(
		"%sPOSTFIX_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Symbol),
		p.Operator.Value,
	)
	context.writeNode(p.Operand)
}

type PointerType struct {
	expression
	Operator token.Token
	Mutable  *token.Token
	Operand  Expression
}

func (p *PointerType) Tokens() []token.Token {
	tokens := []token.Token{p.Operator}
	if p.Mutable != nil {
		tokens = append(tokens, *p.Mutable)
	}
	return append(tokens, p.Operand.Tokens()...)
}

func (p *PointerType) Location() text.Location {
	return p.Operator.Location.To(p.Operand.Location())
}

func (p *PointerType) String(context printContext) {
	context.write("%sPTR_TYPE", context.colour(colour.NodeName))
	if p.Mutable != nil {
		context.write(" %smut", context.colour(colour.Attribute))
	}

	context.writeNode(p.Operand)
}

type OptionType struct {
	expression
	Operator token.Token
	Operand  Expression
}

func (p *OptionType) Tokens() []token.Token {
	tokens := []token.Token{p.Operator}
	return append(tokens, p.Operand.Tokens()...)
}

func (p *OptionType) Location() text.Location {
	return p.Operator.Location.To(p.Operand.Location())
}

func (p *OptionType) String(context printContext) {
	context.write("%sOPTION_TYPE", context.colour(colour.NodeName))
	context.writeNode(p.Operand)
}

type DerefExpression struct {
	expression
	Operand  Expression
	Operator token.Token
}

func (d *DerefExpression) Tokens() []token.Token {
	return append(d.Operand.Tokens(), d.Operator)
}

func (d *DerefExpression) Location() text.Location {
	return d.Operand.Location().To(d.Operator.Location)
}

func (d *DerefExpression) String(context printContext) {
	context.write("%sDEREF_EXPR", context.colour(colour.NodeName))
	context.writeNode(d.Operand)
}

type RefExpression struct {
	expression
	Operator token.Token
	Mutable  *token.Token
	Operand  Expression
}

func (r *RefExpression) Tokens() []token.Token {
	tokens := []token.Token{r.Operator}
	if r.Mutable != nil {
		tokens = append(tokens, *r.Mutable)
	}
	return append(tokens, r.Operand.Tokens()...)
}

func (r *RefExpression) Location() text.Location {
	return r.Operator.Location.To(r.Operand.Location())
}

func (r *RefExpression) String(context printContext) {

	context.write("%sREF_EXPR", context.colour(colour.NodeName))
	if r.Mutable != nil {
		context.write(" %smut", context.colour(colour.Attribute))
	}
	context.writeNode(r.Operand)
}

// We don't store the tokens of the commas because they probably won't be needed
type ListLiteral struct {
	expression
	LeftSquare  token.Token
	Values      []Expression
	RightSquare token.Token
}

func (l *ListLiteral) Tokens() []token.Token {
	tokens := []token.Token{l.LeftSquare}
	for _, value := range l.Values {
		tokens = append(tokens, value.Tokens()...)
	}
	tokens = append(tokens, l.RightSquare)
	return tokens
}

func (l *ListLiteral) Location() text.Location {
	return l.LeftSquare.Location.To(l.RightSquare.Location)
}

func (l *ListLiteral) String(context printContext) {
	context.write("%sLIST_EXPR", context.colour(colour.NodeName))
	if len(l.Values) != 0 {
		writeNodeList(context.withNest(), l.Values)
	}
}

type KeyValue struct {
	Key   Expression
	Colon token.Token
	Value Expression
}

func (kv *KeyValue) Tokens() []token.Token {
	tokens := kv.Key.Tokens()
	tokens = append(tokens, kv.Colon)
	tokens = append(tokens, kv.Value.Tokens()...)

	return tokens
}

func (kv KeyValue) String(context printContext) {
	context.write("%sKEY_VALUE", context.colour(colour.NodeName))
	context.writeNode(kv.Key)
	context.writeNode(kv.Value)
}

type MapLiteral struct {
	expression
	LeftBrace  token.Token
	KeyValues  []KeyValue
	RightBrace token.Token
}

func (m *MapLiteral) Location() text.Location {
	return m.LeftBrace.Location.To(m.RightBrace.Location)
}

func (m *MapLiteral) Tokens() []token.Token {
	tokens := []token.Token{m.LeftBrace}

	for _, kv := range m.KeyValues {
		tokens = append(tokens, kv.Tokens()...)
	}

	tokens = append(tokens, m.RightBrace)

	return tokens
}

func (m *MapLiteral) String(context printContext) {
	context.write("%sMAP_EXPR", context.colour(colour.NodeName))
	if len(m.KeyValues) != 0 {
		writeNodeList(context.withNest(), m.KeyValues)
	}
}

type FunctionCall struct {
	expression
	Callee     Expression
	LeftParen  token.Token
	Arguments  []Expression
	RightParen token.Token
}

func (call *FunctionCall) Tokens() []token.Token {
	tokens := append(call.Callee.Tokens(), call.LeftParen)

	for _, arg := range call.Arguments {
		tokens = append(tokens, arg.Tokens()...)
	}

	tokens = append(tokens, call.RightParen)

	return tokens
}

func (call *FunctionCall) Location() text.Location {
	return call.Callee.Location().To(call.RightParen.Location)
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
	Left        Expression
	LeftSquare  token.Token
	Index       Expression
	RightSquare token.Token
}

func (index *IndexExpression) Tokens() []token.Token {
	tokens := append(index.Left.Tokens(), index.LeftSquare)

	tokens = append(tokens, index.Index.Tokens()...)

	tokens = append(tokens, index.RightSquare)

	return tokens
}

func (index *IndexExpression) Location() text.Location {
	return index.Left.Location().To(index.RightSquare.Location)
}

func (index *IndexExpression) String(context printContext) {
	context.write("%sINDEX_EXPR", context.colour(colour.NodeName))
	context.writeNode(index.Index)
	context.writeNode(index.Left)
}

type AssignmentExpression struct {
	expression
	Assignee Expression
	Operator token.Token
	Value    Expression
}

func (a *AssignmentExpression) Tokens() []token.Token {
	tokens := append(a.Assignee.Tokens(), a.Operator)
	return append(tokens, a.Value.Tokens()...)
}

func (a *AssignmentExpression) Location() text.Location {
	return a.Assignee.Location().To(a.Value.Location())
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
	LeftParen  token.Token
	Values     []Expression
	RightParen token.Token
}

func (t *TupleExpression) Tokens() []token.Token {
	tokens := []token.Token{t.LeftParen}

	for _, value := range t.Values {
		tokens = append(tokens, value.Tokens()...)
	}

	tokens = append(tokens, t.RightParen)

	return tokens
}

func (t *TupleExpression) Location() text.Location {
	return t.LeftParen.Location.To(t.RightParen.Location)
}

func (t *TupleExpression) String(context printContext) {
	context.write("%sTUPLE_EXPR", context.colour(colour.NodeName))
	if len(t.Values) != 0 {
		writeNodeList(context.withNest(), t.Values)
	}
}

type MemberExpression struct {
	expression
	Left   Expression
	Dot    token.Token
	Member token.Token
}

func (m *MemberExpression) Tokens() []token.Token {
	if m.Left != nil {
		return append(m.Left.Tokens(), m.Dot, m.Member)
	}
	return []token.Token{m.Dot, m.Member}
}

func (m *MemberExpression) Location() text.Location {
	return m.Left.Location().To(m.Member.Location)
}

func (m *MemberExpression) String(context printContext) {
	context.write(
		"%sMEMBER_EXPR %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		m.Member.Value,
	)
	if m.Left != nil {
		context.writeNode(m.Left)
	}
}

type StructMember struct {
	Name  *token.Token
	Colon *token.Token
	Value Expression
}

func (sm *StructMember) Tokens() []token.Token {
	tokens := []token.Token{}
	if sm.Name != nil {
		tokens = append(tokens, *sm.Name)
	}
	if sm.Colon != nil {
		tokens = append(tokens, *sm.Colon)
	}
	return append(tokens, sm.Value.Tokens()...)
}

func (sm StructMember) String(context printContext) {
	context.write(
		"%sSTRUCT_MEMBER %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		sm.Name.Value,
	)
	context.writeNode(sm.Value)
}

type InferredExpression struct {
	expression
	Token token.Token
}

func (i *InferredExpression) Tokens() []token.Token {
	return []token.Token{i.Token}
}

func (i *InferredExpression) Location() text.Location {
	return i.Token.Location
}

func (i *InferredExpression) String(context printContext) {
	context.write("%sINFERRED_EXPR", context.colour(colour.NodeName))
}

type StructExpression struct {
	expression
	Struct     Expression
	LeftBrace  token.Token
	Members    []StructMember
	RightBrace token.Token
}

func (s *StructExpression) Tokens() []token.Token {
	tokens := append(s.Struct.Tokens(), s.LeftBrace)

	for _, member := range s.Members {
		tokens = append(tokens, member.Tokens()...)
	}

	return append(tokens, s.RightBrace)
}

func (s *StructExpression) Location() text.Location {
	return s.Struct.Location().To(s.RightBrace.Location)
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
	Left  Expression
	Arrow token.Token
	Type  Expression
}

func (ce *CastExpression) Tokens() []token.Token {
	tokens := ce.Left.Tokens()
	tokens = append(tokens, ce.Arrow)
	tokens = append(tokens, ce.Type.Tokens()...)

	return tokens
}

func (ce *CastExpression) Location() text.Location {
	return ce.Left.Location().To(ce.Type.Location())
}

func (ce *CastExpression) String(context printContext) {
	context.write("%sCAST_EXPR", context.colour(colour.NodeName))
	context.writeNode(ce.Left)
	context.writeNode(ce.Type)
}

type TypeCheckExpression struct {
	expression
	Left     Expression
	Operator token.Token
	Type     Expression
}

func (tc *TypeCheckExpression) Tokens() []token.Token {
	tokens := tc.Left.Tokens()
	tokens = append(tokens, tc.Operator)
	tokens = append(tokens, tc.Type.Tokens()...)

	return tokens
}

func (ce *TypeCheckExpression) Location() text.Location {
	return ce.Left.Location().To(ce.Type.Location())
}

func (tc *TypeCheckExpression) String(context printContext) {
	context.write("%sTYPE_CHECK_EXPR", context.colour(colour.NodeName))
	context.writeNode(tc.Left)
	context.writeNode(tc.Type)
}

type RangeExpression struct {
	expression
	Start    Expression
	Operator token.Token
	End      Expression
}

func (r *RangeExpression) Location() text.Location {
	return r.Start.Location().To(r.End.Location())
}

func (r *RangeExpression) Tokens() []token.Token {
	tokens := r.Start.Tokens()
	tokens = append(tokens, r.Operator)
	tokens = append(tokens, r.End.Tokens()...)

	return tokens
}

func (r *RangeExpression) String(context printContext) {
	context.write("%sRANGE_EXPR", context.colour(colour.NodeName))
	context.writeNode(r.Start)
	context.writeNode(r.End)
}

type FunctionExpression struct {
	expression
	Keyword    token.Token
	LeftParen  token.Token
	Parameters []Parameter
	RightParen token.Token
	ReturnType *TypeAnnotation
	Body       *Block
}

func (f *FunctionExpression) Location() text.Location {
	return f.Keyword.Location
}

func (f *FunctionExpression) Tokens() []token.Token {
	tokens := []token.Token{f.Keyword, f.LeftParen}
	for _, param := range f.Parameters {
		tokens = append(tokens, param.Tokens()...)
	}

	tokens = append(tokens, f.RightParen)
	if f.ReturnType != nil {
		tokens = append(tokens, f.ReturnType.Tokens()...)
	}
	if f.Body != nil {
		tokens = append(tokens, f.Body.Tokens()...)
	}

	return tokens
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
	LeftBrace  token.Token
	Statements []Statement
	RightBrace token.Token
}

func (b *Block) Tokens() []token.Token {
	tokens := []token.Token{b.LeftBrace}

	for _, stmt := range b.Statements {
		tokens = append(tokens, stmt.Tokens()...)
	}

	tokens = append(tokens, b.RightBrace)
	return tokens
}

func (b *Block) String(context printContext) {
	context.write("%sBLOCK", context.colour(colour.NodeName))
	if len(b.Statements) != 0 {
		writeNodeList(context.withNest(), b.Statements)
	}
}

func (b *Block) Location() text.Location {
	return b.LeftBrace.Location
}

type IfExpression struct {
	expression
	Keyword    token.Token
	Condition  Expression
	Body       *Block
	ElseBranch *ElseBranch
}

func (i *IfExpression) Tokens() []token.Token {
	tokens := []token.Token{i.Keyword}
	tokens = append(tokens, i.Condition.Tokens()...)
	tokens = append(tokens, i.Body.Tokens()...)

	if i.ElseBranch != nil {
		tokens = append(tokens, i.ElseBranch.ElseKeyword)
		tokens = append(tokens, i.ElseBranch.Statement.Tokens()...)
	}

	return tokens
}

func (i *IfExpression) String(context printContext) {
	context.write("%sIF_EXPR", context.colour(colour.NodeName))
	context.writeNode(i.Condition)
	context.writeNode(i.Body)

	if i.ElseBranch != nil {
		context = context.withNest()
		context.write("%sELSE_BRANCH", context.colour(colour.NodeName))
		context.writeNode(i.ElseBranch.Statement)
	}
}

func (i *IfExpression) Location() text.Location {
	return i.Keyword.Location
}

type ElseBranch struct {
	ElseKeyword token.Token
	Statement   Statement
}

type WhileLoop struct {
	expression
	Keyword   token.Token
	Condition Expression
	Body      *Block
}

func (wl *WhileLoop) Tokens() []token.Token {
	tokens := []token.Token{wl.Keyword}
	tokens = append(tokens, wl.Condition.Tokens()...)
	tokens = append(tokens, wl.Body.Tokens()...)
	return tokens
}

func (wl *WhileLoop) String(context printContext) {
	context.write("%sWHILE_LOOP", context.colour(colour.NodeName))
	context.writeNode(wl.Condition)
	context.writeNode(wl.Body)
}

func (w *WhileLoop) Location() text.Location {
	return w.Keyword.Location
}

type ForLoop struct {
	expression
	ForKeyword token.Token
	Variable   token.Token
	InKeyword  token.Token
	Iterator   Expression
	Body       *Block
}

func (fl *ForLoop) Tokens() []token.Token {
	tokens := []token.Token{fl.ForKeyword, fl.Variable, fl.InKeyword}
	tokens = append(tokens, fl.Iterator.Tokens()...)
	tokens = append(tokens, fl.Body.Tokens()...)
	return tokens
}

func (fl *ForLoop) String(context printContext) {
	context.write(
		"%sFOR_LOOP %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		fl.Variable.Value,
	)
	context.writeNode(fl.Iterator)
	context.writeNode(fl.Body)
}

func (f *ForLoop) Location() text.Location {
	return f.ForKeyword.Location
}

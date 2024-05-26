package ast

import (
	"bytes"

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
func (il *IntegerLiteral) String() string {
	return il.Token.Value
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
func (fl *FloatLiteral) String() string {
	return fl.Token.Value
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
func (bl *BooleanLiteral) String() string {
	return bl.Token.Value
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
func (sl *StringLiteral) String() string {
	return `"` + sl.Token.Value + `"`
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
func (i *Identifier) String() string {
	return i.Name
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

func (b *BinaryExpression) String() string {
	var result bytes.Buffer

	result.WriteString(b.Left.String())
	result.WriteByte(' ')
	result.WriteString(b.Operator.Value)
	result.WriteByte(' ')
	result.WriteString(b.Right.String())

	return result.String()
}

func (b *BinaryExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('(')

	result.WriteString(maybePrecedence(b.Left))

	result.WriteByte(' ')
	result.WriteString(b.Operator.Value)
	result.WriteByte(' ')

	result.WriteString(maybePrecedence(b.Right))

	result.WriteByte(')')

	return result.String()
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

func (p *ParenthesisedExpression) String() string {
	var result bytes.Buffer

	result.WriteByte('(')
	result.WriteString(p.Expression.String())
	result.WriteByte(')')

	return result.String()
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

func (p *PrefixExpression) String() string {
	return p.Operator.Value + p.Operand.String()
}

func (p *PrefixExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteString(p.Operator.Value)
	result.WriteByte('(')
	result.WriteString(maybePrecedence(p.Operand))
	result.WriteByte(')')

	return result.String()
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

func (p *PostfixExpression) String() string {
	return p.Operand.String() + p.Operator.Value
}

func (p *PostfixExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('(')
	result.WriteString(maybePrecedence(p.Operand))
	result.WriteByte(')')
	result.WriteString(p.Operator.Value)

	return result.String()
}

type DerefExpression struct {
	expression
	Operator token.Token
	Mutable  *token.Token
	Operand  Expression
}

func (d *DerefExpression) Tokens() []token.Token {
	tokens := []token.Token{d.Operator}
	if d.Mutable != nil {
		tokens = append(tokens, *d.Mutable)
	}
	return append(tokens, d.Operand.Tokens()...)
}

func (d *DerefExpression) Location() text.Location {
	return d.Operator.Location.To(d.Operand.Location())
}

func (d *DerefExpression) String() string {
	var result bytes.Buffer
	result.WriteByte('*')
	if d.Mutable != nil {
		result.WriteString("mut ")
	}

	result.WriteString(d.Operand.String())
	return result.String()
}

func (d *DerefExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('*')
	if d.Mutable != nil {
		result.WriteString("mut ")
	}

	result.WriteByte('(')
	result.WriteString(maybePrecedence(d.Operand))
	result.WriteByte(')')

	return result.String()
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

func (r *RefExpression) String() string {
	var result bytes.Buffer
	result.WriteByte('&')
	if r.Mutable != nil {
		result.WriteString("mut ")
	}

	result.WriteString(r.Operand.String())
	return result.String()
}

func (r *RefExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('&')
	if r.Mutable != nil {
		result.WriteString("mut ")
	}

	result.WriteByte('(')
	result.WriteString(maybePrecedence(r.Operand))
	result.WriteByte(')')

	return result.String()
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

func (l *ListLiteral) String() string {
	var result bytes.Buffer

	result.WriteByte('[')
	for i, value := range l.Values {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(value.String())
	}
	result.WriteByte(']')

	return result.String()
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

func (kv *KeyValue) String() string {
	var result bytes.Buffer

	result.WriteString(kv.Key.String())
	result.WriteString(": ")
	result.WriteString(kv.Value.String())

	return result.String()
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

func (m *MapLiteral) String() string {
	var result bytes.Buffer

	result.WriteByte('{')

	for i, kv := range m.KeyValues {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(kv.String())
	}

	result.WriteByte('}')

	return result.String()
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

func (call *FunctionCall) String() string {
	var result bytes.Buffer

	result.WriteString(call.Callee.String())
	result.WriteByte('(')

	for i, arg := range call.Arguments {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(arg.String())
	}

	result.WriteByte(')')

	return result.String()
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

func (index *IndexExpression) String() string {
	var result bytes.Buffer

	result.WriteString(index.Left.String())
	result.WriteByte('[')

	result.WriteString(index.Index.String())

	result.WriteByte(']')

	return result.String()
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

func (a *AssignmentExpression) String() string {
	var result bytes.Buffer

	result.WriteString(a.Assignee.String())
	result.WriteByte(' ')
	result.WriteString(a.Operator.Value)
	result.WriteByte(' ')
	result.WriteString(a.Value.String())

	return result.String()
}

func (a *AssignmentExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('(')
	result.WriteString(maybePrecedence(a.Assignee))
	result.WriteByte(' ')
	result.WriteString(a.Operator.Value)
	result.WriteByte(' ')
	result.WriteString(maybePrecedence(a.Value))
	result.WriteByte(')')

	return result.String()
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

func (t *TupleExpression) String() string {
	var result bytes.Buffer

	result.WriteByte('(')

	for i, value := range t.Values {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(value.String())
	}

	result.WriteByte(')')

	return result.String()
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

func (m *MemberExpression) String() string {
	var result bytes.Buffer

	if m.Left != nil {
		result.WriteString(m.Left.String())
	}
	result.WriteByte('.')
	result.WriteString(m.Member.Value)

	return result.String()
}

type StructMember struct {
	Name  token.Token
	Colon token.Token
	Value Expression
}

func (sm *StructMember) Tokens() []token.Token {
	return append([]token.Token{sm.Name, sm.Colon}, sm.Value.Tokens()...)
}

func (sm *StructMember) String() string {
	var result bytes.Buffer

	result.WriteString(sm.Name.Value)
	result.WriteString(": ")
	result.WriteString(sm.Value.String())

	return result.String()
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

func (i *InferredExpression) String() string {
	return i.Token.Value
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

func (s *StructExpression) String() string {
	var result bytes.Buffer

	result.WriteString(s.Struct.String())
	if _, isInferred := s.Struct.(*InferredExpression); !isInferred {
		result.WriteByte(' ')
	}
	result.WriteString("{ ")

	for i, member := range s.Members {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(member.String())
	}

	result.WriteString(" }")

	return result.String()
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

func (ce *CastExpression) String() string {
	var result bytes.Buffer

	result.WriteString(ce.Left.String())
	result.WriteString(" -> ")
	result.WriteString(ce.Type.String())

	return result.String()
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

func (tc *TypeCheckExpression) String() string {
	var result bytes.Buffer

	result.WriteString(tc.Left.String())
	result.WriteString(" is ")
	result.WriteString(tc.Type.String())

	return result.String()
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

func (r *RangeExpression) String() string {
	var result bytes.Buffer

	result.WriteString(r.Start.String())
	result.WriteString("..")
	result.WriteString(r.End.String())

	return result.String()
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

func (f *FunctionExpression) String() string {
	var result bytes.Buffer

	result.WriteString("fn(")

	for i, param := range f.Parameters {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.String())
	}

	result.WriteByte(')')
	if f.ReturnType != nil {
		result.WriteString(f.ReturnType.String())
	}
	
	if f.Body != nil {
		result.WriteByte(' ')
		result.WriteString(f.Body.String())
	}

	return result.String()
}

type HasPrecedence interface {
	Expression
	PrecedenceString() string
}

func maybePrecedence(expr Expression) string {
	if prec, ok := expr.(HasPrecedence); ok {
		return prec.PrecedenceString()
	}

	return expr.String()
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

func (b *Block) String() string {
	var result bytes.Buffer

	result.WriteByte('{')
	for _, stmt := range b.Statements {
		result.WriteByte('\n')
		result.WriteString(stmt.String())
	}
	result.WriteString("\n}")

	return result.String()
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

func (i *IfExpression) String() string {
	var result bytes.Buffer

	result.WriteString("if ")
	result.WriteString(i.Condition.String())
	result.WriteByte(' ')
	result.WriteString(i.Body.String())

	if i.ElseBranch != nil {
		result.WriteString(" else ")
		result.WriteString(i.ElseBranch.Statement.String())
	}

	return result.String()
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

func (wl *WhileLoop) String() string {
	var result bytes.Buffer

	result.WriteString("while ")
	result.WriteString(wl.Condition.String())
	result.WriteByte(' ')
	result.WriteString(wl.Body.String())

	return result.String()
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

func (fl *ForLoop) String() string {
	var result bytes.Buffer

	result.WriteString("for ")
	result.WriteString(fl.Variable.Value)
	result.WriteString(" in ")
	result.WriteString(fl.Iterator.String())
	result.WriteByte(' ')
	result.WriteString(fl.Body.String())

	return result.String()
}

func (f *ForLoop) Location() text.Location {
	return f.ForKeyword.Location
}

package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type TypeAnnotation struct {
	Colon token.Token
	Type  Expression
}

func (ta *TypeAnnotation) Tokens() []token.Token {
	tokens := []token.Token{ta.Colon}
	tokens = append(tokens, ta.Type.Tokens()...)

	return tokens
}

func (ta *TypeAnnotation) String() string {
	var result bytes.Buffer

	result.WriteString(": ")
	result.WriteString(ta.Type.String())

	return result.String()
}

type VariableDeclaration struct {
	Keyword    token.Token
	Identifier token.Token
	Type       *TypeAnnotation
	Equals     token.Token
	Value      Expression
	Attributes DeclarationAttributes
}

func (varDec *VariableDeclaration) Tokens() []token.Token {
	tokens := []token.Token{varDec.Keyword, varDec.Identifier}
	if varDec.Type != nil {
		tokens = append(tokens, varDec.Type.Tokens()...)
	}
	tokens = append(tokens, varDec.Equals)
	tokens = append(tokens, varDec.Value.Tokens()...)

	return tokens
}

func (varDec *VariableDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString(varDec.Keyword.Value)
	result.WriteByte(' ')
	result.WriteString(varDec.Identifier.Value)

	if varDec.Type != nil {
		result.WriteString(varDec.Type.String())
	}

	result.WriteString(" = ")
	result.WriteString(varDec.Value.String())

	return result.String()
}

func (v *VariableDeclaration) tryAddAttribute(attribute Attribute) bool {
	return v.Attributes.tryAddAttribute(attribute)
}

type DefaultValue struct {
	Equals token.Token
	Value  Expression
}

type TypeOrIdent struct {
	Name  *token.Token
	Colon *token.Token
	Type  Expression
}

func (t *TypeOrIdent) Tokens() []token.Token {
	tokens := []token.Token{}
	if t.Name != nil {
		tokens = append(tokens, *t.Name)
	}
	if t.Colon != nil {
		tokens = append(tokens, *t.Colon)
	}
	if t.Type != nil {
		tokens = append(tokens, t.Type.Tokens()...)
	}

	return tokens
}

func (t *TypeOrIdent) String() string {
	var result bytes.Buffer

	if t.Name != nil {
		result.WriteString(t.Name.Value)
	}
	if t.Colon != nil {
		result.WriteString(": ")
	}
	if t.Type != nil {
		result.WriteString(t.Type.String())
	}

	return result.String()
}

type Parameter struct {
	Mutable *token.Token
	TypeOrIdent
	Default *DefaultValue
}

func (p *Parameter) Tokens() []token.Token {
	tokens := []token.Token{}
	if p.Mutable != nil {
		tokens = append(tokens, *p.Mutable)
	}
	tokens = append(tokens, p.TypeOrIdent.Tokens()...)
	if p.Default != nil {
		tokens = append(tokens, p.Default.Equals)
		tokens = append(tokens, p.Default.Value.Tokens()...)
	}

	return tokens
}

func (p *Parameter) String() string {
	var result bytes.Buffer

	if p.Mutable != nil {
		result.WriteString("mut ")
	}
	result.WriteString(p.TypeOrIdent.String())
	if p.Default != nil {
		result.WriteString(" = ")
		result.WriteString(p.Default.Value.String())
	}

	return result.String()
}

type MethodOf struct {
	LeftParen  token.Token
	Mutable    *token.Token
	Type       Expression
	RightParen token.Token
}

func (m *MethodOf) Tokens() []token.Token {
	tokens := []token.Token{m.LeftParen}
	if m.Mutable != nil {
		tokens = append(tokens, *m.Mutable)
	}
	tokens = append(tokens, m.Type.Tokens()...)
	tokens = append(tokens, m.RightParen)

	return tokens
}

func (m *MethodOf) String() string {
	var result bytes.Buffer

	result.WriteByte('(')
	if m.Mutable != nil {
		result.WriteString("mut ")
	}
	result.WriteString(m.Type.String())
	result.WriteByte(')')

	return result.String()
}

type MemberOf struct {
	Name token.Token
	Dot  token.Token
}

func (m *MemberOf) Tokens() []token.Token {
	return []token.Token{m.Name, m.Dot}
}

func (m *MemberOf) String() string {
	var result bytes.Buffer

	result.WriteString(m.Name.Value)
	result.WriteByte('.')

	return result.String()
}

type FunctionDeclaration struct {
	decl
	Keyword    token.Token
	MethodOf   *MethodOf
	MemberOf   *MemberOf
	Name       token.Token
	LeftParen  token.Token
	Parameters []Parameter
	RightParen token.Token
	ReturnType *TypeAnnotation
	Body       *Block
	Implements *string
	Attributes DeclarationAttributes
}

func (fd *FunctionDeclaration) Tokens() []token.Token {
	tokens := []token.Token{fd.Keyword, fd.Name, fd.LeftParen}
	for _, param := range fd.Parameters {
		tokens = append(tokens, param.Tokens()...)
	}

	tokens = append(tokens, fd.RightParen)
	if fd.ReturnType != nil {
		tokens = append(tokens, fd.ReturnType.Tokens()...)
	}
	tokens = append(tokens, fd.Body.Tokens()...)

	return tokens
}

func (fd *FunctionDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString("fn ")
	if fd.MethodOf != nil {
		result.WriteString(fd.MethodOf.String())
		result.WriteByte(' ')
	}

	if fd.MemberOf != nil {
		result.WriteString(fd.MemberOf.String())
	}

	result.WriteString(fd.Name.Value)
	result.WriteByte('(')

	for i, param := range fd.Parameters {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.String())
	}

	result.WriteByte(')')
	if fd.ReturnType != nil {
		result.WriteString(fd.ReturnType.String())
	}
	result.WriteByte(' ')

	result.WriteString(fd.Body.String())

	return result.String()
}

func (f *FunctionDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "impl" {
		f.Implements = &attribute.(*TextAttribute).Text
		return true
	}

	return f.Attributes.tryAddAttribute(attribute)
}

type ReturnStatement struct {
	Keyword token.Token
	Value   Expression
}

func (r *ReturnStatement) Tokens() []token.Token {
	tokens := []token.Token{r.Keyword}
	if r.Value != nil {
		tokens = append(tokens, r.Value.Tokens()...)
	}
	return tokens
}

func (r *ReturnStatement) String() string {
	var result bytes.Buffer
	result.WriteString("return")

	if r.Value != nil {
		result.WriteByte(' ')
		result.WriteString(r.Value.String())
	}

	return result.String()
}

type YieldStatement struct {
	Keyword token.Token
	Value   Expression
}

func (y *YieldStatement) Tokens() []token.Token {
	tokens := []token.Token{y.Keyword}
	tokens = append(tokens, y.Value.Tokens()...)
	return tokens
}

func (y *YieldStatement) String() string {
	var result bytes.Buffer
	result.WriteString(" ")

	result.WriteString(y.Value.String())

	return result.String()
}

type BreakStatement struct {
	Keyword token.Token
	Value   Expression
}

func (b *BreakStatement) Tokens() []token.Token {
	tokens := []token.Token{b.Keyword}
	if b.Value != nil {
		tokens = append(tokens, b.Value.Tokens()...)
	}
	return tokens
}

func (b *BreakStatement) String() string {
	var result bytes.Buffer
	result.WriteString("break")

	if b.Value != nil {
		result.WriteByte(' ')
		result.WriteString(b.Value.String())
	}

	return result.String()
}

type ContinueStatement struct {
	Keyword token.Token
}

func (c *ContinueStatement) Tokens() []token.Token {
	return []token.Token{c.Keyword}
}

func (*ContinueStatement) String() string {
	return "continue"
}

type TypeDeclaration struct {
	decl
	expl
	Keyword    token.Token
	Name       token.Token
	Equals     token.Token
	Type       Expression
	Attributes DeclarationAttributes
}

func (t *TypeDeclaration) Tokens() []token.Token {
	tokens := []token.Token{t.Keyword, t.Name, t.Equals}
	tokens = append(tokens, t.Type.Tokens()...)
	return tokens
}

func (t *TypeDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString("type ")
	result.WriteString(t.Name.Value)
	result.WriteString(" = ")
	result.WriteString(t.Type.String())

	return result.String()
}

func (t *TypeDeclaration) tryAddAttribute(attribute Attribute) bool {
	return t.Attributes.tryAddAttribute(attribute)
}

type StructField struct {
	Pub *token.Token
	TypeOrIdent
}

func (s *StructField) Tokens() []token.Token {
	tokens := []token.Token{}
	if s.Pub != nil {
		tokens = append(tokens, *s.Pub)
	}
	tokens = append(tokens, s.TypeOrIdent.Tokens()...)

	return tokens
}

func (s *StructField) String() string {
	var result bytes.Buffer

	if s.Pub != nil {
		result.WriteString("pub ")
	}
	result.WriteString(s.TypeOrIdent.String())

	return result.String()
}

type StructDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	Body       *StructBody
	Tag        *string
	Attributes DeclarationAttributes
}

type StructBody struct {
	LeftBrace  token.Token
	Fields     []StructField
	RightBrace token.Token
}

type TupleStruct struct {
	LeftParen  token.Token
	Types      []Expression
	RightParen token.Token
}

func (s *StructDeclaration) Tokens() []token.Token {
	tokens := []token.Token{s.Keyword, s.Name}

	if s.Body != nil {
		tokens = append(tokens, s.Body.LeftBrace)
		for _, field := range s.Body.Fields {
			tokens = append(tokens, field.Tokens()...)
		}
		tokens = append(tokens, s.Body.RightBrace)
	}

	return tokens
}

func (s *StructDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString("struct ")
	result.WriteString(s.Name.Value)
	if s.Body != nil {
		result.WriteString(" {\n")
		for i, field := range s.Body.Fields {
			if i != 0 {
				result.WriteString(",\n")
			}
			result.WriteString(field.String())
		}
		result.WriteString("\n}")
	}

	return result.String()
}

func (s *StructDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "tag" {
		s.Tag = &attribute.(*TextAttribute).Text
		return true
	}

	return s.Attributes.tryAddAttribute(attribute)
}

type InterfaceMember struct {
	Name       token.Token
	LeftParen  token.Token
	Parameters []Expression
	RightParen token.Token
	ReturnType *TypeAnnotation
}

func (i *InterfaceMember) Tokens() []token.Token {
	tokens := []token.Token{i.Name, i.LeftParen}
	for _, param := range i.Parameters {
		tokens = append(tokens, param.Tokens()...)
	}
	tokens = append(tokens, i.RightParen)
	if i.ReturnType != nil {
		tokens = append(tokens, i.ReturnType.Tokens()...)
	}

	return tokens
}

func (i *InterfaceMember) String() string {
	var result bytes.Buffer

	result.WriteString(i.Name.Value)
	result.WriteRune('(')
	for i, param := range i.Parameters {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.String())
	}

	result.WriteByte(')')
	if i.ReturnType != nil {
		result.WriteString(i.ReturnType.String())
	}

	return result.String()
}

type InterfaceDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	LeftBrace  token.Token
	Members    []InterfaceMember
	RightBrace token.Token
	Attributes DeclarationAttributes
}

func (i *InterfaceDeclaration) Tokens() []token.Token {
	tokens := []token.Token{i.Keyword, i.Name, i.LeftBrace}

	for _, member := range i.Members {
		tokens = append(tokens, member.Tokens()...)
	}

	tokens = append(tokens, i.RightBrace)
	return tokens
}

func (i *InterfaceDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString("interface ")
	result.WriteString(i.Name.Value)
	result.WriteString(" {\n")

	for i, member := range i.Members {
		if i != 0 {
			result.WriteString(",\n")
		}
		result.WriteString(member.String())
	}

	result.WriteString("\n}")
	return result.String()
}

func (i *InterfaceDeclaration) tryAddAttribute(attribute Attribute) bool {
	return i.Attributes.tryAddAttribute(attribute)
}

type ImportAll struct {
	Star token.Token
	From token.Token
}

type ImportAlias struct {
	As    token.Token
	Alias token.Token
}

type ImportedSymbols struct {
	LeftBrace  token.Token
	Symbols    []token.Token
	RightBrace token.Token
	From       token.Token
}

func (s *ImportedSymbols) Tokens() []token.Token {
	tokens := []token.Token{s.LeftBrace}
	tokens = append(tokens, s.Symbols...)
	tokens = append(tokens, s.RightBrace, s.From)

	return tokens
}

func (s *ImportedSymbols) String() string {
	var result bytes.Buffer

	result.WriteString("{ ")
	for i, symbol := range s.Symbols {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(symbol.Value)
	}

	result.WriteString(" } from")

	return result.String()
}

type ImportStatement struct {
	Keyword token.Token
	Symbols *ImportedSymbols
	All     *ImportAll
	Module  token.Token
	Alias   *ImportAlias
}

func (i *ImportStatement) Tokens() []token.Token {
	tokens := []token.Token{i.Keyword}
	if i.Symbols != nil {
		tokens = append(tokens, i.Symbols.Tokens()...)
	}
	if i.All != nil {
		tokens = append(tokens, i.All.Star, i.All.From)
	}
	tokens = append(tokens, i.Module)
	if i.Alias != nil {
		tokens = append(tokens, i.Alias.As, i.Alias.Alias)
	}
	return tokens
}

func (i *ImportStatement) String() string {
	var result bytes.Buffer

	result.WriteString("import ")
	if i.Symbols != nil {
		result.WriteString(i.Symbols.String())
		result.WriteByte(' ')
	}
	if i.All != nil {
		result.WriteString("* from ")
	}
	result.WriteByte('"')
	result.WriteString(i.Module.Value)
	result.WriteByte('"')
	if i.Alias != nil {
		result.WriteString(" as ")
		result.WriteString(i.Alias.Alias.Value)
	}

	return result.String()
}

type EnumMember struct {
	Name  token.Token
	Value *DefaultValue
}

func (e *EnumMember) Tokens() []token.Token {
	tokens := []token.Token{e.Name}
	if e.Value != nil {
		tokens = append(tokens, e.Value.Equals)
		tokens = append(tokens, e.Value.Value.Tokens()...)
	}

	return tokens
}

func (e *EnumMember) String() string {
	var result bytes.Buffer
	result.WriteString(e.Name.Value)

	if e.Value != nil {
		result.WriteString(" = ")
		result.WriteString(e.Value.Value.String())
	}

	return result.String()
}

type EnumDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	ValueType  *TypeAnnotation
	LeftBrace  token.Token
	Members    []EnumMember
	RightBrace token.Token
	Tag        *string
	Attributes DeclarationAttributes
}

func (e *EnumDeclaration) Tokens() []token.Token {
	tokens := []token.Token{e.Keyword, e.Name, e.LeftBrace}
	if e.ValueType != nil {
		tokens = append(tokens, e.ValueType.Tokens()...)
	}

	for _, member := range e.Members {
		tokens = append(tokens, member.Tokens()...)
	}
	tokens = append(tokens, e.RightBrace)
	return tokens
}

func (e *EnumDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString(e.Keyword.Value)
	result.WriteByte(' ')
	result.WriteString(e.Name.Value)

	if e.ValueType != nil {
		result.WriteString(e.ValueType.String())
	}

	result.WriteString(" {\n")
	for i, member := range e.Members {
		if i != 0 {
			result.WriteString(",\n")
		}
		result.WriteString(member.String())
	}
	result.WriteString("\n}")

	return result.String()
}

func (e *EnumDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "tag" {
		e.Tag = &attribute.(*TextAttribute).Text
		return true
	}

	return e.Attributes.tryAddAttribute(attribute)
}

type UnionMember struct {
	Name     token.Token
	Type     *TypeAnnotation
	Compound *StructBody
}

func (u *UnionMember) Tokens() []token.Token {
	tokens := []token.Token{u.Name}
	if u.Type != nil {
		tokens = append(tokens, u.Type.Tokens()...)
	}
	if u.Compound != nil {
		tokens = append(tokens, u.Compound.LeftBrace)
		for _, field := range u.Compound.Fields {
			tokens = append(tokens, field.Tokens()...)
		}
		tokens = append(tokens, u.Compound.RightBrace)
	}

	return tokens
}

func (u *UnionMember) String() string {
	var result bytes.Buffer
	result.WriteString(u.Name.Value)

	if u.Type != nil {
		result.WriteString(u.Type.String())
	}
	if u.Compound != nil {
		result.WriteString(" {")
		if len(u.Compound.Fields) != 0 {
			result.WriteByte(' ')
		}
		for i, field := range u.Compound.Fields {
			if i != 0 {
				result.WriteString(", ")
			}
			result.WriteString(field.String())
		}
		if len(u.Compound.Fields) != 0 {
			result.WriteByte(' ')
		}
		result.WriteByte('}')
	}

	return result.String()
}

type UnionDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	LeftBrace  token.Token
	Members    []UnionMember
	RightBrace token.Token
	Untagged   bool
	Tag        *string
	Attributes DeclarationAttributes
}

func (u *UnionDeclaration) Tokens() []token.Token {
	tokens := []token.Token{u.Keyword, u.Name, u.LeftBrace}

	for _, member := range u.Members {
		tokens = append(tokens, member.Tokens()...)
	}
	tokens = append(tokens, u.RightBrace)
	return tokens
}

func (u *UnionDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString(u.Keyword.Value)
	result.WriteByte(' ')
	result.WriteString(u.Name.Value)

	result.WriteString(" {\n")
	for i, member := range u.Members {
		if i != 0 {
			result.WriteString(",\n")
		}
		result.WriteString(member.String())
	}
	result.WriteString("\n}")

	return result.String()
}

func (u *UnionDeclaration) tryAddAttribute(attribute Attribute) bool {
	switch attribute.GetName() {
	case "untagged":
		u.Untagged = true
	case "tag":
		u.Tag = &attribute.(*TextAttribute).Text
	default:
		return u.Attributes.tryAddAttribute(attribute)
	}

	return true
}

type TagDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	Attributes DeclarationAttributes
}

func (t *TagDeclaration) Tokens() []token.Token {
	return []token.Token{t.Keyword, t.Name}
}

func (t *TagDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString(t.Keyword.Value)
	result.WriteByte(' ')
	result.WriteString(t.Name.Value)

	return result.String()
}

func (t *TagDeclaration) tryAddAttribute(attribute Attribute) bool {
	return t.Attributes.tryAddAttribute(attribute)
}

type Declaration interface {
	Statement
	MarkExport()
}

type decl struct {
	Exported bool
}

func (d *decl) MarkExport() {
	d.Exported = true
}

type Explicit interface {
	Statement
	MarkExplicit()
}

type expl struct {
	Explicit bool
}

func (d *expl) MarkExplicit() {
	d.Explicit = true
}

type AcceptsAttributes interface {
	Statement
	tryAddAttribute(Attribute) bool
}

func TryAddAttribute(stmt Statement, attribute Attribute) bool {
	if acceptsAttributes, ok := stmt.(AcceptsAttributes); ok {
		return acceptsAttributes.tryAddAttribute(attribute)
	}
	return false
}

type DeclarationAttributes struct {
	TodoMessage       *string
	Documentation     string
	DeprecatedMessage *string
}

func (d *DeclarationAttributes) tryAddAttribute(attribute Attribute) bool {
	switch attribute.GetName() {
	case "todo":
		d.TodoMessage = &attribute.(*TextAttribute).Text
	case "doc":
		d.Documentation = attribute.(*TextAttribute).Text
	case "deprecated":
		d.DeprecatedMessage = &attribute.(*TextAttribute).Text
	default:
		return false
	}
	return true
}

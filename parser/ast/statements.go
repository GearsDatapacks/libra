package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/colour"
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

func (ta *TypeAnnotation) String(context printContext) {
	context.write("%sTYPE_ANNOTATION", context.colour(colour.NodeName))
	context.writeNode(ta.Type)
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

func (varDec *VariableDeclaration) String(context printContext) {
	context.write(
		"%sVAR_DECL %s%s %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		varDec.Keyword.Value,
		context.colour(colour.Attribute),
		varDec.Identifier.Value,
	)
	if varDec.Type != nil {
		context.writeNode(varDec.Type)
	}
	context.writeNode(varDec.Value)
	varDec.Attributes.String(context)
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

func (t *TypeOrIdent) String(context printContext) {
	context.write("%sTYPE_OR_IDENT", context.colour(colour.NodeName))
	if t.Name != nil {
		context.write(" %s%s", context.colour(colour.Name), t.Name.Value)
	}
	if t.Type != nil {
		context.writeNode(t.Type)
	}
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

func (p Parameter) String(context printContext) {
	context.write("%sPARAM", context.colour(colour.NodeName))

	if p.Mutable != nil {
		context.write(" %smut", context.colour(colour.Attribute))
	}

	context.writeNode(&p.TypeOrIdent)
	if p.Default != nil {
		context.writeNode(p.Default.Value)
	}
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

func (m *MethodOf) String(context printContext) {
	context.write("%sMETHOD_OF", context.colour(colour.NodeName))
	if m.Mutable != nil {
		context.write(" %smut", context.colour(colour.Name))
	}

	context.writeNode(m.Type)
}

type MemberOf struct {
	Name token.Token
	Dot  token.Token
}

func (m *MemberOf) Tokens() []token.Token {
	return []token.Token{m.Name, m.Dot}
}

func (m *MemberOf) String(context printContext) {
	context.write(
		"%sMEMBER_OF %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		m.Name.Value,
	)
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

func (fd *FunctionDeclaration) String(context printContext) {
	context.write(
		"%sFUNC_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		fd.Name.Value,
	)

	if fd.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}

	if fd.Implements != nil {
		context.write(
			" %simpl %s%s",
			context.colour(colour.Attribute),
			context.colour(colour.Name),
			*fd.Implements,
		)
	}

	if fd.MethodOf != nil {
		context.writeNode(fd.MethodOf)
	}

	if fd.MemberOf != nil {
		context.writeNode(fd.MemberOf)
	}

	if len(fd.Parameters) != 0 {
		writeNodeList(context.withNest(), fd.Parameters)
	}

	if fd.ReturnType != nil {
		context.writeNode(fd.ReturnType)
	}

	context.writeNode(fd.Body)
	fd.Attributes.String(context)
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

func (r *ReturnStatement) String(context printContext) {
	context.write("%sRETURN", context.colour(colour.NodeName))

	if r.Value != nil {
		context.writeNode(r.Value)
	}
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

func (y *YieldStatement) String(context printContext) {
	context.write("%sYIELD", context.colour(colour.NodeName))

	if y.Value != nil {
		context.writeNode(y.Value)
	}
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

func (b *BreakStatement) String(context printContext) {
	context.write("%sBREAK", context.colour(colour.NodeName))

	if b.Value != nil {
		context.writeNode(b.Value)
	}
}

type ContinueStatement struct {
	Keyword token.Token
}

func (c *ContinueStatement) Tokens() []token.Token {
	return []token.Token{c.Keyword}
}

func (*ContinueStatement) String(context printContext) {
	context.write("%sCONTINUE", context.colour(colour.NodeName))
}

type TypeDeclaration struct {
	decl
	expl
	Keyword    token.Token
	Name       token.Token
	Equals     token.Token
	Type       Expression
	Tag        Expression
	Attributes DeclarationAttributes
}

func (t *TypeDeclaration) Tokens() []token.Token {
	tokens := []token.Token{t.Keyword, t.Name, t.Equals}
	tokens = append(tokens, t.Type.Tokens()...)
	return tokens
}

func (t *TypeDeclaration) String(context printContext) {
	context.write(
		"%sTYPE_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		t.Name.Value,
	)

	if t.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}
	if t.Explicit {
		context.write(" %sexplicit", context.colour(colour.Attribute))
	}

	context.writeNode(t.Type)
	t.Attributes.String(context)
	if t.Tag != nil {
		nested := context.withNest()
		nested.write("%stag", nested.colour(colour.Attribute))
		nested.writeNode(t.Tag)
	}
}

func (t *TypeDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "tag" {
		if t.Explicit {
			t.Tag = attribute.(*ExpressionAttribute).Expression
			return true
		}
		// TODO: Add a proper error message for this
		return false
	}

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

func (s StructField) String(context printContext) {
	context.write("%sSTRUCT_FIELD", context.colour(colour.NodeName))
	if s.Pub != nil {
		context.write(" %spub", context.colour(colour.Attribute))
	}
	context.writeNode(&s.TypeOrIdent)
}

type StructDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	Body       *StructBody
	Tag        Expression
	Attributes DeclarationAttributes
}

type StructBody struct {
	LeftBrace  token.Token
	Fields     []StructField
	RightBrace token.Token
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

func (s *StructDeclaration) String(context printContext) {
	context.write(
		"%sSTRUCT_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		s.Name.Value,
	)

	if s.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}

	if s.Body != nil {
		if len(s.Body.Fields) != 0 {
			writeNodeList(context.withNest(), s.Body.Fields)
		}
	}
	s.Attributes.String(context)
	if s.Tag != nil {
		nested := context.withNest()
		nested.write("%stag", nested.colour(colour.Attribute))
		nested.writeNode(s.Tag)
	}
}

func (s *StructDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "tag" {
		s.Tag = attribute.(*ExpressionAttribute).Expression
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

func (i InterfaceMember) String(context printContext) {
	context.write(
		"%sINTERFACE_MEMBER %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		i.Name.Value,
	)
	if len(i.Parameters) != 0 {
		writeNodeList(context.withNest(), i.Parameters)
	}

	if i.ReturnType != nil {
		context.writeNode(i.ReturnType)
	}
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

func (i *InterfaceDeclaration) String(context printContext) {
	context.write(
		"%sINTERFACE_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		i.Name.Value,
	)

	if i.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}

	if len(i.Members) != 0 {
		writeNodeList(context.withNest(), i.Members)
	}
	i.Attributes.String(context)
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

func (s *ImportedSymbols) String(context printContext) {
	context.write("%sIMPORTED_SYMBOLS", context.colour(colour.NodeName))
	for _, symbol := range s.Symbols {
		context.write(" %s%s", context.colour(colour.Name), symbol.Value)
	}
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

func (i *ImportStatement) String(context printContext) {
	context.write("%sIMPORT ", context.colour(colour.NodeName))
	if i.All != nil {
		context.write("%s* ", context.colour(colour.Symbol))
	}
	context.write("%s%q", context.colour(colour.Literal), i.Module.Value)
	if i.Alias != nil {
		context.write(" %s%s", context.colour(colour.Name), i.Alias.Alias.Value)
	}

	if i.Symbols != nil {
		context.writeNode(i.Symbols)
	}
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

func (e EnumMember) String(context printContext) {
	context.write(
		"%sENUM_MEMBER %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		e.Name.Value,
	)

	if e.Value != nil {
		context.writeNode(e.Value.Value)
	}
}

type EnumDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	ValueType  *TypeAnnotation
	LeftBrace  token.Token
	Members    []EnumMember
	RightBrace token.Token
	Tag        Expression
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

func (e *EnumDeclaration) String(context printContext) {
	context.write(
		"%sENUM_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		e.Name.Value,
	)

	if e.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}


	if e.ValueType != nil {
		context.writeNode(e.ValueType)
	}

	if len(e.Members) != 0 {
		writeNodeList(context.withNest(), e.Members)
	}
	e.Attributes.String(context)
	if e.Tag != nil {
		nested := context.withNest()
		nested.write("%stag", nested.colour(colour.Attribute))
		nested.writeNode(e.Tag)
	}
}

func (e *EnumDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "tag" {
		e.Tag = attribute.(*ExpressionAttribute).Expression
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

func (u UnionMember) String(context printContext) {
	context.write(
		"%sUNION_MEMBER %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		u.Name.Value,
	)

	if u.Type != nil {
		context.writeNode(u.Type)
	}
	if u.Compound != nil {
		if len(u.Compound.Fields) != 0 {
			writeNodeList(context.withNest(), u.Compound.Fields)
		}
	}
}

type UnionDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	LeftBrace  token.Token
	Members    []UnionMember
	RightBrace token.Token
	Untagged   bool
	Tag        Expression
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

func (u *UnionDeclaration) String(context printContext) {
	context.write("%sUNION_DECL %s", context.colour(colour.NodeName), u.Name.Value)
	
	if u.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}
	if u.Untagged {
		context.write(" %suntagged", context.colour(colour.Attribute))
	}

	if len(u.Members) != 0 {
		writeNodeList(context.withNest(), u.Members)
	}

	u.Attributes.String(context)
	if u.Tag != nil {
		nested := context.withNest()
		nested.write("%stag", nested.colour(colour.Attribute))
		nested.writeNode(u.Tag)
	}
}

func (u *UnionDeclaration) tryAddAttribute(attribute Attribute) bool {
	switch attribute.GetName() {
	case "untagged":
		u.Untagged = true
	case "tag":
		u.Tag = attribute.(*ExpressionAttribute).Expression
	default:
		return u.Attributes.tryAddAttribute(attribute)
	}

	return true
}

type TagBody struct {
	LeftBrace  token.Token
	Types      []Expression
	RightBrace token.Token
}

func (t *TagBody) Tokens() []token.Token {
	tokens := []token.Token{t.LeftBrace}
	for _, ty := range t.Types {
		tokens = append(tokens, ty.Tokens()...)
	}
	tokens = append(tokens, t.RightBrace)
	return tokens
}

type TagDeclaration struct {
	decl
	Keyword    token.Token
	Name       token.Token
	Body       *TagBody
	Attributes DeclarationAttributes
}

func (t *TagDeclaration) Tokens() []token.Token {
	tokens := []token.Token{t.Keyword, t.Name}
	if t.Body != nil {
		tokens = append(tokens, t.Body.Tokens()...)
	}
	return tokens
}

func (t *TagDeclaration) String(context printContext) {
	context.write("%sTAG_DECL %s%s", context.colour(colour.NodeName), context.colour(colour.Name), t.Name.Value)

	if t.Body != nil {
		if len(t.Body.Types) != 0 {
			writeNodeList(context.withNest(), t.Body.Types)
		}
	}
	t.Attributes.String(context)
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

func (d *DeclarationAttributes) String(context printContext) {
	var result bytes.Buffer
	writer := context.writer
	context.writer = &result
	hasAttributes := false

	context = context.withNest()
	context.write("%sATTRIBUTES", context.colour(colour.NodeName))
	context = context.nested()
	if d.TodoMessage != nil {
		context.writeNest()
		context.write(
			"%stodo %s= %s%q",
			context.colour(colour.Attribute),
			context.colour(colour.Symbol),
			context.colour(colour.Literal),
			*d.TodoMessage,
		)
		hasAttributes = true
	}
	if d.Documentation != "" {
		context.writeNest()
		context.write(
			"%sdoc %s= %s%q",
			context.colour(colour.Attribute),
			context.colour(colour.Symbol),
			context.colour(colour.Literal),
			d.Documentation,
		)
		hasAttributes = true
	}
	if d.DeprecatedMessage != nil {
		context.writeNest()
		context.write(
			"%sdeprecated %s= %s%q",
			context.colour(colour.Attribute),
			context.colour(colour.Symbol),
			context.colour(colour.Literal),
			*d.DeprecatedMessage,
		)
		hasAttributes = true
	}

	context.writer = writer

	if hasAttributes {
		context.write(result.String())
	}
}

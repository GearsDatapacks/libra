package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
)

type VariableDeclaration struct {
	Keyword      token.Token
	NameLocation text.Location
	Name         string
	Type         Expression
	Value        Expression
	Attributes   DeclarationAttributes
}

func (varDec *VariableDeclaration) String(context printer.Context) {
	context.Write(
		"%sVAR_DECL %s%s %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Attribute),
		varDec.Keyword.Value,
		context.Colour(colour.Name),
		varDec.Name,
	)
	if varDec.Type != nil {
		context.WriteNode(varDec.Type)
	}
	context.WriteNode(varDec.Value)
	varDec.Attributes.String(context)
}

func (v *VariableDeclaration) GetLocation() text.Location {
	return v.Keyword.Location
}

func (v *VariableDeclaration) tryAddAttribute(attribute Attribute) bool {
	return v.Attributes.tryAddAttribute(attribute)
}

type TypeOrIdent struct {
	Location text.Location
	Name     *string
	Type     Expression
}

func (t *TypeOrIdent) String(context printer.Context) {
	context.Write("%sTYPE_OR_IDENT", context.Colour(colour.NodeName))
	if t.Name != nil {
		context.Write(" %s%s", context.Colour(colour.Name), *t.Name)
	}
	if t.Type != nil {
		context.WriteNode(t.Type)
	}
}

type Parameter struct {
	Location text.Location
	Mutable  bool
	TypeOrIdent
	Default Expression
}

func (p Parameter) String(context printer.Context) {
	context.Write("%sPARAM", context.Colour(colour.NodeName))

	if p.Mutable {
		context.Write(" %smut", context.Colour(colour.Attribute))
	}

	context.WriteNode(&p.TypeOrIdent)
	if p.Default != nil {
		context.WriteNode(p.Default)
	}
}

type MethodOf struct {
	Mutable bool
	Type    Expression
}

func (m *MethodOf) String(context printer.Context) {
	context.Write("%sMETHOD_OF", context.Colour(colour.NodeName))
	if m.Mutable {
		context.Write(" %smut", context.Colour(colour.Name))
	}

	context.WriteNode(m.Type)
}

type MemberOf struct {
	Location text.Location
	Name     string
}

type FunctionDeclaration struct {
	decl
	Location     text.Location
	NameLocation text.Location
	MethodOf     *MethodOf
	MemberOf     *MemberOf
	Name         string
	Parameters   []Parameter
	ReturnType   Expression
	Body         *Block
	Implements   *string
	Attributes   DeclarationAttributes
}

func (fd *FunctionDeclaration) String(context printer.Context) {
	context.Write(
		"%sFUNC_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		fd.Name,
	)

	if fd.Exported {
		context.Write(" %spub", context.Colour(colour.Attribute))
	}

	if fd.Implements != nil {
		context.Write(
			" %simpl %s%s",
			context.Colour(colour.Attribute),
			context.Colour(colour.Name),
			*fd.Implements,
		)
	}

	if fd.MemberOf != nil {
		context.Write(
			" %smethodof %s%s",
			context.Colour(colour.Attribute),
			context.Colour(colour.Name),
			fd.MemberOf.Name,
		)
	}

	if fd.MethodOf != nil {
		context.WriteNode(fd.MethodOf)
	}

	if len(fd.Parameters) != 0 {
		printer.WriteNodeList(context.WithNest(), fd.Parameters)
	}

	if fd.ReturnType != nil {
		context.WriteNode(fd.ReturnType)
	}

	context.WriteNode(fd.Body)
	fd.Attributes.String(context)
}

func (f *FunctionDeclaration) GetLocation() text.Location {
	return f.Location
}

func (f *FunctionDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "impl" {
		f.Implements = &attribute.(*TextAttribute).Text
		return true
	}

	return f.Attributes.tryAddAttribute(attribute)
}

type ReturnStatement struct {
	Location text.Location
	Value    Expression
}

func (r *ReturnStatement) String(context printer.Context) {
	context.Write("%sRETURN", context.Colour(colour.NodeName))

	if r.Value != nil {
		context.WriteNode(r.Value)
	}
}

func (r *ReturnStatement) GetLocation() text.Location {
	return r.Location
}

type YieldStatement struct {
	Location text.Location
	Value    Expression
}

func (y *YieldStatement) String(context printer.Context) {
	context.Write("%sYIELD", context.Colour(colour.NodeName))

	context.WriteNode(y.Value)
}

func (y *YieldStatement) GetLocation() text.Location {
	return y.Location
}

type BreakStatement struct {
	Location text.Location
	Value    Expression
}

func (b *BreakStatement) String(context printer.Context) {
	context.Write("%sBREAK", context.Colour(colour.NodeName))

	if b.Value != nil {
		context.WriteNode(b.Value)
	}
}

func (b *BreakStatement) GetLocation() text.Location {
	return b.Location
}

type ContinueStatement struct{ Location text.Location }

func (*ContinueStatement) String(context printer.Context) {
	context.Write("%sCONTINUE", context.Colour(colour.NodeName))
}

func (c *ContinueStatement) GetLocation() text.Location {
	return c.Location
}

type TypeDeclaration struct {
	decl
	expl
	Location   text.Location
	Name       string
	Type       Expression
	Tag        Expression
	Attributes DeclarationAttributes
}

func (t *TypeDeclaration) String(context printer.Context) {
	context.Write(
		"%sTYPE_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		t.Name,
	)

	if t.Exported {
		context.Write(" %spub", context.Colour(colour.Attribute))
	}
	if t.Explicit {
		context.Write(" %sexplicit", context.Colour(colour.Attribute))
	}

	context.WriteNode(t.Type)
	t.Attributes.String(context)
	if t.Tag != nil {
		nested := context.WithNest()
		nested.Write("%stag", nested.Colour(colour.Attribute))
		nested.WriteNode(t.Tag)
	}
}

func (t *TypeDeclaration) GetLocation() text.Location {
	return t.Location
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
	Location text.Location
	Pub      bool
	TypeOrIdent
}

func (s StructField) String(context printer.Context) {
	context.Write("%sSTRUCT_FIELD", context.Colour(colour.NodeName))
	if s.Pub {
		context.Write(" %spub", context.Colour(colour.Attribute))
	}
	context.WriteNode(&s.TypeOrIdent)
}

type StructDeclaration struct {
	decl
	Location     text.Location
	NameLocation text.Location
	Name         string
	Body         []StructField
	Tag          Expression
	Attributes   DeclarationAttributes
}

func (s *StructDeclaration) String(context printer.Context) {
	context.Write(
		"%sSTRUCT_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		s.Name,
	)

	if s.Exported {
		context.Write(" %spub", context.Colour(colour.Attribute))
	}

	if s.Body != nil && len(s.Body) != 0 {
		printer.WriteNodeList(context.WithNest(), s.Body)
	}
	s.Attributes.String(context)
	if s.Tag != nil {
		nested := context.WithNest()
		nested.Write("%stag", nested.Colour(colour.Attribute))
		nested.WriteNode(s.Tag)
	}
}

func (s *StructDeclaration) GetLocation() text.Location {
	return s.Location
}

func (s *StructDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "tag" {
		s.Tag = attribute.(*ExpressionAttribute).Expression
		return true
	}

	return s.Attributes.tryAddAttribute(attribute)
}

type InterfaceMember struct {
	Name       string
	Parameters []Expression
	ReturnType Expression
}

func (i InterfaceMember) String(context printer.Context) {
	context.Write(
		"%sINTERFACE_MEMBER %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		i.Name,
	)
	if len(i.Parameters) != 0 {
		printer.WriteNodeList(context.WithNest(), i.Parameters)
	}

	if i.ReturnType != nil {
		context.WriteNode(i.ReturnType)
	}
}

type InterfaceDeclaration struct {
	decl
	Location   text.Location
	Name       string
	Members    []InterfaceMember
	Attributes DeclarationAttributes
}

func (i *InterfaceDeclaration) String(context printer.Context) {
	context.Write(
		"%sINTERFACE_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		i.Name,
	)

	if i.Exported {
		context.Write(" %spub", context.Colour(colour.Attribute))
	}

	if len(i.Members) != 0 {
		printer.WriteNodeList(context.WithNest(), i.Members)
	}
	i.Attributes.String(context)
}

func (i *InterfaceDeclaration) GetLocation() text.Location {
	return i.Location
}

func (i *InterfaceDeclaration) tryAddAttribute(attribute Attribute) bool {
	return i.Attributes.tryAddAttribute(attribute)
}

type ImportedSymbol struct {
	Location text.Location
	Name     string
}

type ImportStatement struct {
	Location text.Location
	Symbols  []ImportedSymbol
	All      bool
	Module   token.Token
	Alias    *string
}

func (i *ImportStatement) String(context printer.Context) {
	context.Write("%sIMPORT ", context.Colour(colour.NodeName))
	if i.All {
		context.Write("%s* ", context.Colour(colour.Symbol))
	}
	context.Write("%s%q", context.Colour(colour.Literal), i.Module.Value)
	if i.Alias != nil {
		context.Write(" %s%s", context.Colour(colour.Name), *i.Alias)
	}

	if i.Symbols != nil {
		for _, symbol := range i.Symbols {
			context.Write(" %s%s", context.Colour(colour.Name), symbol.Name)
		}
	}
}

func (i *ImportStatement) GetLocation() text.Location {
	return i.Location
}

type EnumMember struct {
	Name  string
	Value Expression
}

func (e EnumMember) String(context printer.Context) {
	context.Write(
		"%sENUM_MEMBER %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		e.Name,
	)

	if e.Value != nil {
		context.WriteNode(e.Value)
	}
}

type EnumDeclaration struct {
	decl
	Location   text.Location
	Name       string
	ValueType  Expression
	Members    []EnumMember
	Tag        Expression
	Attributes DeclarationAttributes
}

func (e *EnumDeclaration) String(context printer.Context) {
	context.Write(
		"%sENUM_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		e.Name,
	)

	if e.Exported {
		context.Write(" %spub", context.Colour(colour.Attribute))
	}

	if e.ValueType != nil {
		context.WriteNode(e.ValueType)
	}

	if len(e.Members) != 0 {
		printer.WriteNodeList(context.WithNest(), e.Members)
	}
	e.Attributes.String(context)
	if e.Tag != nil {
		nested := context.WithNest()
		nested.Write("%stag", nested.Colour(colour.Attribute))
		nested.WriteNode(e.Tag)
	}
}

func (e *EnumDeclaration) GetLocation() text.Location {
	return e.Location
}

func (e *EnumDeclaration) tryAddAttribute(attribute Attribute) bool {
	if attribute.GetName() == "tag" {
		e.Tag = attribute.(*ExpressionAttribute).Expression
		return true
	}

	return e.Attributes.tryAddAttribute(attribute)
}

type UnionMember struct {
	NameLocation text.Location
	Name         string
	Type         Expression
	Compound     []StructField
}

func (u UnionMember) String(context printer.Context) {
	context.Write(
		"%sUNION_MEMBER %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		u.Name,
	)

	if u.Type != nil {
		context.WriteNode(u.Type)
	}
	if u.Compound != nil && len(u.Compound) != 0 {
		printer.WriteNodeList(context.WithNest(), u.Compound)
	}
}

type UnionDeclaration struct {
	decl
	Location   text.Location
	Name       string
	Members    []UnionMember
	Untagged   bool
	Tag        Expression
	Attributes DeclarationAttributes
}

func (u *UnionDeclaration) String(context printer.Context) {
	context.Write(
		"%sUNION_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		u.Name,
	)

	if u.Exported {
		context.Write(" %spub", context.Colour(colour.Attribute))
	}
	if u.Untagged {
		context.Write(" %suntagged", context.Colour(colour.Attribute))
	}

	if len(u.Members) != 0 {
		printer.WriteNodeList(context.WithNest(), u.Members)
	}

	u.Attributes.String(context)
	if u.Tag != nil {
		nested := context.WithNest()
		nested.Write("%stag", nested.Colour(colour.Attribute))
		nested.WriteNode(u.Tag)
	}
}

func (u *UnionDeclaration) GetLocation() text.Location {
	return u.Location
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

type TagDeclaration struct {
	decl
	Location   text.Location
	Name       string
	Body       []Expression
	Attributes DeclarationAttributes
}

func (t *TagDeclaration) String(context printer.Context) {
	context.Write("%sTAG_DECL %s%s", context.Colour(colour.NodeName), context.Colour(colour.Name), t.Name)

	if t.Body != nil && len(t.Body) != 0 {
		printer.WriteNodeList(context.WithNest(), t.Body)
	}
	t.Attributes.String(context)
}

func (t *TagDeclaration) GetLocation() text.Location {
	return t.Location
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

func (d *DeclarationAttributes) String(context printer.Context) {
	var result bytes.Buffer
	bufferContext := context.WithWriter(&result)
	hasAttributes := false

	bufferContext = bufferContext.WithNest()
	bufferContext.Write("%sATTRIBUTES", bufferContext.Colour(colour.NodeName))
	bufferContext = bufferContext.Nested()
	if d.TodoMessage != nil {
		bufferContext.WriteNest()
		bufferContext.Write(
			"%stodo %s= %s%q",
			bufferContext.Colour(colour.Attribute),
			bufferContext.Colour(colour.Symbol),
			bufferContext.Colour(colour.Literal),
			*d.TodoMessage,
		)
		hasAttributes = true
	}
	if d.Documentation != "" {
		bufferContext.WriteNest()
		bufferContext.Write(
			"%sdoc %s= %s%q",
			bufferContext.Colour(colour.Attribute),
			bufferContext.Colour(colour.Symbol),
			bufferContext.Colour(colour.Literal),
			d.Documentation,
		)
		hasAttributes = true
	}
	if d.DeprecatedMessage != nil {
		bufferContext.WriteNest()
		bufferContext.Write(
			"%sdeprecated %s= %s%q",
			bufferContext.Colour(colour.Attribute),
			bufferContext.Colour(colour.Symbol),
			bufferContext.Colour(colour.Literal),
			*d.DeprecatedMessage,
		)
		hasAttributes = true
	}

	if hasAttributes {
		context.Write(result.String())
	}
}

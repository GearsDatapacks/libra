package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/lexer/token"
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

func (varDec *VariableDeclaration) String(context printContext) {
	context.write(
		"%sVAR_DECL %s%s %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		varDec.Keyword.Value,
		context.colour(colour.Attribute),
		varDec.Name,
	)
	if varDec.Type != nil {
		context.writeNode(varDec.Type)
	}
	context.writeNode(varDec.Value)
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

func (t *TypeOrIdent) String(context printContext) {
	context.write("%sTYPE_OR_IDENT", context.colour(colour.NodeName))
	if t.Name != nil {
		context.write(" %s%s", context.colour(colour.Name), *t.Name)
	}
	if t.Type != nil {
		context.writeNode(t.Type)
	}
}

type Parameter struct {
	Location text.Location
	Mutable  bool
	TypeOrIdent
	Default Expression
}

func (p Parameter) String(context printContext) {
	context.write("%sPARAM", context.colour(colour.NodeName))

	if p.Mutable {
		context.write(" %smut", context.colour(colour.Attribute))
	}

	context.writeNode(&p.TypeOrIdent)
	if p.Default != nil {
		context.writeNode(p.Default)
	}
}

type MethodOf struct {
	Mutable bool
	Type    Expression
}

func (m *MethodOf) String(context printContext) {
	context.write("%sMETHOD_OF", context.colour(colour.NodeName))
	if m.Mutable {
		context.write(" %smut", context.colour(colour.Name))
	}

	context.writeNode(m.Type)
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

func (fd *FunctionDeclaration) String(context printContext) {
	context.write(
		"%sFUNC_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		fd.Name,
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

	if fd.MemberOf != nil {
		context.write(
			" %smethodof %s%s",
			context.colour(colour.Attribute),
			context.colour(colour.Name),
			fd.MemberOf.Name,
		)
	}

	if fd.MethodOf != nil {
		context.writeNode(fd.MethodOf)
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

func (r *ReturnStatement) String(context printContext) {
	context.write("%sRETURN", context.colour(colour.NodeName))

	if r.Value != nil {
		context.writeNode(r.Value)
	}
}

func (r *ReturnStatement) GetLocation() text.Location {
	return r.Location
}

type YieldStatement struct {
	Location text.Location
	Value    Expression
}

func (y *YieldStatement) String(context printContext) {
	context.write("%sYIELD", context.colour(colour.NodeName))

	context.writeNode(y.Value)
}

func (y *YieldStatement) GetLocation() text.Location {
	return y.Location
}

type BreakStatement struct {
	Location text.Location
	Value    Expression
}

func (b *BreakStatement) String(context printContext) {
	context.write("%sBREAK", context.colour(colour.NodeName))

	if b.Value != nil {
		context.writeNode(b.Value)
	}
}

func (b *BreakStatement) GetLocation() text.Location {
	return b.Location
}

type ContinueStatement struct{ Location text.Location }

func (*ContinueStatement) String(context printContext) {
	context.write("%sCONTINUE", context.colour(colour.NodeName))
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

func (t *TypeDeclaration) String(context printContext) {
	context.write(
		"%sTYPE_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		t.Name,
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

func (s StructField) String(context printContext) {
	context.write("%sSTRUCT_FIELD", context.colour(colour.NodeName))
	if s.Pub {
		context.write(" %spub", context.colour(colour.Attribute))
	}
	context.writeNode(&s.TypeOrIdent)
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

func (s *StructDeclaration) String(context printContext) {
	context.write(
		"%sSTRUCT_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		s.Name,
	)

	if s.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}

	if s.Body != nil && len(s.Body) != 0 {
		writeNodeList(context.withNest(), s.Body)
	}
	s.Attributes.String(context)
	if s.Tag != nil {
		nested := context.withNest()
		nested.write("%stag", nested.colour(colour.Attribute))
		nested.writeNode(s.Tag)
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

func (i InterfaceMember) String(context printContext) {
	context.write(
		"%sINTERFACE_MEMBER %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		i.Name,
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
	Location   text.Location
	Name       string
	Members    []InterfaceMember
	Attributes DeclarationAttributes
}

func (i *InterfaceDeclaration) String(context printContext) {
	context.write(
		"%sINTERFACE_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		i.Name,
	)

	if i.Exported {
		context.write(" %spub", context.colour(colour.Attribute))
	}

	if len(i.Members) != 0 {
		writeNodeList(context.withNest(), i.Members)
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

func (i *ImportStatement) String(context printContext) {
	context.write("%sIMPORT ", context.colour(colour.NodeName))
	if i.All {
		context.write("%s* ", context.colour(colour.Symbol))
	}
	context.write("%s%q", context.colour(colour.Literal), i.Module.Value)
	if i.Alias != nil {
		context.write(" %s%s", context.colour(colour.Name), *i.Alias)
	}

	if i.Symbols != nil {
		for _, symbol := range i.Symbols {
			context.write(" %s%s", context.colour(colour.Name), symbol.Name)
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

func (e EnumMember) String(context printContext) {
	context.write(
		"%sENUM_MEMBER %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		e.Name,
	)

	if e.Value != nil {
		context.writeNode(e.Value)
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

func (e *EnumDeclaration) String(context printContext) {
	context.write(
		"%sENUM_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		e.Name,
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

func (u UnionMember) String(context printContext) {
	context.write(
		"%sUNION_MEMBER %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		u.Name,
	)

	if u.Type != nil {
		context.writeNode(u.Type)
	}
	if u.Compound != nil && len(u.Compound) != 0 {
		writeNodeList(context.withNest(), u.Compound)
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

func (u *UnionDeclaration) String(context printContext) {
	context.write(
		"%sUNION_DECL %s%s",
		context.colour(colour.NodeName),
		context.colour(colour.Name),
		u.Name,
	)

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

func (t *TagDeclaration) String(context printContext) {
	context.write("%sTAG_DECL %s%s", context.colour(colour.NodeName), context.colour(colour.Name), t.Name)

	if t.Body != nil && len(t.Body) != 0 {
		writeNodeList(context.withNest(), t.Body)
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

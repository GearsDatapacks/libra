package ast

import (
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

func (varDec *VariableDeclaration) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sVAR_DECL %s%s %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Attribute),
		varDec.Keyword.Value,
		context.Colour(colour.Name),
		varDec.Name,
	)
	if varDec.Type != nil {
		context.QueueNode(varDec.Type)
	}
	context.QueueNode(varDec.Value)
	context.QueueNode(varDec.Attributes)
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

func (t *TypeOrIdent) Print(context *printer.Printer) {
	context.QueueInfo("%sTYPE_OR_IDENT", context.Colour(colour.NodeName))
	if t.Name != nil {
		context.AddInfo(" %s%s", context.Colour(colour.Name), *t.Name)
	}
	if t.Type != nil {
		context.QueueNode(t.Type)
	}
}

type Parameter struct {
	Location text.Location
	Mutable  bool
	TypeOrIdent
	Default Expression
}

func (p Parameter) Print(context *printer.Printer) {
	context.QueueInfo("%sPARAM", context.Colour(colour.NodeName))

	if p.Mutable {
		context.AddInfo(" %smut", context.Colour(colour.Attribute))
	}

	context.QueueNode(&p.TypeOrIdent)
	if p.Default != nil {
		context.QueueNode(p.Default)
	}
}

type MethodOf struct {
	Mutable bool
	Type    Expression
}

func (m *MethodOf) Print(context *printer.Printer) {
	context.QueueInfo("%sMETHOD_OF", context.Colour(colour.NodeName))
	if m.Mutable {
		context.AddInfo(" %smut", context.Colour(colour.Name))
	}

	context.QueueNode(m.Type)
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

func (fd *FunctionDeclaration) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sFUNC_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		fd.Name,
	)

	if fd.Exported {
		context.AddInfo(" %spub", context.Colour(colour.Attribute))
	}

	if fd.Implements != nil {
		context.AddInfo(
			" %simpl %s%s",
			context.Colour(colour.Attribute),
			context.Colour(colour.Name),
			*fd.Implements,
		)
	}

	if fd.MemberOf != nil {
		context.AddInfo(
			" %smethodof %s%s",
			context.Colour(colour.Attribute),
			context.Colour(colour.Name),
			fd.MemberOf.Name,
		)
	}

	if fd.MethodOf != nil {
		context.QueueNode(fd.MethodOf)
	}

	printer.QueueNodeList(context, fd.Parameters)

	if fd.ReturnType != nil {
		context.QueueNode(fd.ReturnType)
	}

	context.QueueNode(fd.Body)
	context.QueueNode(fd.Attributes)
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

func (r *ReturnStatement) Print(context *printer.Printer) {
	context.QueueInfo("%sRETURN", context.Colour(colour.NodeName))

	if r.Value != nil {
		context.QueueNode(r.Value)
	}
}

func (r *ReturnStatement) GetLocation() text.Location {
	return r.Location
}

type YieldStatement struct {
	Location text.Location
	Value    Expression
}

func (y *YieldStatement) Print(context *printer.Printer) {
	context.QueueInfo("%sYIELD", context.Colour(colour.NodeName))

	context.QueueNode(y.Value)
}

func (y *YieldStatement) GetLocation() text.Location {
	return y.Location
}

type BreakStatement struct {
	Location text.Location
	Value    Expression
}

func (b *BreakStatement) Print(context *printer.Printer) {
	context.QueueInfo("%sBREAK", context.Colour(colour.NodeName))

	if b.Value != nil {
		context.QueueNode(b.Value)
	}
}

func (b *BreakStatement) GetLocation() text.Location {
	return b.Location
}

type ContinueStatement struct{ Location text.Location }

func (*ContinueStatement) Print(context *printer.Printer) {
	context.QueueInfo("%sCONTINUE", context.Colour(colour.NodeName))
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

func (t *TypeDeclaration) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sTYPE_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		t.Name,
	)

	if t.Exported {
		context.AddInfo(" %spub", context.Colour(colour.Attribute))
	}
	if t.Explicit {
		context.AddInfo(" %sexplicit", context.Colour(colour.Attribute))
	}

	context.QueueNode(t.Type)
	context.QueueNode(t.Attributes)
	if t.Tag != nil {
		context.Nest()
		context.QueueInfo("%stag", context.Colour(colour.Attribute))
		context.QueueNode(t.Tag)
		context.UnNest()
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

func (s StructField) Print(context *printer.Printer) {
	context.QueueInfo("%sSTRUCT_FIELD", context.Colour(colour.NodeName))
	if s.Pub {
		context.AddInfo(" %spub", context.Colour(colour.Attribute))
	}
	context.QueueNode(&s.TypeOrIdent)
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

func (s *StructDeclaration) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sSTRUCT_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		s.Name,
	)

	if s.Exported {
		context.AddInfo(" %spub", context.Colour(colour.Attribute))
	}

	if s.Body != nil {
		printer.QueueNodeList(context, s.Body)
	}
	context.QueueNode(s.Attributes)
	if s.Tag != nil {
		context.Nest()
		context.QueueInfo("%stag", context.Colour(colour.Attribute))
		context.QueueNode(s.Tag)
		context.UnNest()
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

func (i InterfaceMember) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sINTERFACE_MEMBER %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		i.Name,
	)
	printer.QueueNodeList(context, i.Parameters)

	if i.ReturnType != nil {
		context.QueueNode(i.ReturnType)
	}
}

type InterfaceDeclaration struct {
	decl
	Location   text.Location
	Name       string
	Members    []InterfaceMember
	Attributes DeclarationAttributes
}

func (i *InterfaceDeclaration) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sINTERFACE_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		i.Name,
	)

	if i.Exported {
		context.AddInfo(" %spub", context.Colour(colour.Attribute))
	}

	printer.QueueNodeList(context, i.Members)
	context.QueueNode(i.Attributes)
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

func (i *ImportStatement) Print(context *printer.Printer) {
	context.QueueInfo("%sIMPORT ", context.Colour(colour.NodeName))
	if i.All {
		context.AddInfo("%s* ", context.Colour(colour.Symbol))
	}
	context.AddInfo("%s%q", context.Colour(colour.Literal), i.Module.Value)
	if i.Alias != nil {
		context.AddInfo(" %s%s", context.Colour(colour.Name), *i.Alias)
	}

	if i.Symbols != nil {
		for _, symbol := range i.Symbols {
			context.AddInfo(" %s%s", context.Colour(colour.Name), symbol.Name)
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

func (e EnumMember) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sENUM_MEMBER %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		e.Name,
	)

	if e.Value != nil {
		context.QueueNode(e.Value)
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

func (e *EnumDeclaration) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sENUM_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		e.Name,
	)

	if e.Exported {
		context.AddInfo(" %spub", context.Colour(colour.Attribute))
	}

	if e.ValueType != nil {
		context.QueueNode(e.ValueType)
	}

	printer.QueueNodeList(context, e.Members)
	context.QueueNode(e.Attributes)
	if e.Tag != nil {
		context.Nest()
		context.QueueInfo("%stag", context.Colour(colour.Attribute))
		context.QueueNode(e.Tag)
		context.UnNest()
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

func (u UnionMember) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sUNION_MEMBER %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		u.Name,
	)

	if u.Type != nil {
		context.QueueNode(u.Type)
	}
	if u.Compound != nil && len(u.Compound) != 0 {
		printer.QueueNodeList(context, u.Compound)
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

func (u *UnionDeclaration) Print(context *printer.Printer) {
	context.QueueInfo(
		"%sUNION_DECL %s%s",
		context.Colour(colour.NodeName),
		context.Colour(colour.Name),
		u.Name,
	)

	if u.Exported {
		context.AddInfo(" %spub", context.Colour(colour.Attribute))
	}
	if u.Untagged {
		context.AddInfo(" %suntagged", context.Colour(colour.Attribute))
	}

	printer.QueueNodeList(context, u.Members)

	context.QueueNode(u.Attributes)
	if u.Tag != nil {
		context.Nest()
		context.QueueInfo("%stag", context.Colour(colour.Attribute))
		context.QueueNode(u.Tag)
		context.UnNest()
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

func (t *TagDeclaration) Print(context *printer.Printer) {
	context.QueueInfo("%sTAG_DECL %s%s", context.Colour(colour.NodeName), context.Colour(colour.Name), t.Name)

	if t.Body != nil && len(t.Body) != 0 {
		printer.QueueNodeList(context, t.Body)
	}
	context.QueueNode(t.Attributes)
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

func (d DeclarationAttributes) Print(context *printer.Printer) {
	hasAttributes := false

	context.QueueInfo("%sATTRIBUTES", context.Colour(colour.NodeName))

	context.Nest()
	defer context.UnNest()
	if d.TodoMessage != nil {
		context.QueueInfo(
			"%stodo %s= %s%q",
			context.Colour(colour.Attribute),
			context.Colour(colour.Symbol),
			context.Colour(colour.Literal),
			*d.TodoMessage,
		)
		hasAttributes = true
	}
	if d.Documentation != "" {
		context.QueueInfo(
			"%sdoc %s= %s%q",
			context.Colour(colour.Attribute),
			context.Colour(colour.Symbol),
			context.Colour(colour.Literal),
			d.Documentation,
		)
		hasAttributes = true
	}
	if d.DeprecatedMessage != nil {
		context.QueueInfo(
			"%sdeprecated %s= %s%q",
			context.Colour(colour.Attribute),
			context.Colour(colour.Symbol),
			context.Colour(colour.Literal),
			*d.DeprecatedMessage,
		)
		hasAttributes = true
	}

	if !hasAttributes {
		context.RejectNode()
	}
}

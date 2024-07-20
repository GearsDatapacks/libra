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

func (varDec *VariableDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sVAR_DECL %s%s %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Attribute),
			varDec.Keyword.Value,
			node.Colour(colour.Name),
			varDec.Name,
		).
		Location(varDec).
		OptionalNode(varDec.Type).
		Node(varDec.Value).
		Node(varDec.Attributes)
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

func (t *TypeOrIdent) Print(node *printer.Node) {
	node.Text("%sTYPE_OR_IDENT", node.Colour(colour.NodeName))
	if t.Name != nil {
		node.Text(" %s%s", node.Colour(colour.Name), *t.Name)
	}
	node.OptionalNode(t.Type)
}

type Parameter struct {
	Location text.Location
	Mutable  bool
	TypeOrIdent
	Default Expression
}

func (p Parameter) Print(node *printer.Node) {
	node.
		Text("%sPARAM", node.Colour(colour.NodeName)).
		TextIf(
			p.Mutable,
			" %smut",
			node.Colour(colour.Attribute),
		).
		Node(&p.TypeOrIdent).
		OptionalNode(p.Default)
}

type MethodOf struct {
	Mutable bool
	Type    Expression
}

func (m *MethodOf) Print(node *printer.Node) {
	node.
		Text("%sMETHOD_OF", node.Colour(colour.NodeName)).
		TextIf(
			m.Mutable,
			" %smut",
			node.Colour(colour.Name),
		).
		Node(m.Type)
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

func (fd *FunctionDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sFUNC_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			fd.Name,
		).
		TextIf(fd.Exported, " %spub", node.Colour(colour.Attribute))

	if fd.Implements != nil {
		node.Text(
			" %simpl %s%s",
			node.Colour(colour.Attribute),
			node.Colour(colour.Name),
			*fd.Implements,
		)
	}

	if fd.MemberOf != nil {
		node.Text(
			" %smethodof %s%s",
			node.Colour(colour.Attribute),
			node.Colour(colour.Name),
			fd.MemberOf.Name,
		)
	}
	node.
		Location(fd).
		OptionalNode(fd.MethodOf)

	printer.Nodes(node, fd.Parameters)

	node.
		OptionalNode(fd.ReturnType).
		Node(fd.Body).
		Node(fd.Attributes)
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

func (r *ReturnStatement) Print(node *printer.Node) {
	node.
		Text("%sRETURN", node.Colour(colour.NodeName)).
		Location(r).
		OptionalNode(r.Value)
}

func (r *ReturnStatement) GetLocation() text.Location {
	return r.Location
}

type YieldStatement struct {
	Location text.Location
	Value    Expression
}

func (y *YieldStatement) Print(node *printer.Node) {
	node.
		Text("%sYIELD", node.Colour(colour.NodeName)).
		Location(y).
		Node(y.Value)
}

func (y *YieldStatement) GetLocation() text.Location {
	return y.Location
}

type BreakStatement struct {
	Location text.Location
	Value    Expression
}

func (b *BreakStatement) Print(node *printer.Node) {
	node.
		Text("%sBREAK", node.Colour(colour.NodeName)).
		Location(b).
		OptionalNode(b.Value)
}

func (b *BreakStatement) GetLocation() text.Location {
	return b.Location
}

type ContinueStatement struct{ Location text.Location }

func (c *ContinueStatement) Print(node *printer.Node) {
	node.
		Text("%sCONTINUE", node.Colour(colour.NodeName)).
		Location(c)
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

func (t *TypeDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sTYPE_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			t.Name,
		).
		TextIf(t.Exported, " %spub", node.Colour(colour.Attribute)).
		TextIf(t.Explicit, " %sexplicit", node.Colour(colour.Attribute)).
		Location(t).
		Node(t.Type).
		Node(t.Attributes)
	if t.Tag != nil {
		node.FakeNode("%stag", func(node *printer.Node) {
			node.Node(t.Tag)
		}, node.Colour(colour.Attribute))
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

func (s StructField) Print(node *printer.Node) {
	node.
		Text("%sSTRUCT_FIELD", node.Colour(colour.NodeName)).
		TextIf(s.Pub, " %spub", node.Colour(colour.Attribute)).
		Node(&s.TypeOrIdent)
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

func (s *StructDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sSTRUCT_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			s.Name,
		).
		TextIf(s.Exported, " %spub", node.Colour(colour.Attribute)).
		Location(s)

	if s.Body != nil {
		printer.Nodes(node, s.Body)
	}

	node.Node(s.Attributes)

	if s.Tag != nil {
		node.FakeNode("%stag", func(node *printer.Node) {
			node.Node(s.Tag)
		}, node.Colour(colour.Attribute))
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

func (i InterfaceMember) Print(node *printer.Node) {
	node.Text(
		"%sINTERFACE_MEMBER %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		i.Name,
	)
	printer.Nodes(node, i.Parameters)

	node.OptionalNode(i.ReturnType)
}

type InterfaceDeclaration struct {
	decl
	Location   text.Location
	Name       string
	Members    []InterfaceMember
	Attributes DeclarationAttributes
}

func (i *InterfaceDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sINTERFACE_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			i.Name,
		).
		TextIf(i.Exported, " %spub", node.Colour(colour.Attribute)).
		Location(i)

	printer.Nodes(node, i.Members)
	node.Node(i.Attributes)
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

func (i *ImportStatement) Print(node *printer.Node) {
	node.
		Text("%sIMPORT ", node.Colour(colour.NodeName)).
		TextIf(i.All, "%s* ", node.Colour(colour.Symbol)).
		Text("%s%q", node.Colour(colour.Literal), i.Module.Value)

	if i.Alias != nil {
		node.
			Text(" %s%s", node.Colour(colour.Name), *i.Alias)
	}

	if i.Symbols != nil {
		for _, symbol := range i.Symbols {
			node.
				Text(" %s%s", node.Colour(colour.Name), symbol.Name)
		}
	}
	node.Location(i)
}

func (i *ImportStatement) GetLocation() text.Location {
	return i.Location
}

type EnumMember struct {
	Name  string
	Value Expression
}

func (e EnumMember) Print(node *printer.Node) {
	node.
		Text(
			"%sENUM_MEMBER %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			e.Name,
		).
		OptionalNode(e.Value)
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

func (e *EnumDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sENUM_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			e.Name,
		).
		TextIf(e.Exported, " %spub", node.Colour(colour.Attribute)).
		Location(e).
		OptionalNode(e.ValueType)

	printer.Nodes(node, e.Members)
	node.Node(e.Attributes)
	if e.Tag != nil {
		node.FakeNode("%stag", func(node *printer.Node) {
			node.Node(e.Tag)
		}, node.Colour(colour.Attribute))
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

func (u UnionMember) Print(node *printer.Node) {
	node.
		Text(
			"%sUNION_MEMBER %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			u.Name,
		).
		OptionalNode(u.Type)
	if u.Compound != nil && len(u.Compound) != 0 {
		printer.Nodes(node, u.Compound)
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

func (u *UnionDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sUNION_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			u.Name,
		).
		TextIf(u.Exported, " %spub", node.Colour(colour.Attribute)).
		TextIf(u.Untagged, " %suntagged", node.Colour(colour.Attribute)).
		Location(u)

	printer.Nodes(node, u.Members)

	node.Node(u.Attributes)
	if u.Tag != nil {
		node.FakeNode("%stag", func(node *printer.Node) {
			node.Node(u.Tag)
		}, node.Colour(colour.Attribute))
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

func (t *TagDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sTAG_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			t.Name,
		).
		Location(t)

	if t.Body != nil {
		printer.Nodes(node, t.Body)
	}
	node.Node(t.Attributes)
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

func (d DeclarationAttributes) Print(node *printer.Node) {
	node.
		Text("%sATTRIBUTES", node.Colour(colour.NodeName))
	hasAttributes := false

	if d.TodoMessage != nil {
		node.FakeNode(
			"%stodo %s= %s%q",
			nil,
			node.Colour(colour.Attribute),
			node.Colour(colour.Symbol),
			node.Colour(colour.Literal),
			*d.TodoMessage,
		)
		hasAttributes = true
	}
	if d.Documentation != "" {
		node.FakeNode(
			"%sdoc %s= %s%q",
			nil,
			node.Colour(colour.Attribute),
			node.Colour(colour.Symbol),
			node.Colour(colour.Literal),
			d.Documentation,
		)
		hasAttributes = true
	}
	if d.DeprecatedMessage != nil {
		node.FakeNode(
			"%sdeprecated %s= %s%q",
			nil,
			node.Colour(colour.Attribute),
			node.Colour(colour.Symbol),
			node.Colour(colour.Literal),
			*d.DeprecatedMessage,
		)
		hasAttributes = true
	}

	if !hasAttributes {
		node.Reject()
	}
}

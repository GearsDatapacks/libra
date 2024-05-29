package diagnostics

import (
	"fmt"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/text"
)

type Manager []Diagnostic

func (m *Manager) Report(diagnostics ...Diagnostic) {
	*m = append(*m, diagnostics...)
}

func makeError(msg string, location text.Location) Diagnostic {
	return new(Error, msg, location)
}

func makeInfo(msg string, location text.Location) Diagnostic {
	return new(Info, msg, location)
}

// Lexer Diagnostics

func InvalidCharacter(location text.Location, char byte) Diagnostic {
	msg := fmt.Sprintf("Invalid character: %q", char)
	return makeError(msg, location)
}

func UnterminatedString(location text.Location) Diagnostic {
	const msg = "Unterminated string"
	return makeError(msg, location)
}

func UnterminatedComment(location text.Location) Diagnostic {
	const msg = "Unterminated block comment"
	return makeError(msg, location)
}

func InvalidEscapeSequence(location text.Location, char byte) Diagnostic {
	msg := fmt.Sprintf("Invalid escape sequence: '\\%c'", char)
	return makeError(msg, location)
}

func ExpectedEscapeSequence(location text.Location) Diagnostic {
	const msg = "Expected escape sequence, reached end of file"
	return makeError(msg, location)
}

func InvalidAsciiSequence(location text.Location, sequence string) Diagnostic {
	msg := fmt.Sprintf("Invalid ascii escape sequence: '\\x%s'", sequence)
	return makeError(msg, location)
}

func InvalidUnicodeSequence(location text.Location, sequence string) Diagnostic {
	msg := fmt.Sprintf("Invalid unicode escape sequence: '\\x%s'", sequence)
	return makeError(msg, location)
}

func NumbersCannotEndWithSeparator(location text.Location) Diagnostic {
	const msg = "Numbers cannot end with numeric separators"
	return makeError(msg, location)
}

// Parser Diagnostics

func ExpectedExpression(location text.Location, kind token.Kind) Diagnostic {
	msg := fmt.Sprintf("Expected expression, found %s", kind.String())
	return makeError(msg, location)
}

func ExpectedNewline(location text.Location, kind token.Kind) Diagnostic {
	msg := fmt.Sprintf("Expected newline after statement, found %s", kind.String())
	return makeError(msg, location)
}

func ExpectedToken(location text.Location, expected token.Kind, actual token.Kind) Diagnostic {
	msg := fmt.Sprintf("Expected %s, found %s", expected.String(), actual.String())
	return makeError(msg, location)
}

func ElseStatementWithoutIf(location text.Location) Diagnostic {
	const msg = "Else statement not allowed without preceding if"
	return makeError(msg, location)
}

func ExpectedKeyword(location text.Location, keyword string, foundToken token.Token) Diagnostic {
	tokenValue := foundToken.Kind.String()
	if foundToken.Kind == token.IDENTIFIER {
		tokenValue = foundToken.Value
	}

	msg := fmt.Sprintf("Expected %q keyword, found %s", keyword, tokenValue)
	return makeError(msg, location)
}

func KeywordOverwritten(location text.Location, keyword string, declared text.Location) []Diagnostic {
	errMsg := fmt.Sprintf(
		"Expected %q keyword, but it has been overwritten by a variable",
		keyword)
	const info = "Try removing or renaming this variable"

	return []Diagnostic{makeError(errMsg, location), makeInfo(info, declared)}
}

func LastParameterMustHaveType(location text.Location, fnLocation text.Location) []Diagnostic {
	const msg = "The last parameter of a function must have a type annotation"
	diagnostics := []Diagnostic{makeError(msg, location)}

	if location.Span.StartLine != fnLocation.Span.StartLine {
		const info = "Parameter of this function"
		diagnostics = append(diagnostics, makeInfo(info, fnLocation))
	}
	return diagnostics
}

func MutWithoutParamName(location text.Location) Diagnostic {
	msg := `"mut" must be followed by a parameter name`
	return makeError(msg, location)
}

func LastStructFieldMustHaveType(location text.Location, structLoc text.Location) []Diagnostic {
	const errMsg = "The last field of a struct must have a type annotation"
	diagnostics := []Diagnostic{makeError(errMsg, location)}

	if location.Span.StartLine != structLoc.Span.StartLine {
		const info = "Field in this struct"
		diagnostics = append(diagnostics, makeInfo(info, structLoc))
	}
	return diagnostics
}

func MemberAndMethodNotAllowed(location text.Location) Diagnostic {
	const msg = "Functions cannot be both methods and static members"

	return makeError(msg, location)
}

func ExpectedMemberOrStructBody(location text.Location, tok token.Token) Diagnostic {
	msg := fmt.Sprintf("Invalid right-hand side of expression. Expected identifier or struct body, found %s", tok.Kind.String())

	return makeError(msg, location)
}

func OneImportModifierAllowed(location text.Location) Diagnostic {
	const msg = "Only one import modifier is allowed"

	return makeError(msg, location)
}

func OnlyTopLevelStatement(location text.Location, stmtKind string) Diagnostic {
	msg := fmt.Sprintf("%s not allowed here", stmtKind)

	return makeError(msg, location)
}

func ExpectedType(location text.Location, kind token.Kind) Diagnostic {
	msg := fmt.Sprintf("Expected type, found %s", kind.String())
	return makeError(msg, location)
}

func CannotExport(location text.Location) Diagnostic {
	const msg = "Only top-level declarations can be exported"

	return makeError(msg, location)
}

// Type-checker Diagnostics

type tcType interface {
	String() string
}

func UndefinedType(location text.Location, name string) Diagnostic {
	msg := fmt.Sprintf("Type %q is not defined", name)
	return makeError(msg, location)
}

func NotAssignable(location text.Location, expected, actual tcType) Diagnostic {
	msg := fmt.Sprintf("Value of type %q is not assignable to type %q", actual.String(), expected.String())
	return makeError(msg, location)
}

func VariableDefined(location text.Location, name string) Diagnostic {
	msg := fmt.Sprintf("Variable %q is already defined", name)
	return makeError(msg, location)
}

func VariableUndefined(location text.Location, name string) Diagnostic {
	msg := fmt.Sprintf("Variable %q is not defined", name)
	return makeError(msg, location)
}

func BinaryOperatorUndefined(location text.Location, operator string, left, right tcType) Diagnostic {
	msg := fmt.Sprintf("Operator %q is not defined for types %q and %q", operator, left.String(), right.String())
	return makeError(msg, location)
}

func UnaryOperatorUndefined(location text.Location, operator string, operand tcType) Diagnostic {
	msg := fmt.Sprintf("Operator %q is not defined for operand of type %q", operator, operand.String())
	return makeError(msg, location)
}

func CannotCast(location text.Location, from, to tcType) Diagnostic {
	msg := fmt.Sprintf("Cannot cast value of type %q to type %q", from.String(), to.String())
	return makeError(msg, location)
}

func CannotIncDec(incDec string) *Partial {
	msg := fmt.Sprintf("Cannot %s a non-variable value", incDec)
	return partial(Error, msg)
}

var ValueImmutablePartial = partial(Error, "Cannot modify value, it is immutable")

func ValueImmutable(location text.Location) Diagnostic {
	return ValueImmutablePartial.Location(location)
}

var NotConstPartial = partial(Error, "Value must be known at compile time")

func NotConst(location text.Location) Diagnostic {
	return NotConstPartial.Location(location)
}

func CountMustBeInt(location text.Location) Diagnostic {
	const msg = "Array length must be an integer"
	return makeError(msg, location)
}

func CannotIndex(leftType, indexType tcType) *Partial {
	msg := fmt.Sprintf("Cannot index value of type %q with value of type %q", leftType.String(), indexType.String())
	return partial(Error, msg)
}

func NotHashable(location text.Location, ty tcType) Diagnostic {
	msg := fmt.Sprintf("Value of type %q cannot be used as a key in a map", ty.String())
	return makeError(msg, location)
}

func CannotAssign(location text.Location) Diagnostic {
	const msg = "Cannot assign to a non-variable value"
	return makeError(msg, location)
}

func IndexOutOfBounds(index, len int64) *Partial {
	msg := fmt.Sprintf("Index %d is out of bounds of array of length %d", index, len)
	return partial(Error, msg)
}

func ConditionMustBeBool(location text.Location) Diagnostic {
	const msg = "Condition must be a boolean"
	return makeError(msg, location)
}

func NotIterable(location text.Location) Diagnostic {
	const msg = "Value is not iterable"
	return makeError(msg, location)
}

func NoReturnOutsideFunction(location text.Location) Diagnostic {
	const msg = "Cannot use return outside of a function"
	return makeError(msg, location)
}

func ExpectedReturnValue(location text.Location) Diagnostic {
	const msg = "Expected a return value"
	return makeError(msg, location)
}

func NotCallable(location text.Location, ty tcType) Diagnostic {
	msg := fmt.Sprintf("Value of type %q cannot be called", ty.String())
	return makeError(msg, location)
}

func WrongNumberAgruments(location text.Location, expected, actual int) Diagnostic {
	msg := fmt.Sprintf("Incorrect number of arguments (expected %d, found %d)", expected, actual)
	return makeError(msg, location)
}

func NoMember(leftType tcType, member string) *Partial {
	msg := fmt.Sprintf("Value of type %q does not have member %q", leftType.String(), member)
	return partial(Error, msg)
}

func FieldPrivate(leftType tcType, member string) *Partial {
	msg := fmt.Sprintf("Field %q of type %q is private", member, leftType.String())
	return partial(Error, msg)
}

func OnlyConstructTypes(location text.Location) Diagnostic {
	const msg = "Cannot construct value, not a type"
	return makeError(msg, location)
}

func CannotConstruct(location text.Location, ty tcType) Diagnostic {
	msg := fmt.Sprintf("Cannot construct value of type %q", ty.String())
	return makeError(msg, location)
}

func NoStructMember(location text.Location, name, member string) Diagnostic {
	msg := fmt.Sprintf("Struct %q does not have member %q", name, member)
	return makeError(msg, location)
}

func CannotUseStatementOutsideLoop(location text.Location, stmtKind string) Diagnostic {
	msg := fmt.Sprintf("Cannot use %s outside of a loop", stmtKind)
	return makeError(msg, location)
}

func CannotUseStatementOutsideBlock(location text.Location, stmtKind string) Diagnostic {
	msg := fmt.Sprintf("Cannot use %s outside of a block", stmtKind)
	return makeError(msg, location)
}

func ExpressionNotType(location text.Location, ty tcType) Diagnostic {
	msg := fmt.Sprintf("Expected a type, found value of type %q", ty.String())
	return makeError(msg, location)
}

func NamedParamInFnType(location text.Location) Diagnostic {
	const msg = "Parameters in function types must be unnamed"
	return makeError(msg, location)
}

func UnnamedParameter(location text.Location) Diagnostic {
	const msg = "Unnamed parameters are only allowed in function types"
	return makeError(msg, location)
}

func NoExport(location text.Location, module, member string) Diagnostic {
	msg := fmt.Sprintf("Module %q does not export member %q", module, member)
	return makeError(msg, location)
}

func CannotDeref(location text.Location, ty tcType) Diagnostic {
	msg := fmt.Sprintf("Cannot dereference non-pointer value of type %q", ty.String())
	return makeError(msg, location)
}

func MutDerefNotAllowed(location text.Location) Diagnostic {
	const msg = "`mut` keyword is only allowed in pointer types"
	return makeError(msg, location)
}

func MutRefOfNotMut(location text.Location) Diagnostic {
	const msg = "Cannot take a mutable reference to an immutable value"
	return makeError(msg, location)
}

func MixedNamedUnnamedStructFields(location text.Location) Diagnostic {
	const msg = "Cannot mix named and unnamed struct fields"
	return makeError(msg, location)
}

func WrongNumberTupleValues(location text.Location, expected, found int) Diagnostic {
	msg := fmt.Sprintf("Incorrect number of values supplied to struct (expected %d, found %d)", expected, found)
	return makeError(msg, location)
}

func TupleStructWithNames(location text.Location) Diagnostic {
	const msg = "Field names not allowed when constructing tuple structs"
	return makeError(msg, location)
}

func NoNameStructMember(location text.Location) Diagnostic {
	const msg = "Struct members must all be named"
	return makeError(msg, location)
}

func PubUnnamedStructField(location text.Location) Diagnostic {
	const msg = "`pub` keyword not allowed for unnamed fields"
	return makeError(msg, location)
}

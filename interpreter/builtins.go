package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
)

func toPrintString(value values.RuntimeValue) string {
	printStr := value.ToString()

	if _, ok := value.(*values.StringLiteral); ok {
		printStr = printStr[1 : len(printStr)-1]
	}

	return printStr
}

func print(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	fmt.Println(toPrintString(args[0]))

	return values.MakeNull()
}

func printil(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	fmt.Print(toPrintString(args[0]))

	return values.MakeNull()
}

var reader = bufio.NewReader(os.Stdin)

func prompt(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	fmt.Print(toPrintString(args[0]))

	result, _, _ := reader.ReadLine()

	return values.MakeString(string(result))
}

func to_string(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	return values.MakeString(args[0].ToString())
}

func parse_int(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	stringValue := args[0].(*values.StringLiteral).Value
	intValue, err := strconv.ParseInt(stringValue, 10, 32)

	if err != nil {
		return values.MakeError(fmt.Sprintf("parse_int: Invalid integer syntax: %q", stringValue))
	}

	return values.MakeInteger(int(intValue))
}

func parse_float(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	stringValue := args[0].(*values.StringLiteral).Value
	floatValue, err := strconv.ParseFloat(stringValue, 32)

	if err != nil {
		return values.MakeError(fmt.Sprintf("parse_float: Invalid float syntax: %q", stringValue))
	}

	return values.MakeFloat(floatValue)
}

func read_file(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	fileName := args[0].(*values.StringLiteral).Value
	file, err := os.ReadFile(fileName)
	if err != nil {
		return values.MakeError(err.Error())
	}

	return values.MakeString(string(file))
}

func write_file(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	fileName := args[0].(*values.StringLiteral).Value
	contents := args[1].(*values.StringLiteral).Value
	err := os.WriteFile(fileName, []byte(contents), 0666)
	if err != nil {
		return values.MakeError(err.Error())
	}

	return values.MakeNull()
}

package pdp

import "fmt"

type functionStringListLen struct {
	e  Expression
}

func makeFunctionStringListLen(e Expression) Expression {
	return functionStringListLen{e: e}
}

func makeFunctionStringListLenAlt(args []Expression) Expression {
	if len(args) != 1 {
		panic(fmt.Errorf("function \"len\" for List of Strings needs exactly one arguments but got %d", len(args)))
	}

	return makeFunctionStringListLen(args[0])
}

func (f functionStringListLen) GetResultType() Type {
	return TypeInteger
}

func (f functionStringListLen) describe() string {
	return "len"
}

// Calculate implements Expression interface and returns calculated value
func (f functionStringListLen) Calculate(ctx *Context) (AttributeValue, error) {
	s, err := ctx.calculateListOfStringExpression(f.e)
	if err != nil {
		return UndefinedValue, bindError(bindError(err, "argument"), f.describe())
	}

	return MakeIntegerValue(int64(len(s))), nil
}

func functionStringListLenValidator(args []Expression) functionMaker {
	if len(args) != 1 || args[0].GetResultType() != TypeListOfStrings {
		return nil
	}

	return makeFunctionStringListLenAlt
}

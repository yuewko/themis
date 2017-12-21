package pdp

import "fmt"

type functionIntegerDivide struct {
	first  Expression
	second Expression
}

func makeFunctionIntegerDivide(first, second Expression) Expression {
	return functionIntegerDivide{
		first:  first,
		second: second,
	}
}

func makeFunctionIntegerDivideAlt(args []Expression) Expression {
	if len(args) != 2 {
		panic(fmt.Errorf("function \"divide\" for Integer needs exactly two arguments but got %d", len(args)))
	}

	return makeFunctionIntegerDivide(args[0], args[1])
}

func (f functionIntegerDivide) GetResultType() int {
	return TypeInteger
}

func (f functionIntegerDivide) calculate(ctx *Context) (AttributeValue, error) {
	first, err := ctx.calculateIntegerExpression(f.first)
	if err != nil {
		return undefinedValue, bindError(bindError(err, "first argument"), "equal")
	}

	second, err := ctx.calculateIntegerExpression(f.second)
	if err != nil {
		return undefinedValue, bindError(bindError(err, "second argument"), "equal")
	}

	return MakeIntegerValue(first / second), nil
}

func functionIntegerDivideValidator(args []Expression) functionMaker {
	if len(args) != 2 || args[0].GetResultType() != TypeInteger || args[1].GetResultType() != TypeInteger {
		return nil
	}

	return makeFunctionIntegerDivideAlt
}

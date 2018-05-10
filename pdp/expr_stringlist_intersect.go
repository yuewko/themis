package pdp

import "fmt"

type functionStringListIntersect struct {
	first Expression
	second Expression
}

func makeFunctionStringListIntersect(first, second Expression) Expression {
	return functionStringListIntersect{
		first: first,
		second: second}
}

func makeFunctionStringListIntersectAlt(args []Expression) Expression {
	if len(args) != 2 {
		panic(fmt.Errorf("function \"intersect\" for List of Strings needs exactly two arguments but got %d", len(args)))
	}

	return makeFunctionStringListIntersect(args[0], args[1])
}

func (f functionStringListIntersect) GetResultType() Type {
	return TypeListOfStrings
}

func (f functionStringListIntersect) describe() string {
	return "intersect"
}

// Calculate implements Expression interface and returns calculated value
func (f functionStringListIntersect) Calculate(ctx *Context) (AttributeValue, error) {
	first, err := ctx.calculateListOfStringExpression(f.first)
	if err != nil {
		return UndefinedValue, bindError(bindError(err, "first argument"), f.describe())
	}

	second, err := ctx.calculateListOfStringExpression(f.second)
	if err != nil {
		return UndefinedValue, bindError(bindError(err, "second argument"), f.describe())
	}

	// FIXME: room for improvement here...
	var res []string
	for _, f := range first {
		for _, s := range second {
			if f == s {
				res = append(res, f)
			}
		}
	}

	return MakeListOfStringsValue(res), nil
}

func functionStringListIntersectValidator(args []Expression) functionMaker {
	if len(args) != 2 || args[0].GetResultType() != TypeListOfStrings || args[1].GetResultType() != TypeListOfStrings {
		return nil
	}

	return makeFunctionStringListIntersectAlt
}

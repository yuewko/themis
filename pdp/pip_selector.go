package pdp

import (
	"context"
	"fmt"

	ps "github.com/infobloxopen/themis/pip-service"
	"google.golang.org/grpc"
)

const pIPServicePort = ":5356"

type PIPSelector struct {
	service     string
	serviceAddr string
	queryType   string
	path        []Expression
	t           int
}

func MakePIPSelector(service, queryType string, path []Expression, t int) (PIPSelector, error) {
	addr, err := ps.GetPIPConnection(service)
	if err != nil {
		return PIPSelector{}, err
	}

	return PIPSelector{
		service:     service,
		serviceAddr: addr,
		queryType:   queryType,
		path:        path,
		t:           t}, nil
}

func (s PIPSelector) GetResultType() int {
	return s.t
}

func (s PIPSelector) describe() string {
	return fmt.Sprintf("PIPselector(%s.%s)", s.service, s.queryType)
}

func (s PIPSelector) getAttributeValue(ctx *Context) (AttributeValue, error) {
	attrList := []*ps.Attribute{}
	for _, p := range s.path {
		val, err := p.calculate(ctx)
		if err != nil {
			return undefinedValue, bindError(err, s.describe())
		}
		id := p.(AttributeDesignator).a.id
		t := p.(AttributeDesignator).a.t
		serializedVal, err := val.Serialize()
		if err != nil {
			return undefinedValue, bindError(err, s.describe())
		}
		attr := ps.Attribute{Id: id, Type: int32(t), Value: serializedVal}
		attrList = append(attrList, &attr)
	}

	conn, err := grpc.Dial(s.serviceAddr, grpc.WithInsecure())
	if err != nil {
		return undefinedValue, bindError(fmt.Errorf("cannot connect to pip, err: '%s'", err), s.describe())
	}
	defer conn.Close()

	c := ps.NewPIPClient(conn)
	request := &ps.Request{QueryType: s.queryType, Attributes: attrList}

	r, err := c.GetAttribute(context.Background(), request)
	if err != nil {
		return undefinedValue, bindError(err, s.describe())
	}
	if r.Status != ps.Response_OK {
		return undefinedValue, bindError(fmt.Errorf("Unexpected response status '%s'", r.Status), s.describe())
	}
	res := r.GetValues()
	val := res[0].GetValue()

	fmt.Printf("PIP returned value='%v'\n", val)

	return MakeStringValue(val), nil
}

func (s PIPSelector) calculate(ctx *Context) (AttributeValue, error) {
	return s.getAttributeValue(ctx)
}
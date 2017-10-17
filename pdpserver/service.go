package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
	ps "github.com/infobloxopen/themis/pip-service"
)

func makeEffect(effect int) (pb.Response_Effect, error) {
	switch effect {
	case pdp.EffectDeny:
		return pb.Response_DENY, nil

	case pdp.EffectPermit:
		return pb.Response_PERMIT, nil

	case pdp.EffectNotApplicable:
		return pb.Response_NOTAPPLICABLE, nil

	case pdp.EffectIndeterminate:
		return pb.Response_INDETERMINATE, nil

	case pdp.EffectIndeterminateD:
		return pb.Response_INDETERMINATED, nil

	case pdp.EffectIndeterminateP:
		return pb.Response_INDETERMINATEP, nil

	case pdp.EffectIndeterminateDP:
		return pb.Response_INDETERMINATEDP, nil
	}

	return pb.Response_INDETERMINATE, newUnknownEffectError(effect)
}

func makeFailEffect(effect pb.Response_Effect) (pb.Response_Effect, error) {
	switch effect {
	case pb.Response_DENY:
		return pb.Response_INDETERMINATED, nil

	case pb.Response_PERMIT:
		return pb.Response_INDETERMINATEP, nil

	case pb.Response_NOTAPPLICABLE, pb.Response_INDETERMINATE, pb.Response_INDETERMINATED, pb.Response_INDETERMINATEP, pb.Response_INDETERMINATEDP:
		return effect, nil
	}

	return pb.Response_INDETERMINATE, newUnknownEffectError(int(effect))
}

func (s *Server) newContext(c *pdp.LocalContentStorage, m *ps.ConnectionManager, in *pb.Request) (*pdp.Context, error) {
	ctx, err := pdp.NewContext(c, m, len(in.Attributes), func(i int) (string, pdp.AttributeValue, error) {
		a := in.Attributes[i]

		t, ok := pdp.TypeIDs[strings.ToLower(a.Type)]
		if !ok {
			return "", pdp.AttributeValue{}, bindError(newUnknownAttributeTypeError(a.Type), a.Id)
		}

		v, err := pdp.MakeValueFromString(t, a.Value)
		if err != nil {
			return "", pdp.AttributeValue{}, bindError(err, a.Id)
		}

		return a.Id, v, nil
	})
	if err != nil {
		return nil, newContextCreationError(err)
	}

	return ctx, nil
}

func (s *Server) newAttributes(obligations []pdp.AttributeAssignmentExpression, ctx *pdp.Context) ([]*pb.Attribute, error) {
	attrs := make([]*pb.Attribute, len(obligations))
	for i, e := range obligations {
		ID, t, s, err := e.Serialize(ctx)
		if err != nil {
			return attrs[:i], err
		}

		attrs[i] = &pb.Attribute{
			Id:    ID,
			Type:  t,
			Value: s}
	}

	return attrs, nil
}

func (s *Server) rawValidate(p *pdp.PolicyStorage, c *pdp.LocalContentStorage, m *ps.ConnectionManager, in *pb.Request) (pb.Response_Effect, []error, []*pb.Attribute) {
	ctx, err := s.newContext(c, m, in)
	if err != nil {
		return pb.Response_INDETERMINATE, []error{err}, nil
	}

	errs := []error{}

	r := p.Root().Calculate(ctx)
	effect, obligations, err := r.Status()
	if err != nil {
		errs = append(errs, newPolicyCalculationError(err))
	}

	re, err := makeEffect(effect)
	if err != nil {
		errs = append(errs, newEffectTranslationError(err))
	}

	if len(errs) > 0 {
		re, err = makeFailEffect(re)
		if err != nil {
			errs = append(errs, newEffectCombiningError(err))
		}
	}

	attrs, err := s.newAttributes(obligations, ctx)
	if err != nil {
		errs = append(errs, newObligationTranslationError(err))
	}

	return re, errs, attrs
}

func (s *Server) Validate(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Info("Validating context")

	s.RLock()
	p := s.p
	c := s.c
	m := s.pcm
	s.RUnlock()

	effect, errs, attrs := s.rawValidate(p, c, m, in)

	status := "Ok"
	if len(errs) > 1 {
		status = newMultiError(errs).Error()
	} else if len(errs) > 0 {
		status = errs[0].Error()
	}

	return &pb.Response{
		Effect:     effect,
		Reason:     status,
		Obligation: attrs}, nil
}

package jcon

import (
	"fmt"
	"net"

	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/infobloxopen/go-trees/iptree"
	"github.com/infobloxopen/go-trees/strtree"
	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func (c *contentItem) ppMap(v interface{}, keyIdx int) (interface{}, error) {
	pairs, ok := v.([]jparser.Pair)
	if !ok {
		return nil, newInvalidMapContentItemNodeError(v, fmt.Sprintf("level %d map", keyIdx+1))
	}

	m, err := newTypedMap(c, keyIdx)
	if err != nil {
		return nil, err
	}

	for _, p := range pairs {
		err = m.postProcess(p)
		if err != nil {
			return nil, err
		}
	}

	return m.get(), nil
}

func ppStringSequenceFromPairs(v []jparser.Pair, desc string, f func(s string) error) error {
	for i, p := range v {
		err := f(p.K)
		if err != nil {
			return bindErrorf(err, "%d", i+1)
		}
	}

	return nil
}

func ppStringSequenceFromArray(v []interface{}, desc string, f func(s string) error) error {
	for i, item := range v {
		s, ok := item.(string)
		if !ok {
			return bindErrorf(newStringCastError(item, desc), "%d", i+1)
		}

		err := f(s)
		if err != nil {
			return bindErrorf(err, "%d", i+1)
		}
	}

	return nil
}

func ppStringSequence(v interface{}, desc string, f func(s string) error) error {
	switch v := v.(type) {
	case []jparser.Pair:
		return ppStringSequenceFromPairs(v, desc, f)

	case []interface{}:
		return ppStringSequenceFromArray(v, desc, f)
	}

	return newInvalidSequenceContentItemNodeError(v, desc)
}

func (c *contentItem) ppValue(v interface{}) (interface{}, error) {
	switch c.t {
	case pdp.TypeBoolean:
		b, ok := v.(bool)
		if !ok {
			return nil, newBooleanCastError(v, "value")
		}

		return b, nil

	case pdp.TypeString:
		s, ok := v.(string)
		if !ok {
			return nil, newStringCastError(v, "value")
		}

		return s, nil

	case pdp.TypeInteger:
		x, ok := v.(float64)
		if !ok {
			return nil, newNumberCastError(v, "value")
		}

		if x < -9007199254740992 || x > 9007199254740992 {
			return nil, newIntegerOverflowError(x)
		}

		return int64(x), nil

	case pdp.TypeFloat:
		x, ok := v.(float64)
		if !ok {
			return nil, newNumberCastError(v, "value")
		}

		return x, nil

	case pdp.TypeAddress:
		s, ok := v.(string)
		if !ok {
			return nil, newStringCastError(v, "address value")
		}

		a := net.ParseIP(s)
		if a == nil {
			return nil, newAddressCastError(s)
		}

		return a, nil

	case pdp.TypeNetwork:
		s, ok := v.(string)
		if !ok {
			return nil, newStringCastError(v, "network value")
		}

		_, n, err := net.ParseCIDR(s)
		if err != nil {
			return nil, newNetworkCastError(s, err)
		}

		return n, nil

	case pdp.TypeDomain:
		s, ok := v.(string)
		if !ok {
			return nil, newStringCastError(v, "domain value")
		}

		d, err := pdp.AdjustDomainName(s)
		if err != nil {
			return nil, newDomainCastError(s, err)
		}

		return d, nil

	case pdp.TypeSetOfStrings:
		m := strtree.NewTree()
		i := 0
		err := ppStringSequence(v, "set of strings value", func(s string) error {
			if _, ok := m.Get(s); !ok {
				m.InplaceInsert(s, i)
				i++
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		return m, nil

	case pdp.TypeSetOfNetworks:
		m := iptree.NewTree()
		err := ppStringSequence(v, "set of networks value", func(s string) error {
			a := net.ParseIP(s)
			if a != nil {
				m.InplaceInsertIP(a, nil)
			} else {
				_, n, err := net.ParseCIDR(s)
				if err != nil {
					return newAddressNetworkCastError(s, err)
				}

				m.InplaceInsertNet(n, nil)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		return m, nil

	case pdp.TypeSetOfDomains:
		m := &domaintree.Node{}
		err := ppStringSequence(v, "set of domains value", func(s string) error {
			d, err := pdp.AdjustDomainName(s)
			if err != nil {
				return newDomainCastError(s, err)
			}

			m.InplaceInsert(d, nil)

			return nil
		})
		if err != nil {
			return nil, err
		}

		return m, nil

	case pdp.TypeListOfStrings:
		lst := []string{}
		err := ppStringSequence(v, "list of strings value", func(s string) error {
			lst = append(lst, s)
			return nil
		})
		if err != nil {
			return nil, err
		}

		return lst, nil
	}

	return nil, newInvalidContentItemTypeError(c.t)
}

func (c *contentItem) postProcess(v interface{}, keyIdx int) (interface{}, error) {
	if len(c.k) > keyIdx {
		return c.ppMap(v, keyIdx)
	}

	return c.ppValue(v)
}

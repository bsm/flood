package qfy

import (
	"encoding/json"
	"fmt"
	"hash/crc64"

	"github.com/bsm/intset"
)

var crcTable = crc64.MakeTable(crc64.ECMA)

// Fact is an interface of a fact that may be passed to a qualifier. Each fact must implement
// a Get(string) method which receives the attribute name and must return either a string or
// an int slice, depending on the attribute definition.
type Fact interface {
	Get(string) interface{}
}

// --------------------------------------------------------------------

type AttrType uint8

const (
	TypeUnknown AttrType = iota
	TypeStringSlice
	TypeIntSlice
)

// Attribute defines a qualifiable fact attribute;
// must have a name and a type
type Attribute struct {
	Name string
	Type AttrType
}

// --------------------------------------------------------------------

// RuleDef contains a JSON-parseable rule definition
type RuleDef struct {
	Attr string          `json:"attr"`   // the attribute name
	Op   string          `json:"op"`     // the operation code, either '+' or '-'
	Vals json.RawMessage `json:"values"` // the rule values, can be either an array of strings or ints
}

// DetectType attempts to detect the type of the values
func (r *RuleDef) DetectType() (AttrType, error) {
	var nums []int
	var strs []string

	if json.Unmarshal(r.Vals, &nums) == nil {
		return TypeIntSlice, nil
	} else if json.Unmarshal(r.Vals, &strs) == nil {
		return TypeStringSlice, nil
	}
	return TypeUnknown, r.invalid()
}

func (r *RuleDef) invalid() error {
	return fmt.Errorf("qfy: invalid rule: %s %s%v", r.Attr, r.Op, r.Vals)
}

func (r *RuleDef) toRule(kind AttrType, dict strDict) (Rule, error) {
	var vals []int

	switch kind {
	case TypeStringSlice:
		var strs []string
		if json.Unmarshal(r.Vals, &strs) != nil {
			return nil, r.invalid()
		}
		vals = dict.FetchSlice(strs...)
	case TypeIntSlice:
		if json.Unmarshal(r.Vals, &vals) != nil {
			return nil, r.invalid()
		}
	default:
		return nil, r.invalid()
	}

	if len(vals) == 0 {
		return nil, r.invalid()
	}

	switch r.Op {
	case "+":
		return newPlusRule(vals), nil
	case "-":
		return newMinusRule(vals), nil
	}
	return nil, r.invalid()
}

// --------------------------------------------------------------------

type converter interface {
	convert(interface{}) *intset.Set
}

// the qualification lookup process abstraction
type lookup struct {
	results   []int
	converter converter
	ruleCache map[uint64]bool
	factCache map[string]*intset.Set
}

func newLookup(cvt converter) *lookup {
	return &lookup{
		results:   make([]int, 0, 100),
		converter: cvt,
		ruleCache: make(map[uint64]bool, 1000),
		factCache: make(map[string]*intset.Set, 20),
	}
}

func (l *lookup) Clear() {
	l.results = l.results[:0]
	for k, _ := range l.ruleCache {
		delete(l.ruleCache, k)
	}
	for k, _ := range l.factCache {
		delete(l.factCache, k)
	}
}

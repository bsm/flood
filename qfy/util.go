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

// --------------------------------------------------------------------

// RuleDef contains a JSON-parseable rule definition
type RuleDef struct {
	Attr string      // the attribute name
	Op   string      // the operation code, either '+' or '-'
	Vals interface{} // the rule values, can be either an array of strings or ints
}

// UnmarshalJSON decodes JSON
func (r *RuleDef) UnmarshalJSON(data []byte) error {
	var temp *ruleDef
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	rdef := RuleDef{Attr: temp.Attr, Op: temp.Op, Vals: nil}
	if rdef.Vals, err = temp.DecodeVals(); err != nil {
		return err
	}

	*r = rdef
	_, err = r.DetectType()
	return err
}

// DetectType attempts to detect the type of the values
func (r *RuleDef) DetectType() (AttrType, error) {
	switch r.Vals.(type) {
	case []string, string:
		return TypeStringSlice, nil
	case []int, int:
		return TypeIntSlice, nil
	}
	return TypeUnknown, r.invalid()
}

func (r *RuleDef) invalid() error {
	return fmt.Errorf("qfy: invalid rule: %s %s%v", r.Attr, r.Op, r.Vals)
}

func (r *RuleDef) toRule(dict strDict) (Rule, error) {
	var vals []int

	switch vv := r.Vals.(type) {
	case []string:
		vals = dict.FetchSlice(vv...)
	case []int:
		vals = vv
	case string:
		vals = dict.FetchSlice(vv)
	case int:
		vals = []int{vv}
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

type ruleDef struct {
	Attr string          `json:"attr"`
	Op   string          `json:"op"`
	Vals json.RawMessage `json:"values"`
}

func (r *ruleDef) DecodeVals() (interface{}, error) {
	raw := r.Vals
	if len(raw) == 0 {
		return nil, nil
	}

	switch raw[0] {
	case '[': // slice
		for _, c := range raw[1:] {

			switch c {
			case '"': // string slice
				var vals []string
				err := json.Unmarshal(raw, &vals)
				return vals, err
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // int slice
				var vals []int
				err := json.Unmarshal(raw, &vals)
				return vals, err
			}
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // int
		var val int
		err := json.Unmarshal(raw, &val)
		return val, err
	case '"': // string
		var val string
		err := json.Unmarshal(raw, &val)
		return val, err
	}
	return nil, nil
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

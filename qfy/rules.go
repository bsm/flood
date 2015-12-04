package qfy

import (
	"fmt"
	"strings"
)

// Rule is an abstract logic evaluation
type Rule interface {
	// String returns a human-readable description
	String() string

	// Must return a globally unique and consistent conditions identifier.
	// For performance reasons, this value must be pre-calculated in the
	// constructor.
	crc64() uint64

	perform(fact Fact, state *State) bool
}

func rulesToString(rules []Rule, sep string) string {
	parts := make([]string, len(rules))
	for i, rule := range rules {
		parts[i] = rule.String()
	}
	return fmt.Sprintf("( %s )", strings.Join(parts, sep))
}

// --------------------------------------------------------------------

// A Rule that performs conditions against a fact attribute
type factCheck struct {
	hash uint64
	key  FactKey
	cond Condition
}

// CheckFact constructs a new rule, it accepts a fact key (used to query the fact)
// and an evaluation condition
func CheckFact(key FactKey, cond Condition) Rule {
	hash := NewCRC64('=', 2)
	hash.Add(uint64(key), cond.CRC64())

	return &factCheck{
		hash: hash.Sum64(),
		key:  key,
		cond: cond,
	}
}

// String returns a human-readable description
func (r *factCheck) String() string {
	return fmt.Sprintf("[%d]%s", r.key, r.cond.String())
}

func (r *factCheck) crc64() uint64 { return r.hash }
func (r *factCheck) perform(fact Fact, state *State) bool {
	if match, ok := state.rules[r.hash]; ok {
		return match
	}

	v, ok := state.facts[r.key]
	if !ok {
		switch vv := fact.GetQualifiable(r.key).(type) {
		case int:
			v = int64(vv)
		case int8:
			v = int64(vv)
		case int16:
			v = int64(vv)
		case int32:
			v = int64(vv)
		case int64:
			v = vv
		case uint:
			v = int64(vv)
		case uint8:
			v = int64(vv)
		case uint16:
			v = int64(vv)
		case uint32:
			v = int64(vv)
		case uint64:
			v = int64(vv)
		case float64:
			v = vv
		case []int:
			v = ints64FromInts(vv)
		case []int8:
			v = ints64FromInts8(vv)
		case []int16:
			v = ints64FromInts16(vv)
		case []int32:
			v = ints64FromInts32(vv)
		case []int64:
			v = SortInts64(vv...)
		case []uint:
			v = ints64FromUints(vv)
		case []uint8:
			v = ints64FromUints8(vv)
		case []uint16:
			v = ints64FromUints16(vv)
		case []uint32:
			v = ints64FromUints32(vv)
		case []uint64:
			v = ints64FromUints64(vv)
		default:
			state.rules[r.crc64()] = false
			return false
		}
		state.facts[r.key] = v
	}

	match := r.cond.Match(v)
	state.rules[r.crc64()] = match
	return match
}

// --------------------------------------------------------------------

// Combines two or more rules into by creating a logical AND
type conjunction struct {
	hash  uint64
	rules []Rule
}

// All requires all of the rules to match
func All(rules ...Rule) Rule {
	return &conjunction{hash: crc64FromRules('&', rules...), rules: rules}
}

// String returns a human-readable description
func (r *conjunction) String() string { return rulesToString(r.rules, " && ") }

func (r *conjunction) crc64() uint64 { return r.hash }
func (r *conjunction) perform(fact Fact, state *State) bool {
	if len(r.rules) == 0 {
		return false
	}
	for _, rule := range r.rules {
		if match, ok := state.rules[rule.crc64()]; ok && !match {
			return false
		}
	}
	for _, rule := range r.rules {
		if !rule.perform(fact, state) {
			return false
		}
	}
	return true
}

// --------------------------------------------------------------------

// disjunction combines two or more rules into by creating a logical OR
type disjunction struct {
	hash  uint64
	rules []Rule
}

// Any requires any of the rules to match
func Any(rules ...Rule) Rule {
	return &disjunction{hash: crc64FromRules('|', rules...), rules: rules}
}

// String returns a human-readable description
func (r *disjunction) String() string { return rulesToString(r.rules, " || ") }

func (r *disjunction) crc64() uint64 { return r.hash }
func (r *disjunction) perform(fact Fact, state *State) bool {
	for _, rule := range r.rules {
		if match, ok := state.rules[rule.crc64()]; ok && match {
			return true
		}
	}
	for _, rule := range r.rules {
		if rule.perform(fact, state) {
			return true
		}
	}
	return false
}

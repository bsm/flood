package qfy

import "fmt"

// Condition is an abstract logic evaluation condition
type Condition interface {
	// CRC64 must return a globally unique an consistent rules identifier
	CRC64() uint64

	// Match tests if the rule is qualified. Certain conditions can only
	// support certain types and must return a negative response if an
	// unsupported type is given.
	Match(interface{}) bool

	// String returns a human-readable rule description
	String() string
}

// --------------------------------------------------------------------

// Equality conditions require the fact value to match input.
type Equality struct{ val interface{} }

// EqualTo constructs an Equality condition
// Supports bool, string, intN, uintN and floatN fact values as inputs.
func EqualTo(v interface{}) *Equality { return &Equality{v} }

// Match tests if the condition is qualified
func (r *Equality) Match(v interface{}) bool { return v == r.val }

// String returns a human-readable description
func (r *Equality) String() string { return fmt.Sprintf("=%v", r.val) }

// CRC64 returns a unique ID
func (r *Equality) CRC64() uint64 { return crc64FromValue('=', r.val) }

// --------------------------------------------------------------------

// Inclusion conditions require at least one value to be present in the fact
// Supports only int64 and []int64 fact values as inputs.
type Inclusion struct{ vals Ints64 }

// OneOf constructs an Inclusion
func OneOf(vals []int64) *Inclusion { return &Inclusion{SortInts64(vals...)} }

// Match tests if the condition is qualified
func (r *Inclusion) Match(v interface{}) bool {
	switch vv := v.(type) {
	case int64:
		return r.vals.Exists(vv)
	case Ints64:
		return vv != nil && r.vals.Inter(vv)
	}
	return false
}

// String returns a human-readable description
func (r *Inclusion) String() string { return fmt.Sprintf("+%v", r.vals) }

// CRC64 returns a unique ID
func (r *Inclusion) CRC64() uint64 { return r.vals.crc64('+') }

// --------------------------------------------------------------------

// Exclusion conditions require none of the values to be present in the fact
type Exclusion struct{ vals Ints64 }

// NoneOf constructs an Exclusion
func NoneOf(vals []int64) *Exclusion { return &Exclusion{SortInts64(vals...)} }

// Match tests if the condition is qualified
// Supports only int64 and []int64 fact values as inputs.
func (r *Exclusion) Match(v interface{}) bool {
	switch vv := v.(type) {
	case int64:
		return !r.vals.Exists(vv)
	case Ints64:
		return !(vv != nil && r.vals.Inter(vv))
	}
	return true
}

// String returns a human-readable description
func (r *Exclusion) String() string { return fmt.Sprintf("-%v", r.vals) }

// CRC64 returns a unique ID
func (r *Exclusion) CRC64() uint64 { return r.vals.crc64('-') }

// --------------------------------------------------------------------

// Negation inverts a Condition
type Negation struct{ cond Condition }

// Not creates a negation
func Not(cond Condition) *Negation {
	return &Negation{cond: cond}
}

// CRC64 returns a combined unique ID
func (r *Negation) CRC64() uint64 { return crc64FromConditions('!', r.cond) }

// Match tests if the condition is qualified
func (r *Negation) Match(v interface{}) bool { return !r.cond.Match(v) }

// String returns a human-readable description
func (r *Negation) String() string { return fmt.Sprintf("!%s", r.cond.String()) }

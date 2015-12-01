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

// NumericGreaterOrEqual conditions require the fact value to match input.
type NumericGreaterOrEqual struct{ val float64 }

// GreaterOrEqual constructs a NumericGreaterOrEqual condition
// Supports intN and floatN fact values as inputs.
func GreaterOrEqual(v float64) *NumericGreaterOrEqual { return &NumericGreaterOrEqual{v} }

// Match tests if the condition is qualified
func (r *NumericGreaterOrEqual) Match(v interface{}) bool {
	switch vv := v.(type) {
	case int64:
		return float64(vv) >= r.val
	case float64:
		return vv >= r.val
	}
	return false
}

// String returns a human-readable description
func (r *NumericGreaterOrEqual) String() string { return fmt.Sprintf(">=%v", r.val) }

// CRC64 returns a unique ID
func (r *NumericGreaterOrEqual) CRC64() uint64 { return crc64FromValue('(', r.val) }

// --------------------------------------------------------------------

// NumericLessOrEqual conditions require the fact value to match input.
type NumericLessOrEqual struct{ val float64 }

// LessOrEqual constructs a NumericLessOrEqual condition
// Supports intN and floatN fact values as inputs.
func LessOrEqual(v float64) *NumericLessOrEqual { return &NumericLessOrEqual{v} }

// Match tests if the condition is qualified
func (r *NumericLessOrEqual) Match(v interface{}) bool {
	switch vv := v.(type) {
	case int64:
		return float64(vv) <= r.val
	case float64:
		return vv <= r.val
	}
	return false
}

// String returns a human-readable description
func (r *NumericLessOrEqual) String() string { return fmt.Sprintf("<=%v", r.val) }

// CRC64 returns a unique ID
func (r *NumericLessOrEqual) CRC64() uint64 { return crc64FromValue(')', r.val) }

// --------------------------------------------------------------------

// NumericGreater conditions require the fact value to match input.
type NumericGreater struct{ val float64 }

// GreaterThan constructs a NumericGreater condition
// Supports intN and floatN fact values as inputs.
func GreaterThan(v float64) *NumericGreater { return &NumericGreater{v} }

// Match tests if the condition is qualified
func (r *NumericGreater) Match(v interface{}) bool {
	switch vv := v.(type) {
	case int64:
		return float64(vv) > r.val
	case float64:
		return vv > r.val
	}
	return false
}

// String returns a human-readable description
func (r *NumericGreater) String() string { return fmt.Sprintf(">%v", r.val) }

// CRC64 returns a unique ID
func (r *NumericGreater) CRC64() uint64 { return crc64FromValue('>', r.val) }

// --------------------------------------------------------------------

// NumericLess conditions require the fact value to match input.
type NumericLess struct{ val float64 }

// LessThan constructs a NumericLess condition
// Supports intN and floatN fact values as inputs.
func LessThan(v float64) *NumericLess { return &NumericLess{v} }

// Match tests if the condition is qualified
func (r *NumericLess) Match(v interface{}) bool {
	switch vv := v.(type) {
	case int64:
		return float64(vv) < r.val
	case float64:
		return vv < r.val
	}
	return false
}

// String returns a human-readable description
func (r *NumericLess) String() string { return fmt.Sprintf("<%v", r.val) }

// CRC64 returns a unique ID
func (r *NumericLess) CRC64() uint64 { return crc64FromValue('<', r.val) }

// --------------------------------------------------------------------

// NumericRange conditions require the fact value to match input.
type NumericRange struct{ min, max float64 }

// Between constructs a NumericRange condition
// Supports intN and floatN fact values as inputs.
func Between(min, max float64) *NumericRange { return &NumericRange{min, max} }

// Match tests if the condition is qualified
func (r *NumericRange) Match(v interface{}) bool {
	switch vv := v.(type) {
	case int64:
		fv := float64(vv)
		return fv >= r.min && fv <= r.max
	case float64:
		return vv >= r.min && vv <= r.max
	}
	return false
}

// String returns a human-readable description
func (r *NumericRange) String() string { return fmt.Sprintf("%v..%v", r.min, r.max) }

// CRC64 returns a unique ID
func (r *NumericRange) CRC64() uint64 { return crc64FromValue('~', r.min, r.max) }

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

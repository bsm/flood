package qfy

import (
	"fmt"

	"github.com/bsm/intset"
)

// Condition is an abstract logic evaluation condition
type Condition interface {
	// CRC64 must return a globally unique an consistent rules identifier
	CRC64() uint64

	// Match tests if the rule is qualified
	Match(*intset.Set) bool

	// String returns a human-readable rule description
	String() string
}

// --------------------------------------------------------------------

// Inclusion conditions require at least one value to be present in the fact
type Inclusion struct{ vals *intset.Set }

// OneOf constructs an Inclusion
func OneOf(vals []int) *Inclusion { return &Inclusion{intset.Use(vals...)} }

// Match tests if the condition is qualified
func (r *Inclusion) Match(vv *intset.Set) bool { return vv != nil && r.vals.Intersects(vv) }

// String returns a human-readable description
func (r *Inclusion) String() string { return fmt.Sprintf("+%v", r.vals.Slice()) }

// CRC64 returns a unique ID
func (r *Inclusion) CRC64() uint64 { return crc64FromInts('+', r.vals.Slice()) }

// --------------------------------------------------------------------

// Exclusion conditions require none of the values to be present in the fact
type Exclusion struct{ vals *intset.Set }

// NoneOf constructs an Exclusion
func NoneOf(vals []int) *Exclusion { return &Exclusion{intset.Use(vals...)} }

// Match tests if the condition is qualified
func (r *Exclusion) Match(vv *intset.Set) bool { return !(vv != nil && r.vals.Intersects(vv)) }

// String returns a human-readable description
func (r *Exclusion) String() string { return fmt.Sprintf("-%v", r.vals.Slice()) }

// CRC64 returns a unique ID
func (r *Exclusion) CRC64() uint64 { return crc64FromInts('-', r.vals.Slice()) }

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
func (r *Negation) Match(vals *intset.Set) bool { return !r.cond.Match(vals) }

// String returns a human-readable description
func (r *Negation) String() string { return fmt.Sprintf("!%s", r.cond.String()) }

package qfy

import "github.com/bsm/intset"

// Rule is an abstract logic evaluation.
type Rule interface {
	// UID uniquely identifies the rule
	UID() uint64

	// Match tests if the rule is qulified
	Match(*intset.Set) bool

	// String returns a human-readable rule description
	String() string
}

// --------------------------------------------------------------------

type baseRule struct {
	hash uint64
	vals *intset.Set
}

func newBaseRule(sign byte, vals []int) *baseRule {
	vset := intset.Use(vals...)
	hash := newCRCHash(sign, len(vals))
	for _, val := range vals {
		hash.Add(uint64(val))
	}
	return &baseRule{hash: hash.Sum64(), vals: vset}
}

// UID returns a unique rule identifier
func (r *baseRule) UID() uint64 { return r.hash }

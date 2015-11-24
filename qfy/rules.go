package qfy

import (
	"fmt"
	"strings"

	"github.com/bsm/intset"
)

// --------------------------------------------------------------------

// Inclusion rules require either of values to be present in the fact
type Inclusion struct{ baseRule }

// OneOf constructs an Inclusion
func OneOf(vals []int) *Inclusion { return &Inclusion{*newBaseRule('+', vals)} }

// Match tests if the rule is qualified
func (r *Inclusion) Match(vv *intset.Set) bool { return vv != nil && r.vals.Intersects(vv) }

// String returns a human-readable description
func (r *Inclusion) String() string { return fmt.Sprintf("+%v", r.vals.Slice()) }

// --------------------------------------------------------------------

// Exclusion rules require none of the values to be present in the fact
type Exclusion struct{ baseRule }

// NoneOf constructs an Exclusion
func NoneOf(vals []int) *Exclusion { return &Exclusion{*newBaseRule('-', vals)} }

// Match tests if the rule is qualified
func (r *Exclusion) Match(vv *intset.Set) bool { return !(vv != nil && r.vals.Intersects(vv)) }

// String returns a human-readable description
func (r *Exclusion) String() string { return fmt.Sprintf("-%v", r.vals.Slice()) }

// --------------------------------------------------------------------

// Negation inverts a Rule
type Negation struct {
	hash uint64
	rule Rule
}

// Not creates a negation
func Not(rule Rule) *Negation {
	hash := newCRCHash('N', 1)
	hash.Add(rule.UID())
	return &Negation{rule: rule, hash: hash.Sum64()}
}

// UID returns a combined unique ID
func (r *Negation) UID() uint64 { return r.hash }

// Match tests if the rule is qualified
func (r *Negation) Match(vals *intset.Set) bool { return !r.rule.Match(vals) }

// String returns a human-readable description
func (r *Negation) String() string { return fmt.Sprintf("( NOT %s )", r.rule.String()) }

// --------------------------------------------------------------------

// Conjunction combines two or more rules into by creating a logical AND
type Conjunction struct {
	hash  uint64
	rules []Rule
}

// All requires all of the rules to match
func All(rules ...Rule) *Conjunction {
	hash := newCRCHash('A', len(rules))
	for _, rule := range rules {
		hash.Add(rule.UID())
	}
	return &Conjunction{
		rules: rules,
		hash:  hash.Sum64(),
	}
}

// UID returns a combined unique ID
func (r *Conjunction) UID() uint64 { return r.hash }

// Match tests if the rule is qualified
func (r *Conjunction) Match(vals *intset.Set) bool {
	if len(r.rules) == 0 {
		return false
	}

	for _, rule := range r.rules {
		if !rule.Match(vals) {
			return false
		}
	}
	return true
}

// String returns a human-readable description
func (r *Conjunction) String() string {
	parts := make([]string, len(r.rules))
	for i, rule := range r.rules {
		parts[i] = rule.String()
	}
	return fmt.Sprintf("( %s )", strings.Join(parts, " AND "))
}

// --------------------------------------------------------------------

// Disjunction combines two or more rules into by creating a logical OR
type Disjunction struct {
	hash  uint64
	rules []Rule
}

// Any requires any of the rules to match
func Any(rules ...Rule) *Disjunction {
	hash := newCRCHash('O', len(rules))
	for _, rule := range rules {
		hash.Add(rule.UID())
	}
	return &Disjunction{
		rules: rules,
		hash:  hash.Sum64(),
	}
}

// UID returns a combined unique ID
func (r *Disjunction) UID() uint64 { return r.hash }

// Match tests if the rule is qualified
func (r *Disjunction) Match(vals *intset.Set) bool {
	for _, rule := range r.rules {
		if rule.Match(vals) {
			return true
		}
	}
	return false
}

// String returns a human-readable description
func (r *Disjunction) String() string {
	parts := make([]string, len(r.rules))
	for i, rule := range r.rules {
		parts[i] = rule.String()
	}
	return fmt.Sprintf("( %s )", strings.Join(parts, " OR "))
}

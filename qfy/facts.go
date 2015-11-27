package qfy

// Fact is an interface of a fact that may be passed to a qualifier. Each fact must implement
// a Get(string) method which receives the attribute name and must return either a string or
// an int slice, depending on the attribute definition.
type Fact interface {
	GetQualifiable(FactKey) []int
}

// FactKey is a unique identifier of any fact attribute
type FactKey uint16

// MustBe is syntactic sugar for CheckFact.
//
// Example:
// 	const (
//		AttrA qfy.FactCheck = iota
//		AttrB
//	)
//	ruleA := AttrA.MustBe(qfy.OneOf(1,2,3))
//	ruleB := AttrB.MustBe(qfy.OneOf(6,5,4))
//  ruleC := qfy.And(ruleA, ruleB)
func (k FactKey) MustBe(cond Condition) Rule { return CheckFact(k, cond) }

// MustInclude is an alias for MustBe
func (k FactKey) MustInclude(cond Condition) Rule { return k.MustBe(cond) }

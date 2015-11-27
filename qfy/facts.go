package qfy

// Fact is an interface of a fact that may be passed to a qualifier. Each fact must implement
// a GetQualifiable(FatKey) method which receives a key and must return either a bool, an int64
// a float64 or an []int64 slice.
type Fact interface {
	GetQualifiable(FactKey) interface{}
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

// MustNotBe is the equivalent of MustBe(Not(...))
func (k FactKey) MustNotBe(cond Condition) Rule { return k.MustBe(Not(cond)) }

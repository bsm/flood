package quest

import "fmt"

// Outcome references a single outcome if a quest match
type Outcome int64

// Fact is an individual fact that is presented to the Quest for matching possible outcomes.
type Fact interface {
	// GetFactValue receives a trait name and returns a value. The value type
	// must match the registered trait type.
	GetFactValue(string) interface{}
}

// Rule is a set of conditions.
type Rule struct {
	// The list of conditions associated with that rule. For a rule to match, ANY of the conditions
	// must be fulfilled
	Conditions []Condition
	// If true, the rule is negated
	Negation bool
}

// Condition is a single match criteria.
type Condition struct {
	Trait      string      // the trait name
	Comparator Comparator  // the comparator
	Value      interface{} // the value to compare against, must match the registered trait type
}

type Comparator uint8

func (c Comparator) String() string {
	switch c {
	case ComparatorEqual:
		return "="
	case ComparatorGreater:
		return ">"
	case ComparatorLess:
		return "<"
		// case ComparatorGreaterOrEqual:
		// 	return ">="
		// case ComparatorLessOrEqual:
		// 	return "<="
	}
	return "?"
}

const (
	ComparatorEqual Comparator = iota + 1
	ComparatorGreater
	ComparatorLess
	// ComparatorGreaterOrEqual
	// ComparatorLessOrEqual
)

type traitData struct {
	Data       matchData
	Exclusions []ruleReference
}

type matchData interface {
	Check(*Condition) error
	Store(ruleReference, interface{})
	Find(interface{}) ([]ruleReference, error)
}

type ruleReference struct {
	outcome Outcome
	index   int
}

// --------------------------------------------------------------------

// Quest can register rules and compare with facts to match outcomes
type Quest struct {
	outcomes checksums
	traits   map[string]traitData
}

// New creates a new, blank quest object
func New() *Quest {
	return &Quest{
		outcomes: make(checksums),
		traits:   make(map[string]traitData),
	}
}

// RegisterTrait registers a new fact trait.
func (q *Quest) RegisterTrait(name string, kind TraitKind) error {
	if _, ok := q.traits[name]; ok {
		return fmt.Errorf("quest: trait '%s' is already regitered", name)
	}

	switch kind {
	case StringHash:
		q.traits[name] = traitData{Data: make(stringHash)}
	case Int64Hash:
		q.traits[name] = traitData{Data: make(int64Hash)}
	case Int32Hash:
		q.traits[name] = traitData{Data: make(int32Hash)}
	case BoolHash:
		q.traits[name] = traitData{Data: make(boolHash)}
	default:
		return fmt.Errorf("quest: unknown trait kind %d", kind)
	}
	return nil
}

// AddRule registers a new rule requirement for a particular outcome. An outcome can only be matched
// if ALL the associated rules match the fact criteria.
func (q *Quest) AddRule(outcome Outcome, rule *Rule) error {

	// Check all conditions first
	for _, cond := range rule.Conditions {
		trait, ok := q.traits[cond.Trait]
		if !ok {
			return fmt.Errorf("quest: condition references unknown trait '%s'", cond.Trait)
		}
		if err := trait.Data.Check(&cond); err != nil {
			return err
		}
	}

	// Now store conditions and outcome
	index := q.outcomes[outcome].size
	rlref := ruleReference{outcome, index}
	for _, cond := range rule.Conditions {
		trait := q.traits[cond.Trait]
		trait.Data.Store(rlref, cond.Value)
		if rule.Negation {
			trait.Exclusions = append(trait.Exclusions, rlref)
		}
		q.traits[cond.Trait] = trait
	}
	q.outcomes.Mark(outcome, index, !rule.Negation)

	return nil
}

// Match matches a fact against known rules and collects outcomes
func (q *Quest) Match(fact Fact) ([]Outcome, error) {
	var matches checksums
	if cached := checksumsPool.Get(); cached != nil {
		matches = cached.(checksums)
		matches.Reset()
	} else {
		matches = make(checksums)
	}

	for name, trait := range q.traits {
		val := fact.GetFactValue(name)
		if val == nil {
			continue
		}

		refs, err := trait.Data.Find(val)
		if err != nil {
			return nil, err
		}

		for _, ref := range trait.Exclusions {
			matches.Mark(ref.outcome, ref.index, false)
		}
		for _, ref := range refs {
			matches.Mark(ref.outcome, ref.index, true)
		}
	}

	var res []Outcome
	for outcome, csum := range matches {
		if csum.size != 0 && q.outcomes[outcome].Equal(csum) {
			res = append(res, outcome)
		}
	}
	checksumsPool.Put(matches)
	return res, nil
}

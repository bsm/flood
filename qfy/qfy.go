package qfy

import (
	"sync"

	"github.com/bsm/intset"
)

// Qualifier represents a rule engine that can match a fact against a
// list of rule-sets using a fixed set of attributes
type Qualifier struct {
	registry []target

	statePool, intsPool sync.Pool
}

// New creates a new qualifier with a list of known/qualifiable attributes
func New() *Qualifier { return &Qualifier{} }

// Resolve registers a rule with a numeric id resolved by that rule
func (q *Qualifier) Resolve(rule Rule, id int) {
	q.registry = append(q.registry, target{rule, id})
}

// Select performs the qualification and matches all known rules against a given fact
// returning a list of associated identifiers
func (q *Qualifier) Select(fact Fact) []int {
	if fact == nil {
		return nil
	}

	ints := q.fetchInts()
	state := q.fetchState()

	// Scan registry
	for _, t := range q.registry {
		if t.rule.perform(fact, state) {
			ints = append(ints, t.id)
		}
	}

	res := make([]int, len(ints))
	copy(res, ints)
	q.intsPool.Put(ints)
	q.statePool.Put(state)
	return res
}

func (q *Qualifier) fetchState() *State {
	if c := q.statePool.Get(); c != nil {
		state := c.(*State)
		state.Reset()
		return state
	}
	return NewState()
}

func (q *Qualifier) fetchInts() []int {
	if c := q.intsPool.Get(); c != nil {
		ints := c.([]int)
		ints = ints[:0]
		return ints
	}
	return make([]int, 0, 10)
}

// --------------------------------------------------------------------

type target struct {
	rule Rule
	id   int
}

// State holds the state of the qualification process
type State struct {
	rules map[uint64]bool
	facts map[FactKey]*intset.Set
}

// NewState initiali
func NewState() *State {
	return &State{
		rules: make(map[uint64]bool, 100),
		facts: make(map[FactKey]*intset.Set, 10),
	}
}

// Reset resets the stateory state
func (m *State) Reset() {
	for k := range m.rules {
		delete(m.rules, k)
	}
	for k := range m.facts {
		delete(m.facts, k)
	}
}

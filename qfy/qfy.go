package qfy

import "sync"

// Qualifier represents a rule engine that can match a fact against a
// list of rule-sets using a fixed set of attributes
type Qualifier struct {
	registry []target

	statePool sync.Pool
}

// New creates a new qualifier with a list of known/qualifiable attributes
func New() *Qualifier { return &Qualifier{} }

// Resolve registers a rule with a numeric id resolved by that rule
func (q *Qualifier) Resolve(rule Rule, id int64) {
	q.registry = append(q.registry, target{rule, id})
}

// Select performs the qualification and matches all known rules against a given fact
// returning a list of associated identifiers
func (q *Qualifier) Select(fact Fact) []int64 {
	if fact == nil {
		return nil
	}

	state := q.fetchState()
	for _, t := range q.registry {
		if t.rule.perform(fact, state) {
			state.results = append(state.results, t.id)
		}
	}

	res := make([]int64, len(state.results))
	copy(res, state.results)
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

// --------------------------------------------------------------------

type target struct {
	rule Rule
	id   int64
}

// State holds the state of the qualification process
type State struct {
	results []int64
	rules   map[uint64]bool
	facts   map[FactKey]interface{}
}

// NewState initiali
func NewState() *State {
	return &State{
		rules: make(map[uint64]bool, 100),
		facts: make(map[FactKey]interface{}, 20),
	}
}

// Reset resets the stateory state
func (m *State) Reset() {
	m.results = m.results[:0]
	for k := range m.rules {
		delete(m.rules, k)
	}
	for k := range m.facts {
		delete(m.facts, k)
	}
}

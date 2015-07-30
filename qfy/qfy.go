package qfy

import (
	"fmt"
	"io"
	"sync"

	"github.com/bsm/intset"
)

// Qualifier represents a rule engine that can match a fact against a
// list of rule-sets using a fixed set of attributes
type Qualifier struct {
	attrs []string
	root  *rootNode
	dict  strDict

	cache *sync.Pool
}

// New creates a new qualifier with a list of known/qualifiable attributes
func New(attrs []string) *Qualifier {
	q := &Qualifier{
		attrs: attrs,
		root:  &rootNode{},
		dict:  make(strDict, 100),
	}
	q.cache = &sync.Pool{New: func() interface{} {
		return newLookup(q)
	}}
	return q
}

// Feed registers a new identifier with its associated rule-set
func (q *Qualifier) Feed(id int, ruleSet []RuleDef) error {
	var parent node = q.root
	var child node

	for _, attr := range q.attrs {
		var rules RuleSet

		for _, rdef := range ruleSet {
			if rdef.Attr == attr {
				rule, err := rdef.toRule(q.dict)
				if err != nil {
					return err
				}
				rules = append(rules, rule)
			}
		}

		if len(rules) == 0 {
			child = &passNode{}
		} else {
			child = newClauseNode(attr, rules)
		}
		parent = parent.Merge(child)
	}
	parent.Merge(valueNode{id})
	return nil
}

// Select performs the qualification and matches all known rules against a given fact
// returning a list of associated identifiers
func (q *Qualifier) Select(fact Fact) []int {
	acc := q.cache.Get().(*lookup)
	q.root.Walk(fact, acc)

	nums := make([]int, len(acc.results))
	copy(nums, acc.results)

	acc.Clear()
	q.cache.Put(acc)

	return nums
}

// Graph prints a graph to writer, pass nil to print to stdout
func (q *Qualifier) Graph(w io.Writer) {
	fmt.Fprintf(w, "digraph {\n")
	q.root.Graph(w, 0, "[root]")
	fmt.Fprintf(w, "}\n")
}

func (q *Qualifier) convert(vals interface{}) *intset.Set {
	switch vv := vals.(type) {
	case []string:
		return intset.Use(q.dict.GetSlice(vv...)...)
	case []int:
		return intset.Use(vv...)
	}
	return nil
}

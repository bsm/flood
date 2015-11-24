package qfy

import (
	"fmt"
	"io"
	"sync"
)

// Fact is an interface of a fact that may be passed to a qualifier. Each fact must implement
// a Get(string) method which receives the attribute name and must return either a string or
// an int slice, depending on the attribute definition.
type Fact interface {
	GetQualifiable(string) []int
}

// --------------------------------------------------------------------

// Qualifier represents a rule engine that can match a fact against a
// list of rule-sets using a fixed set of attributes
//
// Example setup:
//
// q := qfy.New([]string{"country", "months"})
// q.Feed(11, map[string]Rule{
//   "country": qfy.Include(q, []string{"GB"}),
//   "months":  qfy.Include(q, []int{2015, 2014}),
// })
// q.Feed(12, map[string]Rule{
//   "country": qfy.Exclude(q, []string{"US"}),
//   "months":  qfy.And(qfy.Between(2010, 2015), qfy.Exclude(q, []int{2012})),
// })
type Qualifier struct {
	attrs []string
	root  *rootNode
	cache sync.Pool
}

// New creates a new qualifier with a list of known/qualifiable attributes
func New(attrs []string) *Qualifier {
	return &Qualifier{attrs: attrs, root: &rootNode{}}
}

// Feed registers a new targetID with a set of rules by attribute name.
func (q *Qualifier) Feed(targetID int, rules map[string]Rule) {
	var parent node = q.root
	var child node

	for _, attr := range q.attrs {
		if rule, ok := rules[attr]; ok {
			child = newClauseNode(attr, rule)
		} else {
			child = &passNode{}
		}
		parent = parent.Merge(child)
	}
	parent.Merge(valueNode{targetID})
}

// Select performs the qualification and matches all known rules against a given fact
// returning a list of associated identifiers
func (q *Qualifier) Select(fact Fact) []int {
	acc := q.makeLookup()
	q.root.Walk(fact, acc)

	nums := make([]int, len(acc.results))
	copy(nums, acc.results)

	q.cache.Put(acc)
	return nums
}

// Graph prints a graph to writer, pass nil to print to stdout
func (q *Qualifier) Graph(w io.Writer) {
	fmt.Fprintf(w, "digraph {\n")
	q.root.Graph(w, 0, "[root]")
	fmt.Fprintf(w, "}\n")
}

func (q *Qualifier) makeLookup() *lookup {
	if c := q.cache.Get(); c != nil {
		lp := c.(*lookup)
		lp.Clear()
		return lp
	}
	return newLookup()
}

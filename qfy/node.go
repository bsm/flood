package qfy

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/bsm/intset"
)

type node interface {
	Walk(Fact, *lookup)
	Merge(node) node
	Len() int
	Graph(io.Writer, int64, string)
}

type rootNode struct{ baseNode }
type passNode struct{ baseNode }

// --------------------------------------------------------------------

type baseNode struct {
	children []node
}

func (n *baseNode) Merge(nNode node) node {
	switch tNode := nNode.(type) {
	case valueNode:
		for i, c := range n.children {
			if vals, ok := c.(valueNode); ok {
				vals = append(vals, tNode...)
				n.children[i] = vals
				return vals
			}
		}
	case *passNode:
		for _, c := range n.children {
			if _, ok := c.(*passNode); ok {
				return c
			}
		}
	case *clauseNode:
		uid := tNode.rule.UID()
		for _, c := range n.children {
			if clause, ok := c.(*clauseNode); ok && clause.rule.UID() == uid {
				return c
			}
		}
	}
	n.children = append(n.children, nNode)
	return nNode
}

func (n *baseNode) Len() int {
	size := len(n.children)
	for _, m := range n.children {
		size += m.Len()
	}
	return size
}

func (n *baseNode) Walk(fact Fact, acc *lookup) {
	for _, m := range n.children {
		m.Walk(fact, acc)
	}
}

func (n *baseNode) Graph(w io.Writer, id int64, label string) {
	fmt.Fprintf(w, "\tN%016x [label = \"%s\"]\n", id, label)
	for _, child := range n.children {
		var childLabel string

		switch cnode := child.(type) {
		case valueNode:
			childLabel = fmt.Sprintf("%+v", cnode)
		case *passNode:
			childLabel = "[pass]"
		case *clauseNode:
			childLabel = cnode.rule.String()
		}

		childID := rand.Int63()
		child.Graph(w, childID, childLabel)
		fmt.Fprintf(w, "\tN%016x -> N%016x\n", id, childID)
	}
}

// --------------------------------------------------------------------

type clauseNode struct {
	baseNode
	attr string
	rule Rule
}

func newClauseNode(attr string, rule Rule) *clauseNode {
	return &clauseNode{
		attr: attr,
		rule: rule,
	}
}

func (n *clauseNode) Walk(fact Fact, acc *lookup) {
	vals, ok := acc.factCache[n.attr]
	if !ok {
		vals = intset.Use(fact.GetQualifiable(n.attr)...)
		acc.factCache[n.attr] = vals
	}

	ruleID := n.rule.UID()
	match, ok := acc.ruleCache[ruleID]
	if !ok {
		match = n.rule.Match(vals)
		acc.ruleCache[ruleID] = match
	}

	if !match {
		return
	}
	n.baseNode.Walk(fact, acc)
}

// --------------------------------------------------------------------

type valueNode []int

func (n valueNode) Len() int          { return 0 }
func (n valueNode) Merge(_ node) node { return n }
func (n valueNode) Walk(_ Fact, acc *lookup) {
	acc.results = append(acc.results, n...)
}
func (n valueNode) Graph(w io.Writer, id int64, label string) {
	fmt.Fprintf(w, "\tN%016x [label = \"%s\"]\n", id, label)
}

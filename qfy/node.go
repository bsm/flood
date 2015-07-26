package qfy

import (
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"io"
	"math/rand"
	"sort"
	"strings"
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
		uid := tNode.UID()
		for _, c := range n.children {
			if clause, ok := c.(*clauseNode); ok && clause.UID() == uid {
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
			childLabel = cnode.String()
		}

		childID := rand.Int63()
		child.Graph(w, childID, childLabel)
		fmt.Fprintf(w, "\tN%016x -> N%016x\n", id, childID)
	}
}

// --------------------------------------------------------------------

type clauseNode struct {
	baseNode
	hash  uint64
	attr  string
	rules RuleSet
}

func newClauseNode(attr string, rules RuleSet) *clauseNode {
	sort.Sort(rules)

	spad := len(attr)
	salt := make([]byte, spad+len(rules)*8)
	copy(salt[0:], attr)
	for i, rule := range rules {
		binary.LittleEndian.PutUint64(salt[spad+i*8:], rule.UID())
	}

	return &clauseNode{
		hash:  crc64.Checksum(salt, crcTable),
		attr:  attr,
		rules: rules,
	}
}

func (n *clauseNode) Walk(fact Fact, acc *lookup) {
	vals, ok := acc.factCache[n.attr]
	if !ok {
		vals = acc.converter.convert(fact.Get(n.attr))
		acc.factCache[n.attr] = vals
	}

	for _, rule := range n.rules {
		ruleID := rule.UID()
		match, ok := acc.ruleCache[ruleID]
		if !ok {
			match = rule.Match(vals)
			acc.ruleCache[ruleID] = match
		}

		if !match {
			return
		}
	}
	n.baseNode.Walk(fact, acc)
}

func (n *clauseNode) UID() uint64 { return n.hash }
func (n *clauseNode) String() string {
	clauses := make([]string, len(n.rules))
	for i, rule := range n.rules {
		clauses[i] = rule.String()
	}
	return fmt.Sprintf("%s %s", n.attr, strings.Join(clauses, ","))
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

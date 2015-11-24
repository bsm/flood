package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("baseNode", func() {
	var subject *baseNode
	var _ node = subject

	BeforeEach(func() {
		subject = &baseNode{}
	})

	It("should merge value children", func() {
		n1 := subject.Merge(valueNode{1, 2, 3})
		Expect(subject.children).To(HaveLen(1))
		Expect(n1).To(Equal(valueNode{1, 2, 3}))
		n2 := subject.Merge(valueNode{4, 5})
		Expect(subject.children).To(HaveLen(1))
		Expect(n2).To(Equal(valueNode{1, 2, 3, 4, 5}))
	})

	It("should merge pass children", func() {
		subject.Merge(&passNode{})
		Expect(subject.children).To(HaveLen(1))
		subject.Merge(&passNode{})
		Expect(subject.children).To(HaveLen(1))
	})

	It("should merge clause children", func() {
		n1 := subject.Merge(&passNode{})
		Expect(subject.children).To(HaveLen(1))
		n2 := subject.Merge(&passNode{})
		Expect(subject.children).To(HaveLen(1))
		Expect(n2).To(Equal(n1))
	})

	It("should walk", func() {
		subject.Merge(valueNode{1, 2, 3})
		acc := newLookup()
		subject.Walk(nil, acc)
		Expect(acc.results).To(HaveLen(3))
	})

})

var _ = Describe("clauseNode", func() {
	var subject *clauseNode
	var _ node = subject

	BeforeEach(func() {
		subject = newClauseNode("x", All(
			OneOf([]int{1, 2}),
			OneOf([]int{3, 4}),
		))
		subject.Merge(valueNode{8})
	})

	It("should skip when no match", func() {
		acc := newLookup()
		subject.Walk(mockFact{"y": {53, 54, 55}}, acc)
		subject.Walk(mockFact{"x": {53, 54, 55}}, acc)
		Expect(acc.results).To(BeEmpty())
	})

	It("should skip when partial match", func() {
		acc := newLookup()
		subject.Walk(mockFact{"x": {1, 2, 55}}, acc)
		Expect(acc.results).To(BeEmpty())
	})

	It("should cache rules", func() {
		acc := newLookup()
		subject.Walk(mockFact{"x": {1, 3}}, acc)
		Expect(acc.ruleCache).To(HaveLen(1))
	})

	It("should cache fact values", func() {
		acc := newLookup()
		subject.Walk(mockFact{"x": {1, 3}}, acc)
		Expect(acc.factCache).To(HaveLen(1))
		Expect(acc.factCache).To(HaveKey("x"))
		Expect(acc.factCache["x"].Len()).To(Equal(2))
	})

	It("should walk when full match", func() {
		acc := newLookup()
		subject.Walk(mockFact{"x": {1, 3}}, acc)
		Expect(acc.results).To(HaveLen(1))
		subject.Walk(mockFact{"x": {2, 3}}, acc)
		Expect(acc.results).To(HaveLen(2))
		subject.Walk(mockFact{"x": {1, 53, 4}}, acc)
		Expect(acc.results).To(HaveLen(3))
	})

})

var _ = Describe("valueNode", func() {
	var subject valueNode
	var _ node = subject

	BeforeEach(func() {
		subject = valueNode{}
	})

	It("should walk", func() {
		acc := newLookup()
		valueNode{1, 2, 3}.Walk(nil, acc)
		Expect(acc.results).To(HaveLen(3))
	})

})

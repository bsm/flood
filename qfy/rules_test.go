package qfy

import (
	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

var _ = Describe("factCheck", func() {
	var subject Rule

	BeforeEach(func() {
		subject = CheckFact(FactKey(33), OneOf([]int64{3, 2, 1}))
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`[33]+[1 2 3]`))
	})

	It("should have an ID", func() {
		g.Expect(subject.crc64()).To(g.Equal(uint64(3048486384098978521)))
	})

	It("should perform", func() {
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{1}}, NewState())).To(g.BeTrue())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{4}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(34): []int64{1}}, NewState())).To(g.BeFalse())
	})

	It("should capture state", func() {
		state := NewState()
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{1}}, state)).To(g.BeTrue())
		g.Expect(state.rules).To(g.Equal(map[uint64]bool{
			3048486384098978521: true,
		}))
		g.Expect(state.facts).To(g.HaveLen(1))
		g.Expect(state.facts).To(g.HaveKey(FactKey(33)))
	})
})

var _ = Describe("conjunction", func() {
	var subject Rule

	BeforeEach(func() {
		subject = All(
			CheckFact(33, OneOf([]int64{3, 2, 1})),
			CheckFact(33, OneOf([]int64{4, 5, 6})),
		)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`( [33]+[1 2 3] && [33]+[4 5 6] )`))
	})

	It("should have an ID", func() {
		g.Expect(subject.crc64()).To(g.Equal(uint64(13701729182879459540)))
	})

	It("should perform", func() {
		g.Expect(All().perform(mockFact{FactKey(33): []int64{1}}, NewState())).To(g.BeFalse())

		g.Expect(subject.perform(mockFact{FactKey(33): []int64{}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{1}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{1, 2}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{5, 2}}, NewState())).To(g.BeTrue())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{7, 2}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{7, 8}}, NewState())).To(g.BeFalse())
	})

	It("should capture state", func() {
		state := NewState()
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{5, 2}}, state)).To(g.BeTrue())
		g.Expect(state.rules).To(g.HaveLen(2))
		g.Expect(state.facts).To(g.HaveLen(1))
	})

})

var _ = Describe("disjunction", func() {
	var subject Rule

	BeforeEach(func() {
		subject = Any(
			CheckFact(33, OneOf([]int64{3, 2, 1})),
			CheckFact(34, OneOf([]int64{4, 5, 6})),
		)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`( [33]+[1 2 3] || [34]+[4 5 6] )`))
	})

	It("should have an ID", func() {
		g.Expect(subject.crc64()).To(g.Equal(uint64(17948886287937560725)))
	})

	It("should perform", func() {
		g.Expect(Any().perform(mockFact{FactKey(33): []int64{1}}, NewState())).To(g.BeFalse())

		g.Expect(subject.perform(mockFact{}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{4}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(34): []int64{7}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{4}, FactKey(34): []int64{7}}, NewState())).To(g.BeFalse())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{}, FactKey(34): []int64{5}}, NewState())).To(g.BeTrue())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{1}, FactKey(34): []int64{}}, NewState())).To(g.BeTrue())
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{1}, FactKey(34): []int64{5}}, NewState())).To(g.BeTrue())
	})

	It("should capture state", func() {
		state := NewState()
		g.Expect(subject.perform(mockFact{FactKey(33): []int64{}, FactKey(34): []int64{5}}, state)).To(g.BeTrue())
		g.Expect(state.rules).To(g.HaveLen(2))
		g.Expect(state.facts).To(g.HaveLen(2))
	})

})

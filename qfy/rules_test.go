package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("factCheck", func() {
	var subject Rule

	BeforeEach(func() {
		subject = CheckFact(FactKey(33), OneOf([]int{3, 2, 1}))
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`[33]+[1 2 3]`))
	})

	It("should have an ID", func() {
		Expect(subject.crc64()).To(Equal(uint64(3048486384098978521)))
	})

	It("should perform", func() {
		Expect(subject.perform(mockFact{FactKey(33): []int{1}}, NewState())).To(BeTrue())
		Expect(subject.perform(mockFact{FactKey(33): []int{4}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(34): []int{1}}, NewState())).To(BeFalse())
	})

	It("should capture state", func() {
		state := NewState()
		Expect(subject.perform(mockFact{FactKey(33): []int{1}}, state)).To(BeTrue())
		Expect(state.rules).To(Equal(map[uint64]bool{
			3048486384098978521: true,
		}))
		Expect(state.facts).To(HaveLen(1))
		Expect(state.facts).To(HaveKey(FactKey(33)))
	})
})

var _ = Describe("conjunction", func() {
	var subject Rule

	BeforeEach(func() {
		subject = All(
			CheckFact(33, OneOf([]int{3, 2, 1})),
			CheckFact(33, OneOf([]int{4, 5, 6})),
		)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`( [33]+[1 2 3] && [33]+[4 5 6] )`))
	})

	It("should have an ID", func() {
		Expect(subject.crc64()).To(Equal(uint64(13701729182879459540)))
	})

	It("should perform", func() {
		Expect(All().perform(mockFact{FactKey(33): []int{1}}, NewState())).To(BeFalse())

		Expect(subject.perform(mockFact{FactKey(33): []int{}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(33): []int{1}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(33): []int{1, 2}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(33): []int{5, 2}}, NewState())).To(BeTrue())
		Expect(subject.perform(mockFact{FactKey(33): []int{7, 2}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(33): []int{7, 8}}, NewState())).To(BeFalse())
	})

	It("should capture state", func() {
		state := NewState()
		Expect(subject.perform(mockFact{FactKey(33): []int{5, 2}}, state)).To(BeTrue())
		Expect(state.rules).To(HaveLen(2))
		Expect(state.facts).To(HaveLen(1))
	})

})

var _ = Describe("disjunction", func() {
	var subject Rule

	BeforeEach(func() {
		subject = Any(
			CheckFact(33, OneOf([]int{3, 2, 1})),
			CheckFact(34, OneOf([]int{4, 5, 6})),
		)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`( [33]+[1 2 3] || [34]+[4 5 6] )`))
	})

	It("should have an ID", func() {
		Expect(subject.crc64()).To(Equal(uint64(17948886287937560725)))
	})

	It("should perform", func() {
		Expect(Any().perform(mockFact{FactKey(33): []int{1}}, NewState())).To(BeFalse())

		Expect(subject.perform(mockFact{}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(33): []int{4}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(34): []int{7}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(33): []int{4}, FactKey(34): []int{7}}, NewState())).To(BeFalse())
		Expect(subject.perform(mockFact{FactKey(33): []int{}, FactKey(34): []int{5}}, NewState())).To(BeTrue())
		Expect(subject.perform(mockFact{FactKey(33): []int{1}, FactKey(34): []int{}}, NewState())).To(BeTrue())
		Expect(subject.perform(mockFact{FactKey(33): []int{1}, FactKey(34): []int{5}}, NewState())).To(BeTrue())
	})

	It("should capture state", func() {
		state := NewState()
		Expect(subject.perform(mockFact{FactKey(33): []int{}, FactKey(34): []int{5}}, state)).To(BeTrue())
		Expect(state.rules).To(HaveLen(2))
		Expect(state.facts).To(HaveLen(2))
	})

})

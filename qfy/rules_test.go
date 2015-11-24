package qfy

import (
	"github.com/bsm/intset"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inclusion", func() {
	var subject *Inclusion
	var _ Rule = subject

	BeforeEach(func() {
		subject = OneOf([]int{3, 2, 1})
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`+[1 2 3]`))
	})

	It("should have an ID", func() {
		Expect(subject.UID()).To(Equal(uint64(4016361724135366094)))
		Expect(OneOf([]int{2, 3, 1}).UID()).To(Equal(uint64(4016361724135366094)))
		Expect(OneOf([]int{7, 8, 9}).UID()).To(Equal(uint64(11523927376963847877)))
	})

	It("should check inclusion", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(intset.Use())).To(BeFalse())
		Expect(subject.Match(intset.Use(1))).To(BeTrue())
		Expect(subject.Match(intset.Use(1, 2))).To(BeTrue())
		Expect(subject.Match(intset.Use(7, 2))).To(BeTrue())
		Expect(subject.Match(intset.Use(7, 8))).To(BeFalse())
	})
})

var _ = Describe("Exclusion", func() {
	var subject *Exclusion
	var _ Rule = subject

	BeforeEach(func() {
		subject = NeitherOf([]int{3, 2, 1})
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`-[1 2 3]`))
	})

	It("should have an ID", func() {
		Expect(subject.UID()).To(Equal(uint64(17194010691906675252)))
		Expect(NeitherOf([]int{2, 3, 1}).UID()).To(Equal(uint64(17194010691906675252)))
		Expect(NeitherOf([]int{7, 8, 9}).UID()).To(Equal(uint64(5101638233538279743)))
	})

	It("should check exclusion", func() {
		Expect(subject.Match(nil)).To(BeTrue())
		Expect(subject.Match(intset.Use())).To(BeTrue())
		Expect(subject.Match(intset.Use(1))).To(BeFalse())
		Expect(subject.Match(intset.Use(1, 2))).To(BeFalse())
		Expect(subject.Match(intset.Use(7, 2))).To(BeFalse())
		Expect(subject.Match(intset.Use(7, 8))).To(BeTrue())
	})

})

var _ = Describe("Negation", func() {
	var subject *Negation
	var _ Rule = subject

	BeforeEach(func() {
		subject = Not(OneOf([]int{3, 2, 1}))
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`( NOT +[1 2 3] )`))
	})

	It("should have an ID", func() {
		Expect(subject.UID()).To(Equal(uint64(15410618248837124505)))
	})

	It("should require the opposite to match", func() {
		Expect(subject.Match(intset.Use())).To(BeTrue())
		Expect(subject.Match(intset.Use(1))).To(BeFalse())
		Expect(subject.Match(intset.Use(1, 2))).To(BeFalse())
		Expect(subject.Match(intset.Use(5, 2))).To(BeFalse())
		Expect(subject.Match(intset.Use(7, 2))).To(BeFalse())
		Expect(subject.Match(intset.Use(7, 8))).To(BeTrue())
	})

})

var _ = Describe("Conjunction", func() {
	var subject *Conjunction
	var _ Rule = subject

	BeforeEach(func() {
		subject = All(
			OneOf([]int{3, 2, 1}),
			OneOf([]int{4, 5, 6}),
		)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`( +[1 2 3] AND +[4 5 6] )`))
	})

	It("should have an ID", func() {
		Expect(subject.UID()).To(Equal(uint64(5902990494815549339)))
	})

	It("should require all to match", func() {
		Expect(All().Match(intset.Use(1))).To(BeFalse())

		Expect(subject.Match(intset.Use())).To(BeFalse())
		Expect(subject.Match(intset.Use(1))).To(BeFalse())
		Expect(subject.Match(intset.Use(1, 2))).To(BeFalse())
		Expect(subject.Match(intset.Use(5, 2))).To(BeTrue())
		Expect(subject.Match(intset.Use(7, 2))).To(BeFalse())
		Expect(subject.Match(intset.Use(7, 8))).To(BeFalse())
	})

})

var _ = Describe("Disjunction", func() {
	var subject *Disjunction
	var _ Rule = subject

	BeforeEach(func() {
		subject = Any(
			OneOf([]int{3, 2, 1}),
			OneOf([]int{4, 5, 6}),
		)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`( +[1 2 3] OR +[4 5 6] )`))
	})

	It("should have an ID", func() {
		Expect(subject.UID()).To(Equal(uint64(4689319512059230670)))
	})

	It("should require any to match", func() {
		Expect(Any().Match(intset.Use(1))).To(BeFalse())

		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(intset.Use())).To(BeFalse())
		Expect(subject.Match(intset.Use(1))).To(BeTrue())
		Expect(subject.Match(intset.Use(1, 2))).To(BeTrue())
		Expect(subject.Match(intset.Use(7, 2))).To(BeTrue())
		Expect(subject.Match(intset.Use(7, 8))).To(BeFalse())
	})

})

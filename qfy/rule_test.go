package qfy

import (
	"github.com/bsm/intset"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("plusRule", func() {
	var subject *plusRule
	var _ Rule = subject

	BeforeEach(func() {
		subject = newPlusRule([]int{3, 2, 1})
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal("+3"))
	})

	It("should have an ID", func() {
		Expect(subject.UID()).To(Equal(uint64(4016361724135366094)))
		Expect(newPlusRule([]int{2, 3, 1}).UID()).To(Equal(uint64(4016361724135366094)))
		Expect(newPlusRule([]int{7, 8, 9}).UID()).To(Equal(uint64(11523927376963847877)))
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

var _ = Describe("minusRule", func() {
	var subject *minusRule
	var _ Rule = subject

	BeforeEach(func() {
		subject = newMinusRule([]int{3, 2, 1})
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal("-3"))
	})

	It("should have an ID", func() {
		Expect(subject.UID()).To(Equal(uint64(17194010691906675252)))
		Expect(newMinusRule([]int{2, 3, 1}).UID()).To(Equal(uint64(17194010691906675252)))
		Expect(newMinusRule([]int{7, 8, 9}).UID()).To(Equal(uint64(5101638233538279743)))
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

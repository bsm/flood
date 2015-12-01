package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Equality", func() {
	var subject *Equality
	var _ Condition = subject

	BeforeEach(func() {
		subject = EqualTo(true)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`=true`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(971422227693832935)))
		Expect(EqualTo(27).CRC64()).To(Equal(uint64(9340596851114011254)))
	})

	It("should check equality", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(1)).To(BeFalse())
		Expect(subject.Match(true)).To(BeTrue())
		Expect(subject.Match(false)).To(BeFalse())
	})
})

var _ = Describe("NumericGreater", func() {
	var subject *NumericGreater
	var _ Condition = subject

	BeforeEach(func() {
		subject = GreaterThan(5.1)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`>5.1`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(16698526710376199168)))
	})

	It("should check", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(1)).To(BeFalse())
		Expect(subject.Match(5.2)).To(BeTrue())
		Expect(subject.Match(5.05)).To(BeFalse())
		Expect(subject.Match(5.1)).To(BeFalse())
		Expect(subject.Match(int64(4))).To(BeFalse())
		Expect(subject.Match(int64(7))).To(BeTrue())
	})
})

var _ = Describe("NumericLess", func() {
	var subject *NumericLess
	var _ Condition = subject

	BeforeEach(func() {
		subject = LessThan(5.1)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`<5.1`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(3414923877884655100)))
	})

	It("should check", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(1)).To(BeFalse())
		Expect(subject.Match(5.2)).To(BeFalse())
		Expect(subject.Match(5.05)).To(BeTrue())
		Expect(subject.Match(5.1)).To(BeFalse())
		Expect(subject.Match(int64(4))).To(BeTrue())
		Expect(subject.Match(int64(7))).To(BeFalse())
	})
})

var _ = Describe("NumericGreaterOrEqual", func() {
	var subject *NumericGreaterOrEqual
	var _ Condition = subject

	BeforeEach(func() {
		subject = GreaterOrEqual(5.1)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`>=5.1`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(2424392493701668213)))
	})

	It("should check", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(1)).To(BeFalse())
		Expect(subject.Match(5.2)).To(BeTrue())
		Expect(subject.Match(5.05)).To(BeFalse())
		Expect(subject.Match(5.1)).To(BeTrue())
		Expect(subject.Match(int64(4))).To(BeFalse())
		Expect(subject.Match(int64(7))).To(BeTrue())
	})
})

var _ = Describe("NumericLessOrEqual", func() {
	var subject *NumericLessOrEqual
	var _ Condition = subject

	BeforeEach(func() {
		subject = LessOrEqual(5.1)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`<=5.1`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(5028754416248954251)))
	})

	It("should check", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(1)).To(BeFalse())
		Expect(subject.Match(5.2)).To(BeFalse())
		Expect(subject.Match(5.05)).To(BeTrue())
		Expect(subject.Match(5.1)).To(BeTrue())
		Expect(subject.Match(int64(4))).To(BeTrue())
		Expect(subject.Match(int64(7))).To(BeFalse())
	})
})

var _ = Describe("NumericRange", func() {
	var subject *NumericRange
	var _ Condition = subject

	BeforeEach(func() {
		subject = Between(4.2, 6.4)
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`4.2..6.4`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(6960669049720772570)))
	})

	It("should check", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(1)).To(BeFalse())
		Expect(subject.Match(4.1)).To(BeFalse())
		Expect(subject.Match(5.2)).To(BeTrue())
		Expect(subject.Match(6.4)).To(BeTrue())
		Expect(subject.Match(6.5)).To(BeFalse())
		Expect(subject.Match(int64(4))).To(BeFalse())
		Expect(subject.Match(int64(5))).To(BeTrue())
		Expect(subject.Match(int64(6))).To(BeTrue())
		Expect(subject.Match(int64(7))).To(BeFalse())
	})
})

var _ = Describe("Inclusion", func() {
	var subject *Inclusion
	var _ Condition = subject

	BeforeEach(func() {
		subject = OneOf([]int64{3, 2, 1})
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`+[1 2 3]`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(4016361724135366094)))
		Expect(OneOf([]int64{2, 3, 1}).CRC64()).To(Equal(uint64(4016361724135366094)))
		Expect(OneOf([]int64{7, 8, 9}).CRC64()).To(Equal(uint64(11523927376963847877)))
	})

	It("should check inclusion", func() {
		Expect(subject.Match(nil)).To(BeFalse())
		Expect(subject.Match(int64(1))).To(BeTrue())
		Expect(subject.Match(int64(7))).To(BeFalse())

		Expect(subject.Match(SortInts64())).To(BeFalse())
		Expect(subject.Match(SortInts64(1))).To(BeTrue())
		Expect(subject.Match(SortInts64(1, 2))).To(BeTrue())
		Expect(subject.Match(SortInts64(7, 2))).To(BeTrue())
		Expect(subject.Match(SortInts64(7, 8))).To(BeFalse())
	})
})

var _ = Describe("Exclusion", func() {
	var subject *Exclusion
	var _ Condition = subject

	BeforeEach(func() {
		subject = NoneOf([]int64{3, 2, 1})
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`-[1 2 3]`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(17194010691906675252)))
		Expect(NoneOf([]int64{2, 3, 1}).CRC64()).To(Equal(uint64(17194010691906675252)))
		Expect(NoneOf([]int64{7, 8, 9}).CRC64()).To(Equal(uint64(5101638233538279743)))
	})

	It("should check exclusion", func() {
		Expect(subject.Match(nil)).To(BeTrue())
		Expect(subject.Match(int64(1))).To(BeFalse())
		Expect(subject.Match(int64(7))).To(BeTrue())

		Expect(subject.Match(SortInts64())).To(BeTrue())
		Expect(subject.Match(SortInts64(1))).To(BeFalse())
		Expect(subject.Match(SortInts64(1, 2))).To(BeFalse())
		Expect(subject.Match(SortInts64(7, 2))).To(BeFalse())
		Expect(subject.Match(SortInts64(7, 8))).To(BeTrue())
	})

})

var _ = Describe("Negation", func() {
	var subject *Negation
	var _ Condition = subject

	BeforeEach(func() {
		subject = Not(OneOf([]int64{3, 2, 1}))
	})

	It("should return a string", func() {
		Expect(subject.String()).To(Equal(`!+[1 2 3]`))
	})

	It("should have an ID", func() {
		Expect(subject.CRC64()).To(Equal(uint64(5849970728172636964)))
	})

	It("should require the opposite to match", func() {
		Expect(subject.Match(nil)).To(BeTrue())
		Expect(subject.Match(int64(1))).To(BeFalse())
		Expect(subject.Match(int64(7))).To(BeTrue())

		Expect(subject.Match(SortInts64())).To(BeTrue())
		Expect(subject.Match(SortInts64(1))).To(BeFalse())
		Expect(subject.Match(SortInts64(1, 2))).To(BeFalse())
		Expect(subject.Match(SortInts64(5, 2))).To(BeFalse())
		Expect(subject.Match(SortInts64(7, 2))).To(BeFalse())
		Expect(subject.Match(SortInts64(7, 8))).To(BeTrue())
	})

})

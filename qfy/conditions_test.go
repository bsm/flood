package qfy

import (
	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

var _ = Describe("Equality", func() {
	var subject *Equality
	var _ Condition = subject

	BeforeEach(func() {
		subject = EqualTo(true)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`=true`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(971422227693832935)))
		g.Expect(EqualTo(27).CRC64()).To(g.Equal(uint64(9340596851114011254)))
	})

	It("should check equality", func() {
		g.Expect(subject.Match(nil)).To(g.BeFalse())
		g.Expect(subject.Match(1)).To(g.BeFalse())
		g.Expect(subject.Match(true)).To(g.BeTrue())
		g.Expect(subject.Match(false)).To(g.BeFalse())
	})
})

var _ = Describe("NumericGreater", func() {
	var subject *NumericGreater
	var _ Condition = subject

	BeforeEach(func() {
		subject = GreaterThan(5.1)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`>5.1`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(16698526710376199168)))
	})

	It("should check", func() {
		g.Expect(subject.Match(nil)).To(g.BeFalse())
		g.Expect(subject.Match(1)).To(g.BeFalse())
		g.Expect(subject.Match(5.2)).To(g.BeTrue())
		g.Expect(subject.Match(5.05)).To(g.BeFalse())
		g.Expect(subject.Match(5.1)).To(g.BeFalse())
		g.Expect(subject.Match(int64(4))).To(g.BeFalse())
		g.Expect(subject.Match(int64(7))).To(g.BeTrue())
	})
})

var _ = Describe("NumericLess", func() {
	var subject *NumericLess
	var _ Condition = subject

	BeforeEach(func() {
		subject = LessThan(5.1)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`<5.1`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(3414923877884655100)))
	})

	It("should check", func() {
		g.Expect(subject.Match(nil)).To(g.BeFalse())
		g.Expect(subject.Match(1)).To(g.BeFalse())
		g.Expect(subject.Match(5.2)).To(g.BeFalse())
		g.Expect(subject.Match(5.05)).To(g.BeTrue())
		g.Expect(subject.Match(5.1)).To(g.BeFalse())
		g.Expect(subject.Match(int64(4))).To(g.BeTrue())
		g.Expect(subject.Match(int64(7))).To(g.BeFalse())
	})
})

var _ = Describe("NumericGreaterOrEqual", func() {
	var subject *NumericGreaterOrEqual
	var _ Condition = subject

	BeforeEach(func() {
		subject = GreaterOrEqual(5.1)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`>=5.1`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(2424392493701668213)))
	})

	It("should check", func() {
		g.Expect(subject.Match(nil)).To(g.BeFalse())
		g.Expect(subject.Match(1)).To(g.BeFalse())
		g.Expect(subject.Match(5.2)).To(g.BeTrue())
		g.Expect(subject.Match(5.05)).To(g.BeFalse())
		g.Expect(subject.Match(5.1)).To(g.BeTrue())
		g.Expect(subject.Match(int64(4))).To(g.BeFalse())
		g.Expect(subject.Match(int64(7))).To(g.BeTrue())
	})
})

var _ = Describe("NumericLessOrEqual", func() {
	var subject *NumericLessOrEqual
	var _ Condition = subject

	BeforeEach(func() {
		subject = LessOrEqual(5.1)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`<=5.1`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(5028754416248954251)))
	})

	It("should check", func() {
		g.Expect(subject.Match(nil)).To(g.BeFalse())
		g.Expect(subject.Match(1)).To(g.BeFalse())
		g.Expect(subject.Match(5.2)).To(g.BeFalse())
		g.Expect(subject.Match(5.05)).To(g.BeTrue())
		g.Expect(subject.Match(5.1)).To(g.BeTrue())
		g.Expect(subject.Match(int64(4))).To(g.BeTrue())
		g.Expect(subject.Match(int64(7))).To(g.BeFalse())
	})
})

var _ = Describe("NumericRange", func() {
	var subject *NumericRange
	var _ Condition = subject

	BeforeEach(func() {
		subject = Between(4.2, 6.4)
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`4.2..6.4`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(6960669049720772570)))
	})

	It("should check", func() {
		g.Expect(subject.Match(nil)).To(g.BeFalse())
		g.Expect(subject.Match(1)).To(g.BeFalse())
		g.Expect(subject.Match(4.1)).To(g.BeFalse())
		g.Expect(subject.Match(5.2)).To(g.BeTrue())
		g.Expect(subject.Match(6.4)).To(g.BeTrue())
		g.Expect(subject.Match(6.5)).To(g.BeFalse())
		g.Expect(subject.Match(int64(4))).To(g.BeFalse())
		g.Expect(subject.Match(int64(5))).To(g.BeTrue())
		g.Expect(subject.Match(int64(6))).To(g.BeTrue())
		g.Expect(subject.Match(int64(7))).To(g.BeFalse())
	})
})

var _ = Describe("Inclusion", func() {
	var subject *Inclusion
	var _ Condition = subject

	BeforeEach(func() {
		subject = OneOf([]int64{3, 2, 1})
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`+[1 2 3]`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(4016361724135366094)))
		g.Expect(OneOf([]int64{2, 3, 1}).CRC64()).To(g.Equal(uint64(4016361724135366094)))
		g.Expect(OneOf([]int64{7, 8, 9}).CRC64()).To(g.Equal(uint64(11523927376963847877)))
	})

	It("should check inclusion", func() {
		g.Expect(subject.Match(nil)).To(g.BeFalse())
		g.Expect(subject.Match(int64(1))).To(g.BeTrue())
		g.Expect(subject.Match(int64(7))).To(g.BeFalse())

		g.Expect(subject.Match(SortInts64())).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(1))).To(g.BeTrue())
		g.Expect(subject.Match(SortInts64(1, 2))).To(g.BeTrue())
		g.Expect(subject.Match(SortInts64(7, 2))).To(g.BeTrue())
		g.Expect(subject.Match(SortInts64(7, 8))).To(g.BeFalse())
	})
})

var _ = Describe("Exclusion", func() {
	var subject *Exclusion
	var _ Condition = subject

	BeforeEach(func() {
		subject = NoneOf([]int64{3, 2, 1})
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`-[1 2 3]`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(17194010691906675252)))
		g.Expect(NoneOf([]int64{2, 3, 1}).CRC64()).To(g.Equal(uint64(17194010691906675252)))
		g.Expect(NoneOf([]int64{7, 8, 9}).CRC64()).To(g.Equal(uint64(5101638233538279743)))
	})

	It("should check exclusion", func() {
		g.Expect(subject.Match(nil)).To(g.BeTrue())
		g.Expect(subject.Match(int64(1))).To(g.BeFalse())
		g.Expect(subject.Match(int64(7))).To(g.BeTrue())

		g.Expect(subject.Match(SortInts64())).To(g.BeTrue())
		g.Expect(subject.Match(SortInts64(1))).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(1, 2))).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(7, 2))).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(7, 8))).To(g.BeTrue())
	})

})

var _ = Describe("Negation", func() {
	var subject *Negation
	var _ Condition = subject

	BeforeEach(func() {
		subject = Not(OneOf([]int64{3, 2, 1}))
	})

	It("should return a string", func() {
		g.Expect(subject.String()).To(g.Equal(`!+[1 2 3]`))
	})

	It("should have an ID", func() {
		g.Expect(subject.CRC64()).To(g.Equal(uint64(5849970728172636964)))
	})

	It("should require the opposite to match", func() {
		g.Expect(subject.Match(nil)).To(g.BeTrue())
		g.Expect(subject.Match(int64(1))).To(g.BeFalse())
		g.Expect(subject.Match(int64(7))).To(g.BeTrue())

		g.Expect(subject.Match(SortInts64())).To(g.BeTrue())
		g.Expect(subject.Match(SortInts64(1))).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(1, 2))).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(5, 2))).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(7, 2))).To(g.BeFalse())
		g.Expect(subject.Match(SortInts64(7, 8))).To(g.BeTrue())
	})

})

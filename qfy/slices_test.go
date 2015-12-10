package qfy

import (
	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

var _ = Describe("Ints64", func() {
	var subject Ints64

	BeforeEach(func() {
		subject = SortInts64(4, 6, 2)
	})

	It("should normalize", func() {
		g.Expect(subject).To(g.Equal(Ints64{2, 4, 6}))
		g.Expect(subject.crc64('x')).To(g.Equal(uint64(6934466117131854228)))
	})

	It("should check if exists", func() {
		g.Expect(subject.Exists(1)).To(g.BeFalse())
		g.Expect(subject.Exists(2)).To(g.BeTrue())
		g.Expect(subject.Exists(3)).To(g.BeFalse())
		g.Expect(subject.Exists(4)).To(g.BeTrue())
	})

	It("should check for intersections", func() {
		g.Expect(subject.Inter(SortInts64(3))).To(g.BeFalse())
		g.Expect(subject.Inter(SortInts64(3, 5))).To(g.BeFalse())
		g.Expect(subject.Inter(SortInts64(3, 4, 5, 7))).To(g.BeTrue())
	})

})

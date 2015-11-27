package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ints64", func() {
	var subject Ints64

	BeforeEach(func() {
		subject = SortInts64(4, 6, 2)
	})

	It("should normalize", func() {
		Expect(subject).To(Equal(Ints64{2, 4, 6}))
		Expect(subject.crc64('x')).To(Equal(uint64(6934466117131854228)))
	})

	It("should check if exists", func() {
		Expect(subject.Exists(1)).To(BeFalse())
		Expect(subject.Exists(2)).To(BeTrue())
		Expect(subject.Exists(3)).To(BeFalse())
		Expect(subject.Exists(4)).To(BeTrue())
	})

	It("should check for intersections", func() {
		Expect(subject.Inter(SortInts64(3))).To(BeFalse())
		Expect(subject.Inter(SortInts64(3, 5))).To(BeFalse())
		Expect(subject.Inter(SortInts64(3, 4, 5, 7))).To(BeTrue())
	})

})

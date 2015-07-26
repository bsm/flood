package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("strDict", func() {
	var subject strDict

	BeforeEach(func() {
		subject = strDict{}
	})

	It("should fetch items", func() {
		Expect(subject.Fetch("a")).To(Equal(1))
		Expect(subject.Fetch("b")).To(Equal(2))
		Expect(subject.Fetch("c")).To(Equal(3))
		Expect(subject.Fetch("b")).To(Equal(2))
		Expect(subject.Fetch("a")).To(Equal(1))
		Expect(subject).To(Equal(strDict{"a": 1, "b": 2, "c": 3}))
	})

	It("should fetch slices", func() {
		Expect(subject.FetchSlice("a", "b")).To(Equal([]int{1, 2}))
		Expect(subject.FetchSlice("c", "b")).To(Equal([]int{3, 2}))
	})

	It("should get slices", func() {
		Expect(subject.FetchSlice("a", "b")).To(Equal([]int{1, 2}))
		Expect(subject.GetSlice("c", "b", "d", "a")).To(Equal([]int{2, 1}))
	})

})

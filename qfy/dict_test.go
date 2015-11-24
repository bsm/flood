package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dict", func() {
	var subject Dict

	BeforeEach(func() {
		subject = NewDict()
	})

	It("should find or create items", func() {
		Expect(subject.Add("a")).To(Equal(1))
		Expect(subject.Add("b")).To(Equal(2))
		Expect(subject.Add("c")).To(Equal(3))
		Expect(subject.Add("b")).To(Equal(2))
		Expect(subject.Add("a")).To(Equal(1))
		Expect(subject).To(Equal(Dict{"a": 1, "b": 2, "c": 3}))
	})

	It("should add slices", func() {
		Expect(subject.AddSlice("a", "b")).To(Equal([]int{1, 2}))
		Expect(subject.AddSlice("c", "b")).To(Equal([]int{3, 2}))
	})

	It("should get slices", func() {
		Expect(subject.AddSlice("a", "b")).To(Equal([]int{1, 2}))
		Expect(subject.GetSlice("c", "b", "d", "a")).To(Equal([]int{2, 1}))
	})

})

var _ = Describe("ConcurrentDict", func() {
	var subject *ConcurrentDict

	BeforeEach(func() {
		subject = NewConcurrentDict()
	})

	It("should find or create items", func() {
		Expect(subject.Add("a")).To(Equal(1))
		Expect(subject.Add("b")).To(Equal(2))
		Expect(subject.Add("c")).To(Equal(3))
		Expect(subject.Add("b")).To(Equal(2))
		Expect(subject.Add("a")).To(Equal(1))
		Expect(subject.dict).To(Equal(Dict{"a": 1, "b": 2, "c": 3}))
	})

	It("should add slices", func() {
		Expect(subject.AddSlice("a", "b")).To(Equal([]int{1, 2}))
		Expect(subject.AddSlice("c", "b")).To(Equal([]int{3, 2}))
	})

	It("should get slices", func() {
		Expect(subject.AddSlice("a", "b")).To(Equal([]int{1, 2}))
		Expect(subject.GetSlice("c", "b", "d", "a")).To(Equal([]int{2, 1}))
	})

})

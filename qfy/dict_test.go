package qfy

import (
	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

var _ = Describe("Dict", func() {
	var subject Dict

	BeforeEach(func() {
		subject = NewDict()
	})

	It("should find or create items", func() {
		g.Expect(subject.Add("a")).To(g.Equal(int64(1)))
		g.Expect(subject.Add("b")).To(g.Equal(int64(2)))
		g.Expect(subject.Add("c")).To(g.Equal(int64(3)))
		g.Expect(subject.Add("b")).To(g.Equal(int64(2)))
		g.Expect(subject.Add("a")).To(g.Equal(int64(1)))
		g.Expect(subject).To(g.Equal(Dict{"a": 1, "b": 2, "c": 3}))
	})

	It("should add slices", func() {
		g.Expect(subject.AddSlice("a", "b")).To(g.Equal([]int64{1, 2}))
		g.Expect(subject.AddSlice("c", "b")).To(g.Equal([]int64{3, 2}))
	})

	It("should get slices", func() {
		g.Expect(subject.AddSlice("a", "b")).To(g.Equal([]int64{1, 2}))
		g.Expect(subject.GetSlice("c", "b", "d", "a")).To(g.Equal([]int64{2, 1}))
	})

})

var _ = Describe("ConcurrentDict", func() {
	var subject *ConcurrentDict

	BeforeEach(func() {
		subject = NewConcurrentDict()
	})

	It("should find or create items", func() {
		g.Expect(subject.Add("a")).To(g.Equal(int64(1)))
		g.Expect(subject.Add("b")).To(g.Equal(int64(2)))
		g.Expect(subject.Add("c")).To(g.Equal(int64(3)))
		g.Expect(subject.Add("b")).To(g.Equal(int64(2)))
		g.Expect(subject.Add("a")).To(g.Equal(int64(1)))
		g.Expect(subject.dict).To(g.Equal(Dict{"a": 1, "b": 2, "c": 3}))
	})

	It("should add slices", func() {
		g.Expect(subject.AddSlice("a", "b")).To(g.Equal([]int64{1, 2}))
		g.Expect(subject.AddSlice("c", "b")).To(g.Equal([]int64{3, 2}))
	})

	It("should get slices", func() {
		g.Expect(subject.AddSlice("a", "b")).To(g.Equal([]int64{1, 2}))
		g.Expect(subject.GetSlice("c", "b", "d", "a")).To(g.Equal([]int64{2, 1}))
	})

})

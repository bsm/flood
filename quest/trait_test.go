package quest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("stringHash", func() {
	var subject stringHash
	var _ matchData = subject

	BeforeEach(func() {
		subject = make(stringHash)
	})

	It("should check", func() {
		Expect(subject.Check(&Condition{"name", ComparatorEqual, "value"})).NotTo(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorGreater, "value"})).To(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorEqual, 0})).To(HaveOccurred())
	})

	It("should store/find", func() {
		subject.Store(ruleReference{90, 1}, "value")

		res, err := subject.Find("value")
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal([]ruleReference{{90, 1}}))

		res, err = subject.Find("missing")
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeNil())

		_, err = subject.Find(99)
		Expect(err).To(HaveOccurred())
	})

})

var _ = Describe("int64Hash", func() {
	var subject int64Hash
	var _ matchData = subject

	BeforeEach(func() {
		subject = make(int64Hash)
	})

	It("should check", func() {
		Expect(subject.Check(&Condition{"name", ComparatorEqual, int64(27)})).NotTo(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorGreater, int64(27)})).To(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorEqual, "value"})).To(HaveOccurred())
	})

	It("should store/find", func() {
		subject.Store(ruleReference{90, 1}, int64(27))

		res, err := subject.Find(int64(27))
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal([]ruleReference{{90, 1}}))

		res, err = subject.Find(int64(28))
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeNil())

		_, err = subject.Find("value")
		Expect(err).To(HaveOccurred())
	})

})

var _ = Describe("int32Hash", func() {
	var subject int32Hash
	var _ matchData = subject

	BeforeEach(func() {
		subject = make(int32Hash)
	})

	It("should check", func() {
		Expect(subject.Check(&Condition{"name", ComparatorEqual, int32(27)})).NotTo(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorGreater, int32(27)})).To(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorEqual, "value"})).To(HaveOccurred())
	})

	It("should store/find", func() {
		subject.Store(ruleReference{90, 1}, int32(27))

		res, err := subject.Find(int32(27))
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal([]ruleReference{{90, 1}}))

		res, err = subject.Find(int32(28))
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeNil())

		_, err = subject.Find("value")
		Expect(err).To(HaveOccurred())
	})

})

var _ = Describe("boolHash", func() {
	var subject boolHash
	var _ matchData = subject

	BeforeEach(func() {
		subject = make(boolHash)
	})

	It("should check", func() {
		Expect(subject.Check(&Condition{"name", ComparatorEqual, true})).NotTo(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorGreater, true})).To(HaveOccurred())
		Expect(subject.Check(&Condition{"name", ComparatorEqual, "value"})).To(HaveOccurred())
	})

	It("should store/find", func() {
		subject.Store(ruleReference{90, 1}, true)

		res, err := subject.Find(true)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(Equal([]ruleReference{{90, 1}}))

		res, err = subject.Find(false)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeNil())

		_, err = subject.Find("value")
		Expect(err).To(HaveOccurred())
	})

})

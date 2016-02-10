package quest

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("checksum", func() {
	var subject checksum
	var bits = func(c checksum) string {
		parts := make([]string, len(c.data))
		for i, b := range c.data {
			parts[i] = fmt.Sprintf("%b", b)
		}
		return strings.Join(parts, ".")
	}

	It("should set and clear", func() {
		subject = subject.Mark(0, true)
		subject = subject.Mark(1, true)
		subject = subject.Mark(2, true)
		subject = subject.Mark(3, true)
		subject = subject.Mark(4, true)
		subject = subject.Mark(5, true)
		subject = subject.Mark(6, true)
		Expect(subject.size).To(Equal(7))
		Expect(bits(subject)).To(Equal("1111111"))

		subject = subject.Mark(7, false)
		Expect(subject.size).To(Equal(8))
		Expect(bits(subject)).To(Equal("1111111"))

		subject = subject.Mark(1, false)
		Expect(bits(subject)).To(Equal("1111101"))
		subject = subject.Mark(3, false)
		Expect(bits(subject)).To(Equal("1110101"))
		subject = subject.Mark(5, false)
		Expect(bits(subject)).To(Equal("1010101"))
	})

	It("should expand", func() {
		subject = subject.Mark(52, true)
		subject = subject.Mark(62, true)
		subject = subject.Mark(34, true)
		Expect(subject).To(Equal(checksum{
			data: []byte{85, 0, 0, 0, 4, 0, 16, 64},
			size: 63,
		}))
	})

	It("should compare", func() {
		a, b := checksum{}, checksum{}
		Expect(a.Equal(b)).To(BeTrue())

		a = a.Mark(0, false)
		Expect(a.Equal(b)).To(BeFalse())

		b = b.Mark(0, false)
		Expect(a.Equal(b)).To(BeTrue())

		a = a.Mark(4, true)
		Expect(a.Equal(b)).To(BeFalse())

		b = b.Mark(4, true)
		a = a.Mark(5, false)
		Expect(a.Equal(b)).To(BeFalse())

		b = b.Mark(5, false)
		b = b.Mark(6, true)
		a = a.Mark(6, true)
		Expect(a.Equal(b)).To(BeTrue())
	})

	It("should create unions", func() {
		a := checksum{}
		a = a.Mark(7, true)
		a = a.Mark(8, true)
		a = a.Mark(9, true)
		Expect(a.size).To(Equal(10))
		Expect(a.data).To(Equal([]byte{128, 3}))

		b := checksum{}
		b = b.Mark(8, true)
		b = b.Mark(15, true)
		b = b.Mark(16, true)
		Expect(b.size).To(Equal(17))
		Expect(b.data).To(Equal([]byte{0, 129, 1}))

		Expect(a.Union(b).size).To(Equal(17))
		Expect(a.Union(b).data).To(Equal([]byte{128, 131, 1}))

		Expect(b.Union(a).size).To(Equal(17))
		Expect(b.Union(a).data).To(Equal([]byte{128, 131, 1}))
	})

})

var _ = Describe("checksums", func() {
	var subject checksums

	BeforeEach(func() {
		subject = make(checksums)

		subject.Mark(100, 0, true)
		subject.Mark(100, 1, false)
		subject.Mark(100, 2, true)

		subject.Mark(101, 0, false)
		subject.Mark(101, 1, true)
	})

	It("should update", func() {
		Expect(subject).To(Equal(checksums{
			100: checksum{size: 3, data: []byte{5}},
			101: checksum{size: 2, data: []byte{2}},
		}))
	})

})

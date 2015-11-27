package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("CRC64", func() {

	DescribeTable("Sum64",
		func(prefix byte, factors []uint64, expected uint64) {
			hash := NewCRC64(byte(prefix), len(factors))
			hash.Add(factors...)
			Expect(hash.Sum64()).To(Equal(expected))
		},
		Entry("ordered", byte('+'), []uint64{12, 43, 76, 87}, uint64(16697346874648777555)),
		Entry("shuffled", byte('+'), []uint64{87, 12, 76, 43}, uint64(16697346874648777555)),
		Entry("negative", byte('-'), []uint64{87, 12, 76, 43}, uint64(11938658858315989027)),
		Entry("blank", byte('+'), []uint64{}, uint64(6093685733581172889)),
		Entry("nil", byte('+'), ([]uint64)(nil), uint64(6093685733581172889)),
	)

})

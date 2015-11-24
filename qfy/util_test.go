package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("crcHash", func() {

	DescribeTable("Sum64",
		func(prefix byte, factors []uint64, expected uint64) {
			hash := newCRCHash(byte(prefix), len(factors))
			for _, f := range factors {
				hash.Add(f)
			}
			Expect(hash.Sum64()).To(Equal(expected))
		},
		Entry("ordered", byte('+'), []uint64{12, 43, 76, 87}, uint64(16697346874648777555)),
		Entry("shuffled", byte('+'), []uint64{87, 12, 76, 43}, uint64(16697346874648777555)),
		Entry("negative", byte('-'), []uint64{87, 12, 76, 43}, uint64(11938658858315989027)),
		Entry("blank", byte('+'), []uint64{}, uint64(6093685733581172889)),
	)

})

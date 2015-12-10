package qfy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	g "github.com/onsi/gomega"
)

var _ = Describe("CRC64", func() {

	DescribeTable("Sum64",
		func(prefix byte, factors []uint64, expected uint64) {
			hash := NewCRC64(byte(prefix), len(factors))
			hash.Add(factors...)
			g.Expect(hash.Sum64()).To(g.Equal(expected))
		},
		Entry("ordered", byte('+'), []uint64{12, 43, 76, 87}, uint64(16697346874648777555)),
		Entry("shuffled", byte('+'), []uint64{87, 12, 76, 43}, uint64(16697346874648777555)),
		Entry("negative", byte('-'), []uint64{87, 12, 76, 43}, uint64(11938658858315989027)),
		Entry("blank", byte('+'), []uint64{}, uint64(6093685733581172889)),
		Entry("nil", byte('+'), ([]uint64)(nil), uint64(6093685733581172889)),
	)

	It("should build CRC64s from any value", func() {
		g.Expect(crc64FromValue('x', 1.23)).To(g.Equal(uint64(5538764236368787474)))
		g.Expect(crc64FromValue('x', 1.23)).To(g.Equal(uint64(5538764236368787474)))
		g.Expect(crc64FromValue('x', 3.21)).To(g.Equal(uint64(6459098420131569288)))
		g.Expect(crc64FromValue('y', 3.21)).To(g.Equal(uint64(4453996903255678582)))

		g.Expect(crc64FromValue('x', 123)).To(g.Equal(uint64(17364521479190536253)))
		g.Expect(crc64FromValue('x', -123)).To(g.Equal(uint64(6651416055640416575)))

		g.Expect(crc64FromValue('x', true)).To(g.Equal(uint64(6688660444647243956)))
		g.Expect(crc64FromValue('x', false)).To(g.Equal(uint64(4088829085531488330)))
		g.Expect(crc64FromValue('x', true)).To(g.Equal(uint64(6688660444647243956)))

		g.Expect(crc64FromValue('x', "string")).To(g.Equal(uint64(526724301549245503)))
		g.Expect(crc64FromValue('x', "STRING")).To(g.Equal(uint64(5405123842493546239)))
	})

})

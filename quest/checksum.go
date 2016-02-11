package quest

import (
	"bytes"
	"sync"
)

type checksum struct {
	data []byte
	size int
}

func (c checksum) Equal(o checksum) bool {
	return c.size == o.size && bytes.Equal(c.data, o.data)
}

func (c checksum) Union(o checksum) checksum {
	if c.size > o.size {
		return o.Union(c)
	}

	res := checksum{
		data: make([]byte, len(o.data)),
		size: o.size,
	}
	copy(res.data, o.data)

	for i, b := range c.data {
		res.data[i] = res.data[i] | b
	}
	return res
}

func (c checksum) Mark(n int, value bool) checksum {
	pos := n / 8
	if pos >= len(c.data) {
		data := make([]byte, pos+1)
		copy(data, c.data)
		c.data = data
	}

	rem := byte(1) << uint(n%8)
	if value {
		c.data[pos] |= rem
	} else {
		c.data[pos] &^= rem
	}
	if size := n + 1; size > c.size {
		c.size = size
	}
	return c
}

// --------------------------------------------------------------------

var checksumsPool sync.Pool

type checksums map[Outcome]checksum

// Marks marks rules index of an outcome
func (c checksums) Mark(outcome Outcome, index int, value bool) {
	c[outcome] = c[outcome].Mark(index, value)
}

// Reset quickly resets the checksums
func (c checksums) Reset() {
	for _, csum := range c {
		csum.size = 0
	}
}

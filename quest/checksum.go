package quest

import (
	"bytes"
	"sync"
)

type checksum struct {
	data []byte
	size int
}

func (c checksum) Reset() checksum {
	c.size = 0
	c.data = c.data[:0]
	return c
}

func (c checksum) Equal(o checksum) bool {
	return c.size == o.size && bytes.Equal(c.data, o.data)
}

func (c checksum) Mark(n int, value bool) checksum {
	pos := n / 8
	for i := len(c.data); i <= pos; i++ {
		c.data = append(c.data, 0)
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
	for outcome := range c {
		c[outcome] = c[outcome].Reset()
	}
}

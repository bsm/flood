package qfy

import (
	"encoding/binary"
	"encoding/gob"
	"hash/crc64"
	"sort"
)

var crcTable = crc64.MakeTable(crc64.ECMA)

// CRC64 is used by conditions to calculate unique and consistent CRC64 sum
type CRC64 struct {
	sign    byte
	factors uint64Slice
}

// NewCRC64 initializes a new hash generator
func NewCRC64(sign byte, capacity int) *CRC64 {
	return &CRC64{sign, make(uint64Slice, 0, capacity)}
}

// Add adds factors to the CRC hash
func (h *CRC64) Add(factors ...uint64) { h.factors = append(h.factors, factors...) }

// Sum64 calculates a numeric hash sum
func (h *CRC64) Sum64() uint64 {
	sort.Sort(h.factors)

	hash := crc64.New(crcTable)
	hash.Write([]byte{h.sign})
	for _, factor := range h.factors {
		bin := make([]byte, 8)
		binary.LittleEndian.PutUint64(bin, factor)
		hash.Write(bin)
	}
	return hash.Sum64()
}

// --------------------------------------------------------------------

func crc64FromValue(sign byte, v interface{}) uint64 {
	sub := crc64.New(crcTable)
	gob.NewEncoder(sub).Encode(v)

	hash := NewCRC64(sign, 1)
	hash.Add(sub.Sum64())
	return hash.Sum64()
}

func crc64FromRules(sign byte, sources ...Rule) uint64 {
	hash := NewCRC64(sign, len(sources))
	for _, s := range sources {
		hash.Add(s.crc64())
	}
	return hash.Sum64()
}

func crc64FromConditions(sign byte, sources ...Condition) uint64 {
	hash := NewCRC64(sign, len(sources))
	for _, s := range sources {
		hash.Add(s.CRC64())
	}
	return hash.Sum64()
}

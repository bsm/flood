package qfy

import (
	"encoding/binary"
	"hash/crc64"
	"sort"

	"github.com/bsm/intset"
)

var crcTable = crc64.MakeTable(crc64.ECMA)

type crcHash struct {
	prefix  byte
	factors uint64Slice
}

func newCRCHash(prefix byte, capacity int) *crcHash {
	return &crcHash{prefix, make(uint64Slice, 0, capacity)}
}

func (h *crcHash) Add(factor uint64) { h.factors = append(h.factors, factor) }

func (h *crcHash) Sum64() uint64 {
	sort.Sort(h.factors)

	data := make([]byte, len(h.factors)*8+1)
	data[0] = h.prefix
	for i, factor := range h.factors {
		binary.LittleEndian.PutUint64(data[i*8+1:], factor)
	}
	return crc64.Checksum(data, crcTable)
}

// --------------------------------------------------------------------

type uint64Slice []uint64

func (p uint64Slice) Len() int           { return len(p) }
func (p uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// --------------------------------------------------------------------

// the qualification lookup process abstraction
type lookup struct {
	results   []int
	ruleCache map[string]map[uint64]bool
	factCache map[string]*intset.Set
}

func newLookup() *lookup {
	return &lookup{
		results:   make([]int, 0, 100),
		ruleCache: make(map[string]map[uint64]bool, 1000),
		factCache: make(map[string]*intset.Set, 20),
	}
}

func (l *lookup) Clear() {
	l.results = l.results[:0]
	for _, v := range l.ruleCache {
		for k, _ := range v {
			delete(v, k)
		}
	}
	for k, _ := range l.factCache {
		delete(l.factCache, k)
	}
}

package qfy

import "sort"

// Ints64 is a sorted and optimised []int64 slice
type Ints64 []int64

// SortInts64 creates a new sorted Ints64
func SortInts64(vv ...int64) Ints64 {
	is := Ints64(vv)
	sort.Sort(is)
	return is
}

func (p Ints64) Len() int           { return len(p) }
func (p Ints64) Less(i, j int) bool { return p[i] < p[j] }
func (p Ints64) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Search searches for an item in the slice
func (p Ints64) Search(x int64) int {
	return sort.Search(len(p), func(i int) bool { return p[i] >= x })
}

// Exists checks the existence
func (p Ints64) Exists(v int64) bool {
	pos := p.Search(v)
	return pos < len(p) && p[pos] == v
}

// Inter checks if intersectable
func (p Ints64) Inter(q Ints64) bool {
	lp, lq := len(p), len(q)
	if lq < lp {
		lp, lq = lq, lp
		p, q = q, p
	}
	if lp == 0 || p[0] > q[lq-1] || q[0] > p[lp-1] {
		return false
	}

	offset := 0
	for _, v := range p {
		pos := q[offset:].Search(v) + offset
		if pos < lq && q[pos] == v {
			return true
		} else if pos >= lq {
			return false
		}
		offset = pos
	}
	return false
}

func (p Ints64) crc64(sign byte) uint64 {
	hash := NewCRC64(sign, len(p))
	for _, n := range p {
		hash.Add(uint64(n))
	}
	return hash.Sum64()
}

// --------------------------------------------------------------------

func ints64FromInts(v []int) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromInts8(v []int8) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromInts16(v []int16) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromInts32(v []int32) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromUints(v []uint) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromUints8(v []uint8) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromUints16(v []uint16) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromUints32(v []uint32) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}
func ints64FromUints64(v []uint64) Ints64 {
	vv := make([]int64, len(v))
	for i, n := range v {
		vv[i] = int64(n)
	}
	return SortInts64(vv...)
}

// --------------------------------------------------------------------

type uint64Slice []uint64

func (p uint64Slice) Len() int           { return len(p) }
func (p uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

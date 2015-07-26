package qfy

import (
	"encoding/binary"
	"fmt"
	"hash/crc64"

	"github.com/bsm/intset"
)

type Rule interface {
	UID() uint64
	Match(*intset.Set) bool
	String() string
}

type RuleSet []Rule

func (p RuleSet) Len() int           { return len(p) }
func (p RuleSet) Less(i, j int) bool { return p[i].UID() < p[j].UID() }
func (p RuleSet) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// --------------------------------------------------------------------

type baseRule struct {
	hash uint64
	vals *intset.Set
}

func newBaseRule(sign byte, vals []int) *baseRule {
	vset := intset.Use(vals...)
	salt := make([]byte, len(vals)*8+1)
	salt[0] = sign
	for i, val := range vset.Slice() {
		binary.LittleEndian.PutUint64(salt[i*8+1:], uint64(val))
	}

	return &baseRule{
		hash: crc64.Checksum(salt, crcTable),
		vals: vset,
	}
}

// UID returns a unique rule identifier
func (r *baseRule) UID() uint64 { return r.hash }

// --------------------------------------------------------------------

type plusRule struct{ baseRule }

func newPlusRule(vals []int) *plusRule        { return &plusRule{*newBaseRule('+', vals)} }
func (r *plusRule) Match(vv *intset.Set) bool { return vv != nil && r.vals.Intersects(vv) }
func (r *plusRule) String() string            { return fmt.Sprintf("+%d", r.vals.Len()) }

type minusRule struct{ baseRule }

func newMinusRule(vals []int) *minusRule       { return &minusRule{*newBaseRule('-', vals)} }
func (r *minusRule) Match(vv *intset.Set) bool { return !(vv != nil && r.vals.Intersects(vv)) }
func (r *minusRule) String() string            { return fmt.Sprintf("-%d", r.vals.Len()) }

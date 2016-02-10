package quest

import "fmt"

type TraitKind uint8

const (
	StringHash TraitKind = iota + 1
	Int64Hash
	Int32Hash
	BoolHash
)

func errInvalidType(val interface{}, required interface{}) error {
	return fmt.Errorf("quest: value %v (%T) is not %T", val, val, required)
}

func errUnsupportedValue(cond *Condition, required interface{}) error {
	return fmt.Errorf("quest: condition '%s' value %v (%T) is not %T", cond.Trait, cond.Value, cond.Value, required)
}

func errUnsupportedComparator(cond *Condition, kind string) error {
	return fmt.Errorf("quest: condition '%s' comparator '%s' is not supported by %s trait", cond.Trait, cond.Comparator, kind)
}

// --------------------------------------------------------------------

type stringHash map[string][]ruleReference

func (h stringHash) Find(val interface{}) ([]ruleReference, error) {
	vv, ok := val.(string)
	if !ok {
		return nil, errInvalidType(val, vv)
	}
	return h[vv], nil
}

func (h stringHash) Check(cond *Condition) error {
	if str, ok := cond.Value.(string); !ok {
		return errUnsupportedValue(cond, str)
	}
	if cond.Comparator != ComparatorEqual {
		return errUnsupportedComparator(cond, "StringHash")
	}
	return nil
}

func (h stringHash) Store(ref ruleReference, val interface{}) {
	vv := val.(string)
	h[vv] = append(h[vv], ref)
}

// --------------------------------------------------------------------

type int64Hash map[int64][]ruleReference

func (h int64Hash) Find(val interface{}) ([]ruleReference, error) {
	vv, ok := val.(int64)
	if !ok {
		return nil, errInvalidType(val, vv)
	}
	return h[vv], nil
}

func (h int64Hash) Check(cond *Condition) error {
	if vv, ok := cond.Value.(int64); !ok {
		return errUnsupportedValue(cond, vv)
	}
	if cond.Comparator != ComparatorEqual {
		return errUnsupportedComparator(cond, "Int64Hash")
	}
	return nil
}

func (h int64Hash) Store(ref ruleReference, val interface{}) {
	vv := val.(int64)
	h[vv] = append(h[vv], ref)
}

// --------------------------------------------------------------------

type int32Hash map[int32][]ruleReference

func (h int32Hash) Find(val interface{}) ([]ruleReference, error) {
	vv, ok := val.(int32)
	if !ok {
		return nil, errInvalidType(val, vv)
	}
	return h[vv], nil
}

func (h int32Hash) Check(cond *Condition) error {
	if vv, ok := cond.Value.(int32); !ok {
		return errUnsupportedValue(cond, vv)
	}
	if cond.Comparator != ComparatorEqual {
		return errUnsupportedComparator(cond, "Int32Hash")
	}
	return nil
}

func (h int32Hash) Store(ref ruleReference, val interface{}) {
	vv := val.(int32)
	h[vv] = append(h[vv], ref)
}

// --------------------------------------------------------------------

type boolHash map[bool][]ruleReference

func (h boolHash) Find(val interface{}) ([]ruleReference, error) {
	vv, ok := val.(bool)
	if !ok {
		return nil, errInvalidType(val, vv)
	}
	return h[vv], nil
}

func (h boolHash) Check(cond *Condition) error {
	if vv, ok := cond.Value.(bool); !ok {
		return errUnsupportedValue(cond, vv)
	}
	if cond.Comparator != ComparatorEqual {
		return errUnsupportedComparator(cond, "BoolHash")
	}
	return nil
}

func (h boolHash) Store(ref ruleReference, val interface{}) {
	vv := val.(bool)
	h[vv] = append(h[vv], ref)
}

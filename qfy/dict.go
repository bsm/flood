package qfy

import "sync"

// Dict is a simple helper to to turn convert strings into integers using
// dictionary encoding. Dict instances are NOT thread-safe because they are
// populated sequentially on Qualifier creation and then only read.
//
// If your use case requires rules to call Add operations concurrently
// please use (the slower but) thread-safe ConcurrentDict
type Dict map[string]int

// NewDict creates a new dict
func NewDict() Dict { return make(Dict) }

// Add adds a value to the dictionary if not already present and returns the associtated ID
func (d Dict) Add(v string) int {
	num, ok := d[v]
	if !ok {
		num = len(d) + 1
		d[v] = num
	}
	return num
}

// AddSlice adds a whole slice of values
func (d Dict) AddSlice(vv ...string) []int {
	nn := make([]int, len(vv))
	for i, v := range vv {
		nn[i] = d.Add(v)
	}
	return nn
}

// GetSlice is a read-only operation, only returns known IDs
func (d Dict) GetSlice(vv ...string) []int {
	nn := make([]int, 0, len(vv))
	for _, v := range vv {
		if n, ok := d[v]; ok {
			nn = append(nn, n)
		}
	}
	return nn
}

// ConcurrentDict is just like Dict, but thread-safe and therefore slightly
// slower.
type ConcurrentDict struct {
	dict  Dict
	mutex sync.RWMutex
}

// NewConcurrentDict creates a new dict
func NewConcurrentDict() *ConcurrentDict { return &ConcurrentDict{dict: NewDict()} }

// Add adds a value to the dictionary if not already present and returns the associtated ID
func (d *ConcurrentDict) Add(v string) int {
	d.mutex.Lock()
	num := d.dict.Add(v)
	d.mutex.Unlock()
	return num
}

// AddSlice adds a whole slice of values
func (d *ConcurrentDict) AddSlice(vv ...string) []int {
	nn := make([]int, len(vv))
	for i, v := range vv {
		nn[i] = d.Add(v)
	}
	return nn
}

// GetSlice is a read-only operation, only returns known IDs
func (d *ConcurrentDict) GetSlice(vv ...string) []int {
	d.mutex.RLock()
	nn := d.dict.GetSlice(vv...)
	d.mutex.RUnlock()
	return nn
}

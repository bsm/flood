package qfy

type strDict map[string]int

func (d strDict) Fetch(v string) int {
	num, ok := d[v]
	if !ok {
		num = len(d) + 1
		d[v] = num
	}
	return num
}

func (d strDict) FetchSlice(vv ...string) []int {
	nn := make([]int, len(vv))
	for i, v := range vv {
		nn[i] = d.Fetch(v)
	}
	return nn
}

func (d strDict) GetSlice(vv ...string) []int {
	nn := make([]int, 0, len(vv))
	for _, v := range vv {
		if n, ok := d[v]; ok {
			nn = append(nn, n)
		}
	}
	return nn
}

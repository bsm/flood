package qfy

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func BenchmarkQualifier(b *testing.B) {
	h, err := newBenchHelper()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		fact := h.fcts[i%h.size]
		h.q.Select(fact)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fact := h.fcts[i%h.size]
		h.q.Select(fact)
	}
}

// --------------------------------------------------------------------

func init() {
	// Populate benchFactKeyMap
	for i := 0; i < benchFactType.NumField(); i++ {
		rf := benchFactType.Field(i)
		name := rf.Tag.Get("json")
		name = name[:strings.Index(name, ",")]
		benchFactKeyMap[name] = FactKey(i)
	}

	fact := &benchFact{Dev: "ok", Vcat: []int{2, 6}}
	fact.GetQualifiable(benchFactKeyMap["vcat"])
}

var (
	benchFactKeyMap      = make(map[string]FactKey)
	benchFactType        = reflect.TypeOf(benchFact{})
	benchStringSliceType = reflect.TypeOf([]string{})
	benchIntSliceType    = reflect.TypeOf([]int{})
	benchDict            = NewDict()
)

type benchFact struct {
	Dev    string   `json:"dev,omitempty"`
	Tod    int      `json:"tod,omitempty"`
	Bwsm   string   `json:"bwsm,omitempty"`
	Ctry   string   `json:"ctry,omitempty"`
	Exch   int      `json:"exch,omitempty"`
	Klmn   int      `json:"klmn,omitempty"`
	Wr     int      `json:"wr,omitempty"`
	Pos    int      `json:"pos,omitempty"`
	Wl     int      `json:"wl,omitempty"`
	Isp    string   `json:"isp,omitempty"`
	Pcode  string   `json:"pcode,omitempty"`
	Mob    int      `json:"mob,omitempty"`
	Strm   string   `json:"strm,omitempty"`
	Ws     int      `json:"ws,omitempty"`
	Loc    string   `json:"loc,omitempty"`
	Domain string   `json:"domain,omitempty"`
	Infq   []string `json:"infq,omitempty"`
	Reg    string   `json:"reg,omitempty"`
	Ac     string   `json:"ac,omitempty"`
	Hb     int      `json:"hb,omitempty"`
	Kws    []string `json:"kws,omitempty"`
	Pmnt   string   `json:"pmnt,omitempty"`
	Vcat   []int    `json:"vcat,omitempty"`
	Rcat   []int    `json:"rcat,omitempty"`
}

func (bf *benchFact) GetQualifiable(key FactKey) []int {
	field := reflect.ValueOf(bf).Elem().Field(int(key))

	switch field.Kind() {
	case reflect.String:
		return benchDict.GetSlice(field.String())
	case reflect.Int:
		return []int{int(field.Int())}
	case reflect.Slice:
		switch field.Type() {
		case benchStringSliceType:
			return benchDict.GetSlice(field.Interface().([]string)...)
		case benchIntSliceType:
			return field.Interface().([]int)
		}
	}
	return nil
}

// --------------------------------------------------------------------

type benchTargetValues []json.Number

func (s benchTargetValues) Vals() []int {
	vals, err := s.Ints()
	if err != nil {
		vals = benchDict.AddSlice(s.Strings()...)
	}
	return vals
}

func (s benchTargetValues) Strings() []string {
	res := make([]string, len(s))
	for i, n := range s {
		res[i] = n.String()
	}
	return res
}

func (s benchTargetValues) Ints() ([]int, error) {
	res := make([]int, len(s))
	for i, n := range s {
		num, err := n.Int64()
		if err != nil {
			return nil, err
		}
		res[i] = int(num)
	}
	return res, nil
}

// --------------------------------------------------------------------

type benchHelper struct {
	q    *Qualifier
	fcts []Fact
	size int
}

func newBenchHelper() (*benchHelper, error) {
	h := &benchHelper{}
	if err := h.parseQualifier(); err != nil {
		return nil, err
	} else if err := h.parseFacts(); err != nil {
		return nil, err
	}
	return h, nil
}

func (bh *benchHelper) parseQualifier() error {
	file, err := os.Open("testdata/targets.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var targeting []struct {
		ID    int
		Rules []struct {
			Attr, Op string
			Values   benchTargetValues
		}
	}
	if err := json.NewDecoder(file).Decode(&targeting); err != nil {
		return err
	}

	bh.q = New()
	for _, target := range targeting {
		var rules []Rule
		for _, rdef := range target.Rules {
			key := benchFactKeyMap[rdef.Attr]
			vals := rdef.Values.Vals()

			if rdef.Op == "-" {
				rules = append(rules, key.MustBe(NoneOf(vals)))
			} else {
				rules = append(rules, key.MustBe(OneOf(vals)))
			}
		}
		bh.q.Resolve(All(rules...), target.ID)
	}
	return nil
}

func (bh *benchHelper) parseFacts() error {
	file, err := os.Open("testdata/facts.json")
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	for {
		var fact *benchFact
		err := dec.Decode(&fact)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		bh.fcts = append(bh.fcts, fact)
	}
	bh.size = len(bh.fcts)
	return nil
}

package qfy

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func BenchmarkQualifier(b *testing.B) {
	helper, err := newBenchHelper()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fact := helper.fcts[i%helper.size]
		helper.qlfy.Select(fact)
	}
}

// --------------------------------------------------------------------

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

func (bf *benchFact) GetQualifiable(name string) []int {
	switch name {
	case "dev":
		return benchDict.GetSlice(bf.Dev)
	case "tod":
		return []int{bf.Tod}
	case "bwsm":
		return benchDict.GetSlice(bf.Bwsm)
	case "ctry":
		return benchDict.GetSlice(bf.Ctry)
	case "exch":
		return []int{bf.Exch}
	case "klmn":
		return []int{bf.Klmn}
	case "wr":
		return []int{bf.Wr}
	case "pos":
		return []int{bf.Pos}
	case "wl":
		return []int{bf.Wl}
	case "isp":
		return benchDict.GetSlice(bf.Isp)
	case "pcode":
		return benchDict.GetSlice(bf.Pcode)
	case "mob":
		return []int{bf.Mob}
	case "strm":
		return benchDict.GetSlice(bf.Strm)
	case "ws":
		return []int{bf.Ws}
	case "loc":
		return benchDict.GetSlice(bf.Loc)
	case "domain":
		return benchDict.GetSlice(bf.Domain)
	case "rcat":
		return bf.Rcat
	case "vcat":
		return bf.Vcat
	case "infq":
		return benchDict.GetSlice(bf.Infq...)
	case "reg":
		return benchDict.GetSlice(bf.Reg)
	case "ac":
		return benchDict.GetSlice(bf.Ac)
	case "hb":
		return []int{bf.Hb}
	case "kws":
		return benchDict.GetSlice(bf.Kws...)
	case "pmnt":
		return benchDict.GetSlice(bf.Pmnt)
	}
	return nil
}

// --------------------------------------------------------------------

var benchAttrs = []string{
	"wl", "strm", "dev", "pos", "hb", "wr", "klmn", "mob",
	"ws", "loc", "isp", "pcode", "bwsm", "ctry", "tod", "ac",
	"pmnt", "exch", "reg", "infq", "rcat", "vcat", "kws", "domain",
}

var benchDict = NewDict()

// --------------------------------------------------------------------

type jsonSlice []json.Number

func (s jsonSlice) Strings() []string {
	res := make([]string, len(s))
	for i, n := range s {
		res[i] = n.String()
	}
	return res
}

func (s jsonSlice) Ints() ([]int, error) {
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
	qlfy *Qualifier
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

	var targets []struct {
		ID    int
		Rules []struct {
			Attr, Op string
			Values   jsonSlice
		}
	}
	if err := json.NewDecoder(file).Decode(&targets); err != nil {
		return err
	}

	bh.qlfy = New(benchAttrs)
	for _, target := range targets {
		byAttr := make(map[string][]Rule)
		for _, def := range target.Rules {
			if len(def.Values) == 0 {
				continue
			}
			vals, err := def.Values.Ints()
			if err != nil {
				vals = benchDict.AddSlice(def.Values.Strings()...)
			}
			if def.Op == "-" {
				byAttr[def.Attr] = append(byAttr[def.Attr], NoneOf(vals))
			} else {
				byAttr[def.Attr] = append(byAttr[def.Attr], OneOf(vals))
			}
		}

		rules := make(map[string]Rule, len(byAttr))
		for attr, rs := range byAttr {
			if len(rs) == 1 {
				rules[attr] = rs[0]
			} else {
				rules[attr] = All(rs...)
			}
		}
		bh.qlfy.Feed(target.ID, rules)
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

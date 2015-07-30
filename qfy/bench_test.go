package qfy

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

var benchattrs = []string{
	"wl", "strm", "dev", "pos", "hb", "wr", "klmn", "mob",
	"ws", "loc", "isp", "pcode", "bwsm", "ctry", "tod", "ac",
	"pmnt", "exch", "reg", "infq", "rcat", "vcat", "kws", "domain",
}

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
	Dev    []string `json:"dev,omitempty"`
	Tod    []int    `json:"tod,omitempty"`
	Bwsm   []string `json:"bwsm,omitempty"`
	Ctry   []string `json:"ctry,omitempty"`
	Exch   []int    `json:"exch,omitempty"`
	Klmn   []int    `json:"klmn,omitempty"`
	Wr     []int    `json:"wr,omitempty"`
	Pos    []int    `json:"pos,omitempty"`
	Wl     []int    `json:"wl,omitempty"`
	Isp    []string `json:"isp,omitempty"`
	Pcode  []string `json:"pcode,omitempty"`
	Mob    []int    `json:"mob,omitempty"`
	Strm   []string `json:"strm,omitempty"`
	Ws     []int    `json:"ws,omitempty"`
	Loc    []string `json:"loc,omitempty"`
	Domain []string `json:"domain,omitempty"`
	Infq   []string `json:"infq,omitempty"`
	Reg    []string `json:"reg,omitempty"`
	Ac     []string `json:"ac,omitempty"`
	Hb     []int    `json:"hb,omitempty"`
	Kws    []string `json:"kws,omitempty"`
	Pmnt   []string `json:"pmnt,omitempty"`
	Vcat   []int    `json:"vcat,omitempty"`
	Rcat   []int    `json:"rcat,omitempty"`
}

func (bf *benchFact) Get(name string) interface{} {
	switch name {
	case "dev":
		return bf.Dev
	case "tod":
		return bf.Tod
	case "bwsm":
		return bf.Bwsm
	case "ctry":
		return bf.Ctry
	case "exch":
		return bf.Exch
	case "klmn":
		return bf.Klmn
	case "wr":
		return bf.Wr
	case "pos":
		return bf.Pos
	case "wl":
		return bf.Wl
	case "isp":
		return bf.Isp
	case "pcode":
		return bf.Pcode
	case "mob":
		return bf.Mob
	case "strm":
		return bf.Strm
	case "ws":
		return bf.Ws
	case "loc":
		return bf.Loc
	case "domain":
		return bf.Domain
	case "rcat":
		return bf.Rcat
	case "vcat":
		return bf.Vcat
	case "infq":
		return bf.Infq
	case "reg":
		return bf.Reg
	case "ac":
		return bf.Ac
	case "hb":
		return bf.Hb
	case "kws":
		return bf.Kws
	case "pmnt":
		return bf.Pmnt
	}
	return nil
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
		Rules []RuleDef
	}
	if err := json.NewDecoder(file).Decode(&targets); err != nil {
		return err
	}

	bh.qlfy = New(benchattrs)
	for _, target := range targets {
		if err := bh.qlfy.Feed(target.ID, target.Rules); err != nil {
			return err
		}
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

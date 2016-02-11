package quest_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bsm/flood/quest"
)

func BenchmarkQuest(b *testing.B) {
	subject := quest.New()

	if err := subject.RegisterTrait("dev", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("tod", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("bwsm", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("ctry", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("exch", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("klmn", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("wr", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("pos", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("wl", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("isp", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("pcode", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("mob", quest.BoolHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("strm", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("ws", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("loc", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("domain", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("infq", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("reg", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("ac", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("hb", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("kws", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("pmnt", quest.StringHash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("vcat", quest.Int64Hash); err != nil {
		b.Fatal(err)
	} else if err := subject.RegisterTrait("rcat", quest.Int64Hash); err != nil {
		b.Fatal(err)
	}

	if err := loadRules(subject); err != nil {
		b.Fatal(err)
	}

	facts, err := loadFacts()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fact := facts[i%len(facts)]
		if _, err := subject.Match(&fact); err != nil {
			b.Fatal(err)
		}
	}
}

func loadFacts() (facts []benchFact, _ error) {
	file, err := os.Open("testdata/facts.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&facts); err != nil {
		return nil, err
	}
	return facts, nil
}

func loadRules(q *quest.Quest) error {
	file, err := os.Open("../qfy/testdata/targets.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var outcomes []struct {
		ID    quest.Outcome
		Rules []struct {
			Attr, Op string
			Values   []json.Number
		}
	}
	if err := json.NewDecoder(file).Decode(&outcomes); err != nil {
		return err
	}

	for _, odef := range outcomes {
		for _, rdef := range odef.Rules {
			rule := quest.Rule{Negation: rdef.Op == "-"}

			for _, val := range rdef.Values {
				cond := quest.Condition{
					Trait:      rdef.Attr,
					Comparator: quest.ComparatorEqual,
				}
				switch rdef.Attr {
				case "mob":
					cond.Value = val.String() == "true"
				case "tod", "exch", "klmn", "wr", "pos", "wl", "ws", "hb", "vcat", "rcat":
					cond.Value, _ = val.Int64()
				default:
					cond.Value = val.String()
				}
				rule.Conditions = append(rule.Conditions, cond)
			}
			if err := q.AddRule(odef.ID, &rule); err != nil {
				return err
			}
		}
	}
	return nil
}

// --------------------------------------------------------------------

type benchFact struct {
	Dev    string
	Tod    int64
	Bwsm   string
	Ctry   string
	Exch   int64
	Klmn   int64
	Wr     int64
	Pos    int64
	Wl     int64
	Isp    string
	Pcode  string
	Mob    bool
	Strm   string
	Ws     int64
	Loc    string
	Domain string
	Infq   []string
	Reg    string
	Ac     string
	Hb     int64
	Kws    []string
	Pmnt   string
	Vcat   []int64
	Rcat   []int64
}

func (bf *benchFact) GetFactValue(name string) interface{} {
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
	case "vcat":
		return bf.Vcat
	case "rcat":
		return bf.Rcat
	}
	return nil
}

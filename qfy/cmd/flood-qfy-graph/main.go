package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bsm/flood/qfy"
)

type targeting struct {
	ID    int `json:"id"`
	Rules []qfy.RuleDef
}

var argv struct {
	input, attrs string
	limit        int
}

func init() {
	flag.StringVar(&argv.input, "in", "targets.json", "Input file path")
	flag.StringVar(&argv.attrs, "attrs", "dev,pos,hb,wr,wl,tod,bwsm,ctry,reg,exch,klmn,isp,pcode,mob,strm,ws,loc,rcat,vcat,infq,ac,kws,domain,pmnt", "Comma-separated list of attributes")
	flag.IntVar(&argv.limit, "n", -1, "Limit output to N targets")
}

func abortOn(err error) {
	if err != nil {
		fmt.Println("ERROR", err.Error())
		os.Exit(2)
	}
}

func main() {
	flag.Parse()

	targets, err := readTargetsFile()
	abortOn(err)

	qualifier := qfy.New(strings.Split(argv.attrs, ","))
	for _, tdef := range targets {
		qualifier.Feed(tdef.ID, tdef.Rules)
	}
	qualifier.Graph(os.Stdout)
}

// --------------------------------------------------------------------

func readTargetsFile() ([]targeting, error) {
	var res []targeting

	input, err := os.Open(argv.input)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	if err := json.NewDecoder(input).Decode(&res); err != nil {
		return nil, err
	}
	if argv.limit > 0 && len(res) > argv.limit {
		res = res[:argv.limit]
	}
	return res, nil
}

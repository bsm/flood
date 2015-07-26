package qfy

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/bsm/intset"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Qualifier", func() {
	var subject *Qualifier

	BeforeEach(func() {
		subject = New([]Attribute{
			{"unused", TypeStringSlice},
			{"country", TypeStringSlice},
			{"browser", TypeStringSlice},
			{"domains", TypeStringSlice},
		})

		err := subject.Feed(91, []RuleDef{
			{Attr: "country", Op: "+", Vals: json.RawMessage(`["US"]`)},
			{Attr: "browser", Op: "-", Vals: json.RawMessage(`["IE"]`)},
			{Attr: "domains", Op: "+", Vals: json.RawMessage(`["a.com", "b.com"]`)},
			{Attr: "domains", Op: "-", Vals: json.RawMessage(`["c.com"]`)},
		})
		Expect(err).NotTo(HaveOccurred())

		err = subject.Feed(92, []RuleDef{
			{Attr: "country", Op: "-", Vals: json.RawMessage(`["CA"]`)},
			{Attr: "domains", Op: "+", Vals: json.RawMessage(`["b.com", "c.com"]`)},
			{Attr: "domains", Op: "+", Vals: json.RawMessage(`["d.com", "a.com"]`)},
		})
		Expect(err).NotTo(HaveOccurred())

		err = subject.Feed(93, []RuleDef{
			{Attr: "country", Op: "+", Vals: json.RawMessage(`["US"]`)},
			{Attr: "browser", Op: "-", Vals: json.RawMessage(`["OP"]`)},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("should feed with data", func() {
		Expect(subject.root.children).To(HaveLen(1))
		Expect(subject.root.children[0].(*passNode).children).To(HaveLen(2))
	})

	It("should reject bad feed inputs", func() {
		Expect(subject.Feed(96, []RuleDef{
			{Attr: "country", Op: "*", Vals: json.RawMessage(`["US"]`)},
		})).To(HaveOccurred())

		Expect(subject.Feed(96, []RuleDef{
			{Attr: "country", Op: "+", Vals: json.RawMessage(`{"a":1}`)},
		})).To(HaveOccurred())

		Expect(subject.Feed(96, []RuleDef{
			{Attr: "country", Op: "-", Vals: json.RawMessage(`[]`)},
		})).To(HaveOccurred())
	})

	It("should graph", func() {
		w := &bytes.Buffer{}
		subject.Graph(w)
		Expect(w.String()).To(ContainSubstring("\tN0000000000000000 [label = \"[root]\"]\n"))
	})

	It("should match", func() {
		tests := []struct {
			fact Fact
			vals []int
		}{
			{&mockFactStruct{}, []int{}},
			{&mockFactStruct{Country: []string{"US"}}, []int{93}},
			{&mockFactStruct{Country: []string{"US"}, Domains: []string{"a.com", "d.com"}}, []int{91, 93}},
			{&mockFactStruct{Country: []string{"US"}, Domains: []string{"a.com", "c.com", "d.com"}}, []int{92, 93}},
			{&mockFactStruct{Country: []string{"CA"}, Domains: []string{"a.com", "c.com", "d.com"}}, []int{}},
			{&mockFactStruct{Domains: []string{"a.com", "c.com", "d.com"}}, []int{92}},
			{&mockFactStruct{Country: []string{"US"}, Browser: []string{"OP"}, Domains: []string{"b.com", "x.com"}}, []int{91}},
		}

		for _, test := range tests {
			res := subject.Select(test.fact)
			Expect(res).To(ConsistOf(test.vals), "for %+v", test.fact)
		}
	})

})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "flood/qfy")
}

type mockConverter struct{}

func (mockConverter) convert(v interface{}) *intset.Set {
	return intset.Use(v.([]int)...)
}

type mockFact map[string][]int

func (m mockFact) Get(attr string) interface{} { return m[attr] }

type mockFactStruct struct{ Country, Browser, Domains []string }

func (m mockFactStruct) Get(attr string) interface{} {
	switch attr {
	case "country":
		return m.Country
	case "browser":
		return m.Browser
	case "domains":
		return m.Domains
	}
	return nil
}

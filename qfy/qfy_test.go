package qfy

import (
	"bytes"
	"testing"

	"github.com/bsm/intset"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Qualifier", func() {
	var subject *Qualifier
	var dict Dict

	BeforeEach(func() {
		dict = NewDict()
		subject = New([]string{"unused", "country", "browser", "domains"})

		subject.Feed(91, map[string]Rule{
			"country": OneOf(dict.AddSlice("US")),
			"browser": NoneOf(dict.AddSlice("IE")),
			"domains": All(
				OneOf(dict.AddSlice("a.com", "b.com")),
				NoneOf(dict.AddSlice("c.com")),
			),
		})

		subject.Feed(92, map[string]Rule{
			"country": NoneOf(dict.AddSlice("CA")),
			"domains": All(
				OneOf(dict.AddSlice("b.com", "c.com")),
				OneOf(dict.AddSlice("d.com", "a.com")),
			),
		})

		subject.Feed(93, map[string]Rule{
			"country": OneOf(dict.AddSlice("US")),
			"browser": NoneOf(dict.AddSlice("OP")),
		})
	})

	It("should feed with data", func() {
		Expect(subject.root.children).To(HaveLen(1))
		Expect(subject.root.children[0].(*passNode).children).To(HaveLen(2))
	})

	It("should graph", func() {
		w := &bytes.Buffer{}
		subject.Graph(w)
		Expect(w.String()).To(ContainSubstring("\tN0000000000000000 [label = \"[root]\"]\n"))
	})

	DescribeTable("matching",
		func(fact *mockFactStruct, expected []int) {
			fact.D = dict // assign dict
			Expect(subject.Select(fact)).To(ConsistOf(expected))
		},

		Entry("blank",
			&mockFactStruct{}, []int{}),
		Entry("91 & 92 have domain inclusions, 93 matches",
			&mockFactStruct{Country: "US"}, []int{93}),
		Entry("91 & 93 match, 92 has only one matching domain rule",
			&mockFactStruct{Country: "US", Domains: []string{"a.com", "d.com"}}, []int{91, 93}),
		Entry("92 & 93 match, 91 excludes c.com",
			&mockFactStruct{Country: "US", Domains: []string{"a.com", "c.com", "d.com"}}, []int{92, 93}),
		Entry("91 & 93 require US, 92 excludes CA",
			&mockFactStruct{Country: "CA", Domains: []string{"a.com", "c.com", "d.com"}}, []int{}),
		Entry("91 & 93 have explicit country inclusions, 92 matches",
			&mockFactStruct{Domains: []string{"a.com", "c.com", "d.com"}}, []int{92}),
		Entry("92 requires more domains, 93 excludes OP, 91 matches",
			&mockFactStruct{Country: "US", Browser: "OP", Domains: []string{"b.com", "x.com"}}, []int{91}),
	)

})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "flood/qfy")
}

type mockFact map[string][]int

func (m mockFact) Get(attr string) *intset.Set {
	if vv, ok := m[attr]; ok {
		return intset.Use(vv...)
	}
	return nil
}

type mockFactStruct struct {
	D Dict

	Country, Browser string
	Domains          []string
}

func (m *mockFactStruct) Get(attr string) *intset.Set {
	switch attr {
	case "country":
		return intset.Use(m.D.GetSlice(m.Country)...)
	case "browser":
		return intset.Use(m.D.GetSlice(m.Browser)...)
	case "domains":
		return intset.Use(m.D.GetSlice(m.Domains...)...)
	}
	return nil
}

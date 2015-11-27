package qfy

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Qualifier", func() {
	var subject *Qualifier
	var dict Dict

	BeforeEach(func() {
		dict = NewDict()
		subject = New()

		subject.Resolve(All(
			mockFactCountry.MustBe(EqualTo(dict.Add("US"))),
			mockFactBrowser.MustNotBe(EqualTo(dict.Add("IE"))),
			mockFactDomains.MustInclude(OneOf(dict.AddSlice("a.com", "b.com"))),
			mockFactDomains.MustInclude(NoneOf(dict.AddSlice("c.com"))),
		), 91)

		subject.Resolve(All(
			mockFactCountry.MustNotBe(EqualTo(dict.Add("CA"))),
			mockFactDomains.MustInclude(OneOf(dict.AddSlice("b.com", "c.com"))),
			mockFactDomains.MustInclude(OneOf(dict.AddSlice("d.com", "a.com"))),
		), 92)

		subject.Resolve(All(
			mockFactCountry.MustBe(OneOf(dict.AddSlice("US"))),
			mockFactBrowser.MustBe(NoneOf(dict.AddSlice("OP"))),
		), 93)
	})

	It("should register targets", func() {
		Expect(subject.registry).To(HaveLen(3))
	})

	DescribeTable("matching",
		func(fact *mockFactStruct, expected []int64) {
			fact.D = dict // assign dict
			Expect(subject.Select(fact)).To(ConsistOf(expected))
		},

		Entry("blank",
			&mockFactStruct{}, []int64{}),
		Entry("91 & 92 have domain inclusions, 93 matches",
			&mockFactStruct{Country: "US"}, []int64{93}),
		Entry("91 & 93 match, 92 has only one matching domain rule",
			&mockFactStruct{Country: "US", Domains: []string{"a.com", "d.com"}}, []int64{91, 93}),
		Entry("92 & 93 match, 91 excludes c.com",
			&mockFactStruct{Country: "US", Domains: []string{"a.com", "c.com", "d.com"}}, []int64{92, 93}),
		Entry("91 & 93 require US, 92 excludes CA",
			&mockFactStruct{Country: "CA", Domains: []string{"a.com", "c.com", "d.com"}}, []int64{}),
		Entry("91 & 93 have explicit country inclusions, 92 matches",
			&mockFactStruct{Domains: []string{"a.com", "c.com", "d.com"}}, []int64{92}),
		Entry("92 requires more domains, 93 excludes OP, 91 matches",
			&mockFactStruct{Country: "US", Browser: "OP", Domains: []string{"b.com", "x.com"}}, []int64{91}),
	)

})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "flood/qfy")
}

// --------------------------------------------------------------------

const (
	mockFactCountry FactKey = iota
	mockFactBrowser
	mockFactDomains
)

type mockFact map[FactKey][]int64

func (m mockFact) GetQualifiable(key FactKey) interface{} {
	vv, _ := m[key]
	return vv
}

type mockFactStruct struct {
	D Dict

	Country, Browser string
	Domains          []string
}

func (m *mockFactStruct) GetQualifiable(key FactKey) interface{} {
	switch key {
	case mockFactCountry:
		return m.D.Get(m.Country)
	case mockFactBrowser:
		return m.D.Get(m.Browser)
	case mockFactDomains:
		return m.D.GetSlice(m.Domains...)
	}
	return nil
}

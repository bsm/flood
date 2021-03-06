package quest

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Quest", func() {
	var subject *Quest
	var seed = func(o Outcome, rs []Rule) error {
		for _, r := range rs {
			if err := subject.AddRule(o, &r); err != nil {
				return err
			}
		}
		return nil
	}

	BeforeEach(func() {
		subject = New()

		err := subject.RegisterTrait("country", StringHash)
		Expect(err).NotTo(HaveOccurred())
		err = subject.RegisterTrait("browser", StringHash)
		Expect(err).NotTo(HaveOccurred())
		err = subject.RegisterTrait("os", StringHash)
		Expect(err).NotTo(HaveOccurred())

		Expect(seed(91, []Rule{
			{Conditions: []Condition{
				{"country", ComparatorEqual, "GB"},
			}},
			{Conditions: []Condition{
				{"browser", ComparatorEqual, "firefox"},
				{"browser", ComparatorEqual, "chrome"},
				{"browser", ComparatorEqual, "safari"},
			}},
		})).NotTo(HaveOccurred())

		Expect(seed(92, []Rule{
			{Negation: true, Conditions: []Condition{
				{"country", ComparatorEqual, "US"},
				{"country", ComparatorEqual, "CA"},
			}},
			{Conditions: []Condition{
				{"browser", ComparatorEqual, "safari"},
				{"os", ComparatorEqual, "ios"},
			}},
		})).NotTo(HaveOccurred())

		Expect(seed(93, []Rule{
			{Conditions: []Condition{
				{"country", ComparatorEqual, "GB"},
			}},
			{Negation: true, Conditions: []Condition{
				{"browser", ComparatorEqual, "chrome"},
			}},
			{Negation: true, Conditions: []Condition{
				{"os", ComparatorEqual, "linux"},
			}},
		})).NotTo(HaveOccurred())
	})

	It("should register traits", func() {
		Expect(subject.traits).To(HaveLen(3))
		err := subject.RegisterTrait("country", StringHash)
		Expect(err).To(MatchError(`quest: trait 'country' is already regitered`))
	})

	It("should append rules/outcomes", func() {
		Expect(subject.outcomes).To(HaveLen(3))

		err := subject.AddRule(99, &Rule{
			Conditions: []Condition{{"unknown", ComparatorEqual, "value"}},
		})
		Expect(err).To(MatchError(`quest: condition references unknown trait 'unknown'`))

		err = subject.AddRule(99, &Rule{
			Conditions: []Condition{{"browser", ComparatorGreater, "value"}},
		})
		Expect(err).To(MatchError(`quest: condition 'browser' comparator '>' is not supported by StringHash trait`))

		err = subject.AddRule(99, &Rule{
			Conditions: []Condition{{"browser", ComparatorEqual, 2}},
		})
		Expect(err).To(MatchError(`quest: condition 'browser' value 2 (int) is not string`))
	})

	DescribeTable("should match outcomes",
		func(fact Fact, expected []Outcome) {
			res, err := subject.Match(fact)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(ConsistOf(expected))
		},
		Entry("country and browser match 91",
			&mockFact{Country: "GB", Browser: "chrome", OS: "linux"},
			[]Outcome{91}),
		Entry("country and browser match 91 + 92",
			&mockFact{Country: "GB", Browser: "safari", OS: "linux"},
			[]Outcome{91, 92}),
		Entry("country and browser match 91, OS matches 92",
			&mockFact{Country: "GB", Browser: "firefox", OS: "ios"},
			[]Outcome{91, 92, 93}),
		Entry("browser/country doesn not match 91, OS matches 92",
			&mockFact{Country: "DE", Browser: "ie", OS: "ios"},
			[]Outcome{92}),
		Entry("country excluded",
			&mockFact{Country: "US", Browser: "firefox", OS: "ios"},
			[]Outcome{}),
		Entry("multiple exclusions",
			&mockFact{Country: "GB", Browser: "firefox", OS: "osx"},
			[]Outcome{91, 93}),
	)

})

// --------------------------------------------------------------------

type mockFact struct {
	Country, Browser, OS string
}

func (f *mockFact) GetFactValue(name string) interface{} {
	switch name {
	case "country":
		return f.Country
	case "browser":
		return f.Browser
	case "os":
		return f.OS
	}
	return nil
}

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "flood/quest")
}

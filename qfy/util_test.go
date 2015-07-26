package qfy

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuleDef", func() {
	var strs, ints *RuleDef

	BeforeEach(func() {
		strs = &RuleDef{Attr: "x", Op: "+", Vals: json.RawMessage(`["a", "b", "c"]`)}
		ints = &RuleDef{Attr: "y", Op: "+", Vals: json.RawMessage(`[1, 2, 3]`)}
	})

	It("should build string rules", func() {
		dict := strDict{}
		rule, err := strs.toRule(TypeStringSlice, dict)
		Expect(err).NotTo(HaveOccurred())
		Expect(rule).NotTo(BeNil())
		Expect(dict).To(HaveLen(3))
	})

	It("should build int rules", func() {
		dict := strDict{}
		rule, err := ints.toRule(TypeIntSlice, dict)
		Expect(err).NotTo(HaveOccurred())
		Expect(rule).NotTo(BeNil())
		Expect(dict).To(BeEmpty())
	})

	It("should ignore bad operators", func() {
		_, err := (&RuleDef{Attr: "x", Op: "*", Vals: json.RawMessage(`[1]`)}).toRule(TypeIntSlice, strDict{})
		Expect(err).To(HaveOccurred())
	})

	It("should ignore bad values", func() {
		_, err := (&RuleDef{Attr: "x", Op: "+", Vals: json.RawMessage(`{"a":1}`)}).toRule(TypeIntSlice, strDict{})
		Expect(err).To(HaveOccurred())

		_, err = (&RuleDef{Attr: "x", Op: "+", Vals: json.RawMessage(`[]`)}).toRule(TypeIntSlice, strDict{})
		Expect(err).To(HaveOccurred())

		_, err = ints.toRule(TypeStringSlice, strDict{})
		Expect(err).To(HaveOccurred())
	})

})

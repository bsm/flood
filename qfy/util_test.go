package qfy

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuleDef", func() {

	It("should detect type", func() {
		tests := []struct {
			t AttrType
			v interface{}
		}{
			{TypeIntSlice, []int{1, 2, 3}},
			{TypeIntSlice, 1},

			{TypeStringSlice, []string{"a", "b"}},
			{TypeStringSlice, "b"},
		}

		for _, test := range tests {
			rdef := &RuleDef{Attr: "x", Op: "+", Vals: test.v}
			kind, err := rdef.DetectType()
			Expect(err).NotTo(HaveOccurred(), "for %+v", rdef)
			Expect(kind).To(Equal(test.t), "for %+v", rdef)
		}
	})

	It("should convert to rules", func() {
		tests := []struct {
			v interface{}
			n int
			c int
		}{
			{[]string{"a", "b", "c"}, 3, 3},
			{"a", 1, 1},
			{[]int{1, 2, 3}, 3, 0},
			{1, 1, 0},
		}

		for _, test := range tests {
			rdef := &RuleDef{Attr: "x", Op: "+", Vals: test.v}
			dict := strDict{}
			rule, err := rdef.toRule(dict)
			Expect(err).NotTo(HaveOccurred(), "for %+v", rdef)
			Expect(rule.(*plusRule).vals.Len()).To(Equal(test.n), "for %+v", rdef)
			Expect(dict).To(HaveLen(test.c), "for %+v", rdef)
		}
	})

	It("should ignore bad inputs", func() {
		tests := []struct {
			o string
			v interface{}
		}{
			{"*", []int{1}},
			{"+", map[string]int{"a": 1}},
			{"+", []int{}},
			{"+", nil},
		}

		for _, test := range tests {
			rdef := &RuleDef{Attr: "x", Op: test.o, Vals: test.v}
			_, err := rdef.toRule(strDict{})
			Expect(err).To(HaveOccurred(), "for %+v", rdef)
		}
	})

	It("should correctly unmarshal from JSON", func() {
		tests := []struct {
			s string
			r *RuleDef
		}{
			{`{"attr":"x","op":"+","values":[1,2,3]}`, &RuleDef{Attr: "x", Op: "+", Vals: []int{1, 2, 3}}},
			{`{"attr":"x","op":"+","values":1}`, &RuleDef{Attr: "x", Op: "+", Vals: 1}},
			{`{"attr":"x","op":"+","values":["a","b"]}`, &RuleDef{Attr: "x", Op: "+", Vals: []string{"a", "b"}}},
			{`{"attr":"x","op":"+","values":[ "a","b"]}`, &RuleDef{Attr: "x", Op: "+", Vals: []string{"a", "b"}}},
			{`{"attr":"x","op":"+","values":[` + "\n" + `"a","b"]}`, &RuleDef{Attr: "x", Op: "+", Vals: []string{"a", "b"}}},
			{`{"attr":"x","op":"+","values":"a"}`, &RuleDef{Attr: "x", Op: "+", Vals: "a"}},

			{`{"attr":"x","op":"+","values":[]}`, nil},
			{`{"attr":"x","op":"+","values":[1, "a"]}`, nil},
			{`{"attr":"x","op":"+","values":{"a":1}}`, nil},
		}

		for _, test := range tests {
			var dst *RuleDef
			err := json.Unmarshal([]byte(test.s), &dst)
			if test.r == nil {
				Expect(err).To(HaveOccurred(), "for '%s'", test.s)
			} else {
				Expect(err).NotTo(HaveOccurred(), "for '%s'", test.s)
				Expect(dst).To(Equal(test.r), "for '%s'", test.s)
			}
		}
	})

})

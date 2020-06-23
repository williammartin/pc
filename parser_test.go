package pc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/williammartin/pc"
)

var _ = Describe("PC", func() {

	Describe("parseA", func() {
		When("input is empty", func() {
			It("returns false", func() {
				succeeded, _ := pc.ParseA("")
				Expect(succeeded).To(BeFalse())
			})
		})

		When("input matches", func() {
			It("returns remaining string and true", func() {
				succeeded, remaining := pc.ParseA("Abc")
				Expect(succeeded).To(BeTrue())
				Expect(remaining).To(Equal("bc"))
			})
		})

		When("input doesn't match", func() {
			It("returns the same string and false", func() {
				succeeded, remaining := pc.ParseA("foo")
				Expect(succeeded).To(BeFalse())
				Expect(remaining).To(Equal("foo"))
			})
		})
	})

})

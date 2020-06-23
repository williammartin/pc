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

	Describe("parseChar", func() {
		When("input is empty", func() {
			It("returns an error", func() {
				_, _, err := pc.ParseChar('z', "")
				Expect(err).To(MatchError("no more input"))
			})
		})

		When("input matches", func() {
			It("returns matched char and remaining string", func() {
				char, remaining, err := pc.ParseChar('z', "zyx")
				Expect(err).NotTo(HaveOccurred())
				Expect(char).To(BeEquivalentTo('z'))
				Expect(remaining).To(Equal("yx"))
			})
		})

		When("input doesn't match", func() {
			It("returns an error", func() {
				_, _, err := pc.ParseChar('z', "foo")
				Expect(err).To(MatchError("Expected 'z'. Got 'f'"))
			})
		})
	})

})

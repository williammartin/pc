package pc_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/williammartin/pc"
)

var _ = Describe("PC", func() {

	Describe("charParser", func() {
		It("returns a function that parses chars", func() {
			parseA := pc.CharParser("a")

			_, _, err := parseA("")
			Expect(err).To(MatchError("no more input"))

			char, remaining, err := parseA("abc")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(BeEquivalentTo("a"))
			Expect(remaining).To(Equal("bc"))

			_, _, err = parseA("zyx")
			Expect(err).To(MatchError("Expected 'a'. Got 'z'"))
		})

		It("matches any of the chars given", func() {
			parseABC := pc.CharParser("abc")

			char, remaining, err := parseABC("cxyz")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(BeEquivalentTo("c"))
			Expect(remaining).To(Equal("xyz"))
		})
	})

	Describe("andThen", func() {
		It("errors if first parse doesn't match", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.AndThen(parseA, parseB)

			_, _, err := parseAB("xyz")
			Expect(err).To(MatchError("Expected 'a'. Got 'x'"))
		})

		It("errors if second parse doesn't match", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.AndThen(parseA, parseB)

			_, _, err := parseAB("ayz")
			Expect(err).To(MatchError("Expected 'b'. Got 'y'"))
		})

		It("returns both chars and remaining if both match", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.AndThen(parseA, parseB)

			chars, remaining, err := parseAB("abc")
			Expect(err).NotTo(HaveOccurred())
			Expect(chars).To(Equal("ab"))
			Expect(remaining).To(Equal("c"))
		})
	})

	Describe("orElse", func() {
		It("returns the first char and remaining if first parser matches", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.OrElse(parseA, parseB)

			char, remaining, err := parseAB("axy")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(Equal("a"))
			Expect(remaining).To(Equal("xy"))
		})

		It("returns the second char and remaining if second parser matches", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.OrElse(parseA, parseB)

			char, remaining, err := parseAB("bxy")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(Equal("b"))
			Expect(remaining).To(Equal("xy"))
		})

		It("returns the error from the second parser if both parses don't match", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.OrElse(parseA, parseB)

			_, _, err := parseAB("xyz")
			Expect(err).To(MatchError("Expected 'b'. Got 'x'"))
		})
	})

	Describe("map", func() {
		It("returns the mapped parsed value if the parser matches", func() {
			parseA := pc.CharParser("a")
			parseAMap := pc.Map(strings.ToUpper, parseA)

			char, remaining, err := parseAMap("abc")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(Equal("A"))
			Expect(remaining).To(Equal("bc"))
		})

		It("errors if parse doesn't match", func() {
			parseA := pc.CharParser("a")
			parseAMap := pc.Map(strings.ToUpper, parseA)

			_, _, err := parseAMap("xyz")
			Expect(err).To(MatchError("Expected 'a'. Got 'x'"))
		})
	})
})

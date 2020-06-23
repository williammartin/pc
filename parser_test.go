package pc_test

import (
	"strings"

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
				_, _, err := pc.ParseChar("z", "")
				Expect(err).To(MatchError("no more input"))
			})
		})

		When("input matches", func() {
			It("returns matched char and remaining string", func() {
				char, remaining, err := pc.ParseChar("z", "zyx")
				Expect(err).NotTo(HaveOccurred())
				Expect(char).To(BeEquivalentTo("z"))
				Expect(remaining).To(Equal("yx"))
			})
		})

		When("input doesn't match", func() {
			It("returns an error", func() {
				_, _, err := pc.ParseChar("z", "foo")
				Expect(err).To(MatchError("Expected 'z'. Got 'f'"))
			})
		})
	})

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

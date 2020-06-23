package pc_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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
			parseAB := pc.AndThen(parseA, parseB, "ab")

			_, _, err := parseAB("xyz")
			Expect(err).To(MatchError("first group of 'ab' failed: Expected 'a'. Got 'x'"))
		})

		It("errors if second parse doesn't match", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.AndThen(parseA, parseB, "ab")

			_, _, err := parseAB("ayz")
			Expect(err).To(MatchError("second group of 'ab' failed: Expected 'b'. Got 'y'"))
		})

		It("returns both chars and remaining if both match", func() {
			parseA := pc.CharParser("a")
			parseB := pc.CharParser("b")
			parseAB := pc.AndThen(parseA, parseB, "ab")

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

	Describe("oneOrMore", func() {
		It("matches one or more of the parser", func() {
			parseA := pc.CharParser("A")
			parseAAA := pc.OneOrMore(parseA, "a+")

			char, remaining, err := parseAAA("AAABBB")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(Equal("AAA"))
			Expect(remaining).To(Equal("BBB"))
		})
	})

	Describe("zeroOrone", func() {
		It("returns the matched char and remaining when there is a match", func() {
			parseA := pc.CharParser("A")
			parseAMaybe := pc.ZeroOrOne(parseA)

			char, remaining, err := parseAMaybe("AB")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(Equal("A"))
			Expect(remaining).To(Equal("B"))
		})

		It("returns the input when there is no match", func() {
			parseA := pc.CharParser("A")
			parseAMaybe := pc.ZeroOrOne(parseA)

			char, remaining, err := parseAMaybe("B")
			Expect(err).NotTo(HaveOccurred())
			Expect(char).To(Equal(""))
			Expect(remaining).To(Equal("B"))
		})
	})

	Describe("a json structure", func() {

		Describe("parsing numbers", func() {
			var parseNumber pc.ParseFn

			BeforeEach(func() {
				parseOneToNine := pc.CharParser("123456789")
				parseZero := pc.CharParser("0")
				parseDecimalPlace := pc.CharParser(".")
				parseDigit := pc.OrElse(parseZero, parseOneToNine)
				parseDigits := pc.OneOrMore(parseDigit, "digits")
				parseDigits0 := pc.OrElse(parseDigits, pc.ZeroOrOne(parseDigit))
				parseFraction := pc.AndThen(
					pc.AndThen(parseZero, parseDecimalPlace, "zero and decimal"),
					parseDigits,
					"fraction",
				)
				parseNumber = pc.OrElse(
					parseFraction,
					pc.AndThen(
						pc.AndThen(parseOneToNine, parseDigits0, "1-9 followed by digits"),
						pc.ZeroOrOne(pc.AndThen(parseDecimalPlace, parseDigits, "decimal followed by digit")),
						"number greater than one",
					),
				)
			})

			It("matches a single digit number", func() {
				char, remaining, err := parseNumber("1")
				Expect(err).NotTo(HaveOccurred())
				Expect(char).To(Equal("1"))
				Expect(remaining).To(Equal(""))
			})

			DescribeTable("numbers", func(input, eChar, eRemaining string) {
				char, remaining, err := parseNumber(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(char).To(Equal(eChar))
				Expect(remaining).To(Equal(eRemaining))
			},
				Entry("single digit", "1", "1", ""),
				Entry("multiple digits more than one", "10089", "10089", ""),
				Entry("zero fraction", "0.123", "0.123", ""),
				Entry("fraction greater than one", "9230000.00", "9230000.00", ""),
			)
		})
	})
})

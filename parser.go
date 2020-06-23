package pc

import (
	"errors"
	"fmt"
)

func ParseA(input string) (bool, string) {
	if input == "" {
		return false, ""
	}

	if input[0] == 'A' {
		return true, input[1:]
	}

	return false, input
}

func ParseChar(char byte, input string) (byte, string, error) {
	if input == "" {
		return 0, "", errors.New("no more input")
	}

	if input[0] == char {
		return char, input[1:], nil
	}

	return 0, "", fmt.Errorf("Expected '%c'. Got '%c'", char, input[0])
}

type ParseFn func(string) (byte, string, error)

func CharParser(char byte) ParseFn {
	return func(input string) (byte, string, error) {
		if input == "" {
			return 0, "", errors.New("no more input")
		}

		if input[0] == char {
			return char, input[1:], nil
		}

		return 0, "", fmt.Errorf("Expected '%c'. Got '%c'", char, input[0])
	}
}

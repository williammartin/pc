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

func ParseChar(char string, input string) (string, string, error) {
	if input == "" {
		return "", "", errors.New("no more input")
	}

	if string(input[0]) == char {
		return char, input[1:], nil
	}

	return "", "", fmt.Errorf("Expected '%s'. Got '%c'", char, input[0])
}

type ParseFn func(string) (string, string, error)

func CharParser(char string) ParseFn {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("no more input")
		}

		if string(input[0]) == char {
			return char, input[1:], nil
		}

		return "", "", fmt.Errorf("Expected '%s'. Got '%c'", char, input[0])
	}
}

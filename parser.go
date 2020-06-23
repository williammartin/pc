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

func AndThen(first ParseFn, second ParseFn) ParseFn {
	return func(input string) (string, string, error) {
		firstChar, firstRemaining, firstErr := first(input)
		if firstErr != nil {
			return "", "", firstErr
		}

		secondChar, secondRemaining, secondErr := second(firstRemaining)
		if secondErr != nil {
			return "", "", secondErr
		}

		return firstChar + secondChar, secondRemaining, nil
	}
}

func OrElse(first ParseFn, second ParseFn) ParseFn {
	return func(input string) (string, string, error) {
		firstChar, firstRemaining, firstErr := first(input)
		if firstErr == nil {
			return firstChar, firstRemaining, nil
		}

		secondChar, secondRemaining, secondErr := second(input)
		if secondErr == nil {
			return secondChar, secondRemaining, nil
		}

		return "", "", secondErr
	}
}

type MapFn func(string) string

func Map(mapFn MapFn, parse ParseFn) ParseFn {
	return func(input string) (string, string, error) {
		char, remaining, err := parse(input)
		if err != nil {
			return "", "", err
		}

		return mapFn(char), remaining, nil
	}
}

package pc

import (
	"errors"
	"fmt"
)

type ParseFn func(string) (string, string, error)

// Accepts any of the given chars in char
func CharParser(anyChar string) ParseFn {
	return func(input string) (string, string, error) {
		if input == "" {
			return "", "", errors.New("no more input")
		}

		for i := 0; i < len(anyChar); i++ {
			if anyChar[i] == input[0] {
				return string(input[0]), input[1:], nil
			}
		}

		return "", "", fmt.Errorf("Expected '%s'. Got '%c'", anyChar, input[0])
	}
}

func AndThen(first ParseFn, second ParseFn, name string) ParseFn {
	return func(input string) (string, string, error) {
		firstChar, firstRemaining, firstErr := first(input)
		if firstErr != nil {
			return "", "", fmt.Errorf("first group of '%s' failed: %s", name, firstErr)
		}

		secondChar, secondRemaining, secondErr := second(firstRemaining)
		if secondErr != nil {
			return "", "", fmt.Errorf("second group of '%s' failed: %s", name, secondErr)

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

func OneOrMore(parser ParseFn, name string) ParseFn {
	return func(input string) (string, string, error) {
		result := ""
		remaining := input
		for {
			presult, premaining, perr := parser(remaining)
			if perr != nil {
				if remaining == input {
					return "", "", fmt.Errorf("None matched for '%s', '%s'", name, perr)
				}
				return result, remaining, nil
			}
			result += presult
			remaining = premaining
		}
	}
}

func ZeroOrOne(parser ParseFn) ParseFn {
	return func(input string) (string, string, error) {
		presult, premaining, perr := parser(input)
		if perr == nil {
			return presult, premaining, nil
		}
		return "", input, nil
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

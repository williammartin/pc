package pc

func ParseA(input string) (bool, string) {
	if input == "" {
		return false, ""
	}

	if input[0] == 'A' {
		return true, input[1:]
	}

	return false, input
}

package main

import "errors"

// pad characters to start of a string
func padStart(originalString string, desiredLen int, paddingChar rune) (string, error) {
	if len(originalString) > desiredLen {
		return "", errors.New(padLengthErrorString)
	}

	lenDiff := desiredLen - len(originalString)

	if lenDiff == 0 {
		return originalString, nil
	}

	preFix := ""

	for range lenDiff {
		preFix += string(paddingChar)
	}

	return preFix + originalString, nil
}

// pad characters to the end of a string
func padEnd(originalString string, desiredLen int, paddingChar rune) (string, error) {
	if len(originalString) > desiredLen {
		return "", errors.New(padLengthErrorString)
	}

	lenDiff := desiredLen - len(originalString)

	if lenDiff == 0 {
		return originalString, nil
	}

	postFix := ""

	for range lenDiff {
		postFix += string(paddingChar)
	}

	return originalString + postFix, nil
}

// Pad both sides. If odd characters are to be padded, the longer string is padded to the start of the string.
// Usually the longer string is actually padded to the end, but this serves the purpose of this utility class
func padCenter(originalString string, desiredLen int, paddingChar rune) (string, error) {
	if len(originalString) > desiredLen {
		return "", errors.New(padLengthErrorString)
	}

	lenDiff := desiredLen - len(originalString)

	toPadEnd := lenDiff / 2
	toPadStart := lenDiff - toPadEnd

	resStr := originalString
	resStr, err := padEnd(originalString, len(resStr)+toPadEnd, paddingChar)
	if err != nil {
		return "", err
	}

	resStr, err = padStart(resStr, len(resStr)+toPadStart, paddingChar)
	if err != nil {
		return "", err
	}

	return resStr, nil
}

package csv2md

import (
	"errors"
	"unicode/utf8"
)

const padLengthErrorString = "the length of the original string already exceeded desired length"

// pad characters to start of a string
func padStart(originalString string, desiredLen int, paddingChar rune) (string, error) {
	if utf8.RuneCountInString(originalString) > desiredLen {
		return "", errors.New(padLengthErrorString)
	}

	lenDiff := desiredLen - utf8.RuneCountInString(originalString)

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
	if utf8.RuneCountInString(originalString) > desiredLen {
		return "", errors.New(padLengthErrorString)
	}

	lenDiff := desiredLen - utf8.RuneCountInString(originalString)

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
func padCenter(originalString string, desiredLen int, paddingChar rune) (string, error) {
	if utf8.RuneCountInString(originalString) > desiredLen {
		return "", errors.New(padLengthErrorString)
	}

	lenDiff := desiredLen - utf8.RuneCountInString(originalString)

	toPadStart := lenDiff / 2
	toPadEnd := lenDiff - toPadStart

	resStr := originalString
	resStr, err := padEnd(originalString, utf8.RuneCountInString(resStr)+toPadEnd, paddingChar)
	if err != nil {
		return "", err
	}

	resStr, err = padStart(resStr, utf8.RuneCountInString(resStr)+toPadStart, paddingChar)
	if err != nil {
		return "", err
	}

	return resStr, nil
}

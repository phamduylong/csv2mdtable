package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	res, err := convert(
		`UID,First name,Last name,Email,Phone
0,Jane,Smith,jane.smith@email.com,555-555-1212
1,John,Doe,john.doe@email.com,555-555-3434
2,Alice,Wonder,alice@wonderland.com,555-555-5656
3,Aaron,Potter`)
	if err != nil {
		fmt.Printf("An error occurred 🙄: %s", err.Error())
	} else {
		fmt.Printf("Converted successfully ✅ Result:\n%s", res)
	}
	fmt.Printf("Press Enter to continue...")
	fmt.Scanln()
}

// convert csv string into a markdown table
func convert(csv string) (string, error) {
	if csv == "" {
		return "", errors.New("empty CSV input")
	}
	result := ""
	lines := strings.Split(csv, "\n")
	colCount := getColumnCount(lines)

	// max length of each column so we can beautify the table
	maxLenOfCol := getMaxColumnLengths(lines, colCount)

	// constructing each data line
	for idx := range len(lines) {
		originalLine := strings.ReplaceAll(lines[idx], "\n", "")
		// array containing field values in the current line
		colVals := strings.Split(originalLine, ",")
		// fill empty column values with empty strings
		for range colCount - len(colVals) {
			colVals = append(colVals, "")
		}
		convertedLine := ""

		for i := range colCount {
			paddedString, err := padCenter(colVals[i], maxLenOfCol[i], ' ')
			if err != nil {
				return "", errors.New("something happened when padding value " + colVals[i] + " row: " + fmt.Sprint(idx) + " col: " + fmt.Sprint(i) + ". Error message: " + err.Error())
			}
			convertedLine += "| " + paddedString + " "
		}

		// add final column closer and new line
		convertedLine += "|\n"
		// append to result string
		result += convertedLine

		if idx == 0 {
			separatorLine := constructSeparatorLine(colCount, maxLenOfCol)
			result += separatorLine
		}
	}

	return result, nil
}

// Construct a separator line between the header line and data lines
func constructSeparatorLine(colsCount int, maxLens []int) string {
	separatorLine := "| "
	for i := range colsCount {
		dashes := ""
		// loop through max length of each column and add dashes
		for range maxLens[i] {
			dashes += "-"
		}
		separatorLine += dashes + " | "
	}
	separatorLine += "\n"
	return separatorLine
}

// Get the amount of columns in the csv file
func getColumnCount(lines []string) int {
	maxCol := 0
	nrOfCols := 0
	for _, line := range lines {
		if line == "" {
			continue
		}
		nrOfCols = strings.Count(line, ",") + 1
		if nrOfCols > maxCol {
			maxCol = nrOfCols
		}
	}

	return maxCol
}

// Get max length of each columns
func getMaxColumnLengths(lines []string, colCount int) []int {
	maxLens := make([]int, colCount)
	for _, line := range lines {
		colVals := strings.Split(line, ",")
		for colIdx, colVal := range colVals {
			if len(colVal) > maxLens[colIdx] {
				maxLens[colIdx] = len(colVal)
			}
		}
	}

	return maxLens
}

// pad characters to start of a string
func padStart(originalString string, desiredLen int, paddingChar rune) (string, error) {
	if len(originalString) > desiredLen {
		return "", errors.New("the length of the original string already exceeded desired length")
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
		return "", errors.New("the length of the original string already exceeded desired length")
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
		return "", errors.New("the length of the original string already exceeded desired length")
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

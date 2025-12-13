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
	const emptyCol = "| "
	lines := strings.Split(csv, "\n")
	colCount := getColumnCount(lines)

	// max length of each column so we can beautify the table
	maxLenOfCol := getMaxColumnLengths(lines, colCount)
	for idx := range len(lines) {
		originalLine := strings.ReplaceAll(lines[idx], "\n", "")
		colVals := strings.Split(originalLine, ",")
		newLine := ""

		for i := range colCount {
			if i < len(colVals) {
				newLine += fmt.Sprintf("| %s ", colVals[i])
			} else {
				newLine += emptyCol
			}
		}

		// add final column closer and new line
		newLine += "|\n"
		// replace awkward double spaces
		newLine = strings.ReplaceAll(newLine, "  ", " ")
		// append to result string
		result += newLine

		if idx == 0 {
			separatorLine := constructSeparatorLine(colCount)
			result += separatorLine
		}
	}

	return result, nil
}

// Construct a separator line between the header line and data lines
func constructSeparatorLine(colsCount int) string {
	separatorLine := "| "
	for range colsCount {
		separatorLine += "- | "
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

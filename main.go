package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.design/x/clipboard"
)

const padLengthErrorString = "the length of the original string already exceeded desired length"

func main() {

	csvString := ""
	var err error = nil
	cfg, cfgErr := parseConfig()

	if cfgErr != nil {
		log.Fatalf("Configuration error: %s", cfgErr.Error())
	}

	// if http:
	if cfg.URL != "" {
		if cfg.VerboseLogging {
			log.Println("Reading CSV from URL: " + cfg.URL)
		}
		csvString, err = getCSVStringFromUrl(cfg.URL)
	}

	// if csv from file:
	if cfg.FilePath != "" {
		if cfg.VerboseLogging {
			log.Println("Reading from file path: " + cfg.FilePath)
		}
		csvString, err = getCSVStringFromPath(cfg.FilePath)
	}

	if err != nil {
		log.Fatalln(err)
	}

	res, err := convert(csvString, cfg)

	if err != nil {
		fmt.Printf("An error occurred 🙄: %s\n", err.Error())
	}

	if cfg.AutoCopy {
		clipboardInitErr := clipboard.Init()
		if clipboardInitErr != nil {
			log.Fatalln("failed to initiate clipboard")
		}
		clipboard.Write(clipboard.FmtText, []byte(res))
		log.Println("Copied to clipboard 📋")
	}

	fmt.Printf("Press Enter to continue...")
	fmt.Scanln()
}

// convert csv string into a markdown table
func convert(csv string, cfg Config) (string, error) {
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
			paddedString := ""
			var err error = nil

			switch cfg.Align {
			case Left:
				paddedString, err = padEnd(colVals[i], maxLenOfCol[i], ' ')
			case Right:
				paddedString, err = padStart(colVals[i], maxLenOfCol[i], ' ')
			case Center:
				paddedString, err = padCenter(colVals[i], maxLenOfCol[i], ' ')
			}

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

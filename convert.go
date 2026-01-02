package csv2mdtable

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"unicode/utf8"
)

// Convert CSV string into a markdown table. Returns the string representation of the markdown table if converted successfully and an error if failed.
func Convert(csv string, cfg Config) (string, error) {

	if csv == "" {
		return "", fmt.Errorf("csv string is empty")
	}

	cfgErr := ValidateConfig(cfg)

	if cfgErr != nil {
		return "", fmt.Errorf("Configuration error: %s\n", cfgErr)
	}

	if cfg.VerboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	// escape pipe characters
	csv = strings.ReplaceAll(csv, "|", `\|`)

	csvReader := createCSVReader(cfg, csv)

	records, readErr := csvReader.ReadAll()

	if readErr != nil {
		return "", fmt.Errorf("Failed to parse CSV. Error: %s", readErr)
	}

	cfg.excludedColumnsIndices = getIndicesOfExcludedColumns(cfg.ExcludedColumns, records[0])

	if len(cfg.excludedColumnsIndices) > 0 && len(cfg.excludedColumnsIndices) == len(records[0]) {
		slog.Warn("All columns were excluded from conversion. Returning an empty string")
		return "", nil
	}

	cfg = populateColumnIndices(cfg, records[0])

	colCount := csvReader.FieldsPerRecord
	result := ""

	if cfg.Caption != "" {
		result += fmt.Sprintf("<!-- %s -->\n", cfg.Caption)
	}

	// max length of each column so we can beautify the table
	maxLenOfCol := getMaxColumnLengths(records, cfg.Align)

	// constructing each data line
	for idx := range len(records) {
		convertedLine, err := constructDataLine(records[idx], cfg, maxLenOfCol, idx)

		if err != nil {
			return "", err
		}

		convertedLine = strings.TrimSpace(convertedLine)

		// only attach a new line if it's not the last line in the table
		if idx < len(records)-1 {
			convertedLine += "\n"
		}

		// append to result string
		result += convertedLine

		// after first line, we shall get a separator line
		if idx == 0 {
			separatorLine := constructSeparatorLine(colCount, maxLenOfCol, cfg)
			result += separatorLine
		}
	}

	return result, nil
}

// Construct data line
func constructDataLine(colVals []string, cfg Config, maxLenOfCol []int, currRowIdx int) (string, error) {
	if cfg.Compact {
		return constructCompactDataLine(colVals, cfg)
	} else {
		return constructBeautifulDataLine(colVals, cfg, maxLenOfCol, currRowIdx)
	}
}

// Construct a well-formatted data line
func constructBeautifulDataLine(colVals []string, cfg Config, maxLenOfCol []int, currRowIdx int) (string, error) {

	convertedLine := "| "

	for _, i := range cfg.orderedColumnsIndices {
		// If current column is excluded, ignore it
		if slices.Contains(cfg.excludedColumnsIndices, i) {
			continue
		}

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
			return "", errors.New("something happened when padding value " + colVals[i] + " row: " + fmt.Sprint(currRowIdx) +
				" col: " + fmt.Sprint(i) + ". Error message: " + err.Error())
		}

		convertedLine += paddedString + " | "
	}

	return convertedLine, nil
}

// Construct a compact data line
func constructCompactDataLine(colVals []string, cfg Config) (string, error) {

	convertedLine := "|"

	for _, i := range cfg.orderedColumnsIndices {
		// If current column is excluded, ignore it
		if slices.Contains(cfg.excludedColumnsIndices, i) {
			continue
		}
		convertedLine += colVals[i] + "|"
	}

	return convertedLine, nil
}

// Construct a separator line between the header line and data lines
func constructSeparatorLine(colsCount int, maxLenOfCol []int, cfg Config) string {
	if cfg.Compact {
		// since we're in compact mode, all columns separator will look alike. We just care about the number of columns included
		return constructCompactSeparatorLine(colsCount-len(cfg.excludedColumnsIndices), cfg.Align)
	} else {
		return constructBeautifulSeparatorLine(cfg, maxLenOfCol)
	}
}

// Construct a well-formatted separator line
func constructBeautifulSeparatorLine(cfg Config, maxLenOfCol []int) string {

	separatorLine := "| "

	for _, i := range cfg.orderedColumnsIndices {
		// If current column is excluded, ignore it
		if slices.Contains(cfg.excludedColumnsIndices, i) {
			continue
		}

		dashes := ""
		// loop through max length of each column and add dashes
		for range maxLenOfCol[i] {
			dashes += "-"
		}
		switch cfg.Align {
		case Left:
			// replace the first dash with a colon. This makes the rendered table align text on the left hand side
			dashes = strings.Replace(dashes, "-", ":", 1)
		case Right:
			// replace the last dash with a colon. This makes the rendered table align text on the right hand side
			i := strings.LastIndex(dashes, "-")
			excludingLast := dashes[:i] + strings.Replace(dashes[i:], "-", "", 1)
			dashes = excludingLast + ":"
		case Center:
			// replace the first and last dashes with colons
			// first
			dashes = strings.Replace(dashes, "-", ":", 1)

			// last
			i := strings.LastIndex(dashes, "-")
			excludingLast := dashes[:i] + strings.Replace(dashes[i:], "-", "", 1)
			dashes = excludingLast + ":"
		}
		separatorLine += dashes + " | "
	}

	// trim any potential leading/following whitespaces and add new line character
	separatorLine = strings.TrimSpace(separatorLine)
	separatorLine += "\n"

	return separatorLine
}

// Construct a compact separator line
func constructCompactSeparatorLine(colCount int, align Align) string {
	separatorLine := "|"
	for range colCount {
		switch align {
		case Left:
			separatorLine += ":-|"
		case Right:
			separatorLine += "-:|"
		case Center:
			separatorLine += ":-:|"
		}
	}

	// trim any potential leading/following whitespaces and add new line character
	separatorLine = strings.TrimSpace(separatorLine)
	separatorLine += "\n"

	return separatorLine
}

// Get max length of each columns
func getMaxColumnLengths(lines [][]string, align Align) []int {
	maxLens := make([]int, len(lines[0]))
	for _, fields := range lines {
		for fieldIdx, fieldVal := range fields {
			if utf8.RuneCountInString(fieldVal) > maxLens[fieldIdx] {
				maxLens[fieldIdx] = utf8.RuneCountInString(fieldVal)
			}
		}
	}

	for idx, colLen := range maxLens {
		if colLen <= 2 && align == Center {
			// if align is center, we need at least 3 spaces (:-:)
			maxLens[idx] = 3
		} else if colLen < 2 && align != Center {
			maxLens[idx] = 2
		}
	}

	return maxLens
}

// Get the indices of columns that are excluded in config
func getIndicesOfExcludedColumns(excludedColumns []string, headerLine []string) []int {
	var excludedColumnsIndices []int
	if len(excludedColumns) > 0 {
		for colIdx := range len(headerLine) {
			// if column is found in the csv and not duplicated
			if slices.Contains(excludedColumns, headerLine[colIdx]) {
				excludedColumnsIndices = append(excludedColumnsIndices, colIdx)
			}
		}
	}
	return excludedColumnsIndices
}
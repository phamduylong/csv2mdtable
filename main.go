package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

func main() {
	startTime := time.Now()

	csvString := ""
	var err error = nil

	cfg, cfgErr := parseConfig()

	if cfgErr != nil {
		slog.Error(fmt.Sprintf("Configuration error: %s\n", cfgErr.Error()))
		return
	}

	if cfg.VerboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	csvString, err = getCSVStringFromSource(cfg)

	if err != nil {
		slog.Error(fmt.Sprintf("An error occurred when fetching data: %s\n", err.Error()))
		return
	}

	res, err := convert(csvString, cfg)

	if err != nil {
		slog.Error(fmt.Sprintf("An error occurred 🙄: %s\n", err.Error()))
		return
	}

	if cfg.OutputFilePath != "" {
		err = writeMarkdownTableToFile(cfg.OutputFilePath, res)
		if err != nil {
			slog.Error(fmt.Sprintf("Failure when writing to output file: %s\n", err.Error()))
		}
	}

	// auto-copy activated, copy the converted table to clipboard
	// only warn if we failed to initiate clipboard module. Program should finish even if we had an error in this phase
	if cfg.AutoCopy {
		clipboardInitErr := clipboard.Init()
		if clipboardInitErr != nil {
			slog.Warn("Failed to initiate clipboard\n", "Initial error", clipboardInitErr)
		}
		clipboard.Write(clipboard.FmtText, []byte(res))
		slog.Info("Copied to clipboard 📋\n")
	}

	// output to window set to true, dump the table out to console
	if cfg.OutputToWindow {
		slog.Info("Converted table:\n" + res + "\n")
	}

	slog.Info(fmt.Sprintf("Elapsed time: %s", durationToReadableString(time.Since(startTime))))

	slog.Info("Press Enter to continue...\n")
	fmt.Scanln()
}

// convert csv string into a markdown table
func convert(csvString string, cfg Config) (string, error) {
	if csvString == "" {
		return "", errors.New("empty CSV input")
	}

	r := csv.NewReader(strings.NewReader(csvString))
	r.LazyQuotes = true
	records, err := r.ReadAll()

	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse CSV. Error: %s", err))
	}

	colCount := len(records[0])
	result := ""

	// max length of each column so we can beautify the table
	maxLenOfCol := getMaxColumnLengths(records)

	// constructing each data line
	for idx := range len(records) {
		// array containing field values in the current line
		colVals := records[idx]
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

// Get max length of each columns
func getMaxColumnLengths(lines [][]string) []int {
	maxLens := make([]int, len(lines[0]))
	for _, fields := range lines {
		for fieldIdx, fieldVal := range fields {
			if utf8.RuneCountInString(fieldVal) > maxLens[fieldIdx] {
				maxLens[fieldIdx] = utf8.RuneCountInString(fieldVal)
			}
		}
	}

	return maxLens
}

func getCSVStringFromSource(cfg Config) (csvString string, err error) {
	sourceOfData := ""

	// if csv from url
	if cfg.URL != "" {
		if cfg.VerboseLogging {
			slog.Info(fmt.Sprintf("Reading CSV from URL: %s\n", cfg.URL))
		}
		sourceOfData = cfg.URL
		csvString, err = getCSVStringFromUrl(cfg.URL)
	}

	// if csv from file
	if cfg.InputFilePath != "" {
		if cfg.VerboseLogging {
			slog.Debug(fmt.Sprintf("Reading CSV from file: %s\n", cfg.InputFilePath))
		}
		sourceOfData = cfg.InputFilePath
		csvString, err = getCSVStringFromFile(cfg.InputFilePath)
	}

	if err != nil {
		return "", fmt.Errorf("failed to read CSV from source.\nSource: %s\nOriginal error: %s", sourceOfData, err.Error())
	}

	return csvString, nil

}

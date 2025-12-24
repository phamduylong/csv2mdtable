package csv2md

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func createCSVReader(cfg Config, csvString string) *csv.Reader {
	r := csv.NewReader(strings.NewReader(csvString))

	if cfg.Delimiter != 0 {
		r.Comma = cfg.Delimiter
	}

	r.LazyQuotes = true

	return r
}

// Get CSV String from the source specified in the config object
func getCSVStringFromSource(cfg Config) (csvString string, err error) {
	sourceOfData := ""

	// if csv from url
	if cfg.URL != "" {
		sourceOfData = cfg.URL
		csvString, err = getCSVStringFromUrl(cfg)
	}

	// if csv from file
	if cfg.InputFilePath != "" {
		sourceOfData = cfg.InputFilePath
		csvString, err = getCSVStringFromFile(cfg)
	}

	if err != nil {
		return "", fmt.Errorf("failed to read CSV from source.\nSource: %s\nOriginal error: %s", sourceOfData, err.Error())
	}

	return csvString, nil

}

// Get CSV data from a URL. Will make a GET request with accepted content text/csv type
func getCSVStringFromUrl(cfg Config) (csvString string, err error) {
	req, err := http.NewRequest("GET", cfg.URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "text/csv; charset=utf-8")

	client := &http.Client{}

	if cfg.VerboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug(fmt.Sprintf("Reading CSV from URL: %s\n", cfg.URL))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		// of course, we are just relying on assumptions that the server is sending legit status
		return "", errors.New(string(b))
	}

	if len(b) == 0 {
		return "", errors.New("request returned an empty body")
	}

	csvFromUrl := string(b)
	return csvFromUrl, nil
}

// Read csv data from file. File must have .csv extension or else an error will be returned
func getCSVStringFromFile(cfg Config) (csvString string, err error) {
	if _, err := os.Stat(cfg.InputFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("path %s does not exist", cfg.InputFilePath)
	}

	fileExtension := filepath.Ext(cfg.InputFilePath)

	if fileExtension != ".csv" {
		return "", errors.New("given file is not a csv file")
	}

	absolutePath, _ := filepath.Abs(cfg.InputFilePath)

	if cfg.VerboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug(fmt.Sprintf("Reading CSV from file: %s\n", absolutePath))
	}

	fileContent, err := os.ReadFile(cfg.InputFilePath)

	if err != nil {
		return "", err
	}

	return string(fileContent[:]), nil
}
package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

// get CSV data from a URL. Will make a GET request with accepted content text/csv type
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

// read csv data from file. File must have .csv extension or else an error will be returned
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

// write csv data to file
func writeMarkdownTableToFile(path string, csvString string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.WriteString(csvString)

	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// get CSV data from a URL. Will make a GET request with accepted content text/csv type
func getCSVStringFromUrl(url string) (csvString string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "text/csv; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	csvFromUrl := string(b)
	return csvFromUrl, nil
}

// read csv data from file. File must have .csv extension or else an error will be returned
func getCSVStringFromFile(path string) (csvString string, err error) {
	fileExtension := filepath.Ext(path)

	if fileExtension != ".csv" {
		return "", errors.New("given file is not a csv file")
	}

	fileContent, err := os.ReadFile(path)
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

package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

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

func getCSVStringFromPath(path string) (csvString string, err error) {
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

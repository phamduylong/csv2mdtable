package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jaswdr/faker"
)

const TEN_STRING_FIELDS_FORMAT = "\"%s\"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\""

func generateCSVRows(rowCount int) []string {
	var res []string
	fake := faker.New()
	newRow := ""
	for range rowCount {
		person := fake.Person()
		newRow = personObjectToCSVRow(person)
		res = append(res, newRow)
	}

	return res
}

func personObjectToCSVRow(person faker.Person) string {
	fake := faker.New()
	addr := fake.Address()
	return fmt.Sprintf(TEN_STRING_FIELDS_FORMAT,
		person.FirstName(), person.LastName(), person.Gender(), person.SSN(), person.Title(), person.Contact().Email, person.Contact().Phone, addr.Address(), addr.PostCode(), addr.City())
}

func createGenericConfig() Config {
	var cfg Config
	return cfg
}

// call this function to do a test run
func testRun() {
	csvRows := generateCSVRows(100000)
	rows := strings.Join(csvRows, "\n")
	cfg := createGenericConfig()
	startTime := time.Now()
	res, _ := Convert(rows, cfg)
	elapsed := durationToReadableString(time.Since(startTime))
	fmt.Println(res)
	fmt.Println(elapsed)
}

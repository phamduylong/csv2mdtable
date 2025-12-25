package csv2mdtable

import (
	"encoding/csv"
	"strings"
)

func createCSVReader(cfg Config, csvString string) *csv.Reader {
	r := csv.NewReader(strings.NewReader(csvString))

	if cfg.CSVReaderConfig.Comma != 0 {
		r.Comma = cfg.CSVReaderConfig.Comma
	}

	if cfg.CSVReaderConfig.Comment != 0 {
		r.Comment = cfg.CSVReaderConfig.Comment
	}

	if r.FieldsPerRecord > 0 {
		r.FieldsPerRecord = cfg.CSVReaderConfig.FieldsPerRecord
	}

	r.LazyQuotes = cfg.CSVReaderConfig.LazyQuotes
	r.ReuseRecord = cfg.CSVReaderConfig.ReuseRecord
	r.TrimLeadingSpace = cfg.CSVReaderConfig.TrimLeadingSpace

	return r
}

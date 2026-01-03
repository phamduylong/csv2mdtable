package csv2mdtable

import (
	"errors"
	"log/slog"
	"slices"
	"strings"
)

type Align int

const (
	Center Align = 0
	Left   Align = 1
	Right  Align = 2
)

type ColumnSortOption int

type ColumnSortFunction func(a string, b string) int

const (
	None       ColumnSortOption = 0
	Ascending  ColumnSortOption = 1
	Descending ColumnSortOption = 2
	Custom     ColumnSortOption = 3
)

type Config struct {
	// Align the rendered content for the Markdown table. 0 = Center, 1 = Left, 2 = Right
	Align Align

	// Caption of the table (as an HTML comment)
	Caption string

	// Should the markdown table be the compact version
	Compact bool

	// Reader configuration to be used for Reader type.
	// See also https://pkg.go.dev/encoding/csv#Reader
	CSVReaderConfig CSVReaderConfig

	// List of columns to be excluded from table construction
	ExcludedColumns []string

	// Indices of excluded columns (internal)
	excludedColumnsIndices []int

	// Indices of columns to convert to
	orderedColumnsIndices []int

	// Should the columns be sorted and how?
	SortColumns ColumnSortOption

	// Custom sort function
	SortFunction ColumnSortFunction

	// Log detailed diagnostic messages when running the program.
	VerboseLogging bool
}

// Mirroring https://pkg.go.dev/encoding/csv#Reader
type CSVReaderConfig struct {
	// Comma is the field delimiter.
	// It is set to comma (',') by NewReader.
	// Comma must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	Comma rune

	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	// Comment must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	// It must also not be equal to Comma.
	Comment rune

	// FieldsPerRecord is the number of expected fields per record.
	// If FieldsPerRecord is positive, Read requires each record to
	// have the given number of fields. If FieldsPerRecord is 0, Read sets it to
	// the number of fields in the first record, so that future records must
	// have the same field count. If FieldsPerRecord is negative, no check is
	// made and records may have a variable number of fields.
	FieldsPerRecord int

	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool

	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool

	// ReuseRecord controls whether calls to Read may return a slice sharing
	// the backing array of the previous call's returned slice for performance.
	// By default, each call to Read returns newly allocated memory owned by the caller.
	ReuseRecord bool
}

// Validate the Config object passed as parameter.
// An error will be returned in case the configuration was invalid.
func ValidateConfig(cfg Config) error {
	configMalformed := false

	if cfg.VerboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Validating config 🤔")
	}

	if cfg.Align < Center || cfg.Align > Right {
		return errors.New("align value is out of range, please choose in range [0-2]")
	}

	if cfg.SortColumns < None || cfg.SortColumns > Custom {
		return errors.New("sort columns value is out of range, please choose in range [0-3]")
	}

	if cfg.SortColumns == Custom && cfg.SortFunction == nil {
		return errors.New("sort type is set to Custom but SortFunc was not set.")
	}

	// function passed in but not sort type is not custom
	if cfg.SortColumns != Custom && cfg.SortFunction != nil {
		configMalformed = true
		slog.Warn("Sort function only works when SortColumns is set to Custom. SortColumns received is " + sortColumnsToString(cfg.SortColumns) + ", ignoring SortFunc.")
	}

	if cfg.VerboseLogging && !configMalformed {
		slog.Debug("Config is valid ✅")
	}

	return nil
}

func sortColumnsToString(val ColumnSortOption) string {
	switch val {
	case None:
		return "None"
	case Ascending:
		return "Ascending"
	case Descending:
		return "Descending"
	case Custom:
		return "Custom"
	}

	return ""
}

// Populate orderColumnIndices in Config object
func populateColumnIndices(cfg Config, headerLine []string) Config {
	// get the new order of columns after sorted, compared to the original order of them.
	if cfg.SortColumns == None {
		for i := range len(headerLine) {
			cfg.orderedColumnsIndices = append(cfg.orderedColumnsIndices, i)
		}
	} else {
		cfg.orderedColumnsIndices = getIndicesAfterSorting(cfg, headerLine)
	}

	return cfg
}

// Get the indices of columns after sorted
func getIndicesAfterSorting(cfg Config, headerLine []string) []int {
	sortedColumns := make([]string, len(headerLine))
	copy(sortedColumns, headerLine)
	var columnsIndicesAfterSorting []int

	switch cfg.SortColumns {
	case Ascending:
		slices.SortFunc(sortedColumns, func(a, b string) int {
			return strings.Compare(strings.ToLower(a), strings.ToLower(b))
		})
	case Descending:
		slices.SortFunc(sortedColumns, func(a, b string) int {
			return strings.Compare(strings.ToLower(b), strings.ToLower(a))
		})
	case Custom:
		slices.SortFunc(sortedColumns, cfg.SortFunction)
	}

	for i := range sortedColumns {
		columnsIndicesAfterSorting = append(columnsIndicesAfterSorting, slices.Index(headerLine, sortedColumns[i]))
	}

	return columnsIndicesAfterSorting
}

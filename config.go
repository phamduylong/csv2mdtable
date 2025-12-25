package csv2mdtable

import (
	"errors"
	"log/slog"
)

type Align int

const (
	Center Align = 0
	Left   Align = 1
	Right  Align = 2
)

type Config struct {
	// Align the rendered content for the Markdown table. 0 = Center, 1 = Left, 2 = Right
	Align Align

	// Reader configuration to be used for Reader type.
	// See also https://pkg.go.dev/encoding/csv#Reader
	CSVReaderConfig CSVReaderConfig

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
	if cfg.VerboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Validating config 🤔")
	}

	if cfg.Align < Center || cfg.Align > Right {
		return errors.New("align value is out of range, please choose in range [0-2]")
	}

	if cfg.VerboseLogging {
		slog.Debug("Config is valid ✅")
	}

	return nil
}

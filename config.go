package csv2md

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
	Align          Align
	URL            string
	InputFilePath  string
	VerboseLogging bool
	Delimiter      rune
}

func ValidateConfig(cfg Config) error {
	if cfg.VerboseLogging {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Validating config 🤔")
	}

	if cfg.Align < Center || cfg.Align > Right {
		return errors.New("align value is out of range, please choose in range [0-2]")
	}

	if cfg.URL != "" && cfg.InputFilePath != "" {
		return errors.New("both URL and file path are given, please provide only one of them exclusively")
	}
	if cfg.URL == "" && cfg.InputFilePath == "" {
		return errors.New("both URL and file path are missing, please provide either one of them exclusively")
	}

	if cfg.VerboseLogging {
		slog.Debug("Config is valid ✅")
	}
	
	return nil
}

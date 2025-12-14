package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
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
	FilePath       string
	AutoCopy       bool
	VerboseLogging bool
}

func parseConfig() (Config, error) {
	pathPtr := flag.String("path", "", "Path of the CSV file to convert")
	urlPtr := flag.String("url", "", "URL of the CSV file to convert")
	alignPtr := flag.Int("align", int(Center), "How should text be aligned in the table?\n0. Center\n1. Left\n2. Right")
	autoCopyPtr := flag.Bool("autoCopy", false, "Should the converted markdown table be copied to clipboard automatically?")
	verboseLoggingPtr := flag.Bool("verboseLogging", false, "Should detailed diagnostic messages be logged?")
	flag.Parse()

	var cfg Config

	cfg.Align = Align(*alignPtr)
	cfg.FilePath = *pathPtr
	cfg.URL = *urlPtr
	cfg.AutoCopy = *autoCopyPtr
	cfg.VerboseLogging = *verboseLoggingPtr

	cfgErr := validateConfig(cfg)

	if cfgErr != nil {
		return Config{}, cfgErr
	}

	if cfg.VerboseLogging {
		jsonCfg, serializationErr := json.MarshalIndent(cfg, "", "    ")
		if serializationErr == nil {
			log.Printf("Config:\n%s", jsonCfg)
		}
	}

	return cfg, nil
}

func validateConfig(cfg Config) error {
	if cfg.URL != "" && cfg.FilePath != "" {
		return errors.New("both URL and file path are given, please provide only one of them exclusively")
	}
	if cfg.URL == "" && cfg.FilePath == "" {
		return errors.New("both URL and file path are missing, please provide either one of them exclusively")
	}

	return nil
}

// Package config handles loading and parsing the configuration file for the application.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type Config struct {
	ShowIntro                 *bool     `json:"showIntro"`
	CommitTitleCharLimit      *int      `json:"commitTitleCharLimit"`
	CommitBodyCharLimit       *int      `json:"commitBodyCharLimit"`
	CommitBodyLineLength      *int      `json:"commitBodyLineLength"`
	MessageTemplate           *string   `json:"messageTemplate"`
	MessageWithTicketTemplate *string   `json:"messageWithTicketTemplate"`
	Prefixes                  Prefixes  `json:"prefixes"`
	Coauthors                 CoAuthors `json:"coauthors"`
	Boards                    Boards    `json:"boards"`
	Scopes                    Scopes    `json:"scopes"`
	ReadContributorsFromGit   *bool     `json:"readContributorsFromGit"`
}

// New returns a new Config
func New() *Config {
	return &Config{}
}

var ErrNoMatchingPath = errors.New("no matching config for current path")

func (c *Config) LoadFile(filePath string) error {
	log.Debug("loading config file", "path", filePath)

	if filePath == "" {
		return errors.New("no path provided")
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(f, &raw); err != nil {
		return fmt.Errorf("error parsing the json file: %w", err)
	}

	knownKeys := map[string]struct{}{
		"showIntro":                 {},
		"commitTitleCharLimit":      {},
		"commitBodyCharLimit":       {},
		"commitBodyLineLength":      {},
		"messageTemplate":           {},
		"messageWithTicketTemplate": {},
		"prefixes":                  {},
		"coauthors":                 {},
		"boards":                    {},
		"scopes":                    {},
		"readContributorsFromGit":   {},
	}

	isSingleConfig := true
	for key := range raw {
		if _, ok := knownKeys[key]; !ok {
			isSingleConfig = false
			break
		}
	}

	if isSingleConfig {
		if err := json.Unmarshal(f, &c); err != nil {
			return fmt.Errorf("error parsing the json file: %w", err)
		}
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	cwd = filepath.Clean(cwd)

	var pathConfigs map[string]Config
	if err := json.Unmarshal(f, &pathConfigs); err != nil {
		return fmt.Errorf("error parsing the json file: %w", err)
	}

	for path, cfg := range pathConfigs {
		if filepath.Clean(path) == cwd {
			*c = cfg
			return nil
		}
	}

	for path, cfg := range pathConfigs {
		match, err := filepath.Match(filepath.Clean(path), cwd)
		if err == nil && match {
			*c = cfg
			return nil
		}
	}

	// try matching on relative sub-paths (e.g. config key without trailing slash)
	for path, cfg := range pathConfigs {
		if strings.HasPrefix(cwd, filepath.Clean(path)+string(filepath.Separator)) {
			*c = cfg
			return nil
		}
	}

	return ErrNoMatchingPath
}

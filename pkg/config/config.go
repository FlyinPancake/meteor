// Package config handles loading and parsing the configuration file for the application.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	pathLike := false
	for key := range raw {
		if strings.ContainsAny(key, `/\`) || strings.ContainsAny(key, "*?[") || filepath.IsAbs(key) {
			pathLike = true
			break
		}
	}

	if !pathLike {
		if err := json.Unmarshal(f, c); err != nil {
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

	type match struct {
		pattern   string
		cleaned   string
		cfg       Config
		matchType int
	}

	const (
		matchExact = iota
		matchGlob
		matchPrefix
	)

	ancestors := []string{cwd}
	for dir := filepath.Dir(cwd); dir != "/" && dir != "." && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		ancestors = append(ancestors, dir)
		if dir == filepath.Dir(dir) {
			break
		}
	}

	matches := make([]match, 0)

	for path, cfg := range pathConfigs {
		cleaned := filepath.Clean(path)

		for _, candidate := range ancestors {
			if cleaned == candidate {
				matches = append(matches, match{pattern: path, cleaned: cleaned, cfg: cfg, matchType: matchExact})
				break
			}

			if ok, err := filepath.Match(cleaned, candidate); err == nil && ok {
				matches = append(matches, match{pattern: path, cleaned: cleaned, cfg: cfg, matchType: matchGlob})
				break
			}
		}

		if strings.HasPrefix(cwd, cleaned+string(filepath.Separator)) {
			matches = append(matches, match{pattern: path, cleaned: cleaned, cfg: cfg, matchType: matchPrefix})
		}
	}

	if len(matches) == 0 {
		log.Debug("no matching path in config; leaving config unchanged", "cwd", cwd, "path", filePath)
		return nil
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].matchType != matches[j].matchType {
			return matches[i].matchType < matches[j].matchType
		}
		if len(matches[i].cleaned) != len(matches[j].cleaned) {
			return len(matches[i].cleaned) > len(matches[j].cleaned)
		}
		return matches[i].pattern < matches[j].pattern
	})

	*c = matches[0].cfg
	return nil
}

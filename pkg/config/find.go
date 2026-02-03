package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
)

const (
	configFile   = ".meteor.json"
	reposFile    = "repos.json"
	globalConfig = "config.json"
	xdgConfigDir = ".config/meteor"
)

// FindConfigFile finds config files using the following order:
// 1. If the current directory contains configFile (.meteor.json), it will be used.
// 2. Traverse parent directories within the user's home directory for configFile and return the first found.
// 3. Check xdgConfigDir/reposFile (~/.config/meteor/repos.json); if present, return it.
// 4. Check xdgConfigDir/globalConfig (~/.config/meteor/config.json); if present, return it.
// 5. If none are found, return an error.
func FindConfigFile(fs afero.Fs, getWD func() (string, error), getHome func() (string, error)) (string, error) {
	if _, err := fs.Stat(configFile); err == nil {
		return filepath.Join("./", configFile), nil
	}

	homeDir, err := getHome()
	if err != nil {
		return "", fmt.Errorf("error getting home dir: %w", err)
	}

	currentDir, err := getWD()
	if err != nil {
		return "", fmt.Errorf("error getting current dir: %w", err)
	}

	for currentDir != "/" {
		rel, _ := filepath.Rel(homeDir, currentDir)
		if rel == ".." {
			break
		}

		filePath := filepath.Join(currentDir, configFile)
		log.Debug("checking for config file", "path", filePath)

		if _, err := fs.Stat(filePath); err == nil {
			return filePath, nil
		}

		currentDir = filepath.Join(currentDir, "..")
	}

	xdgReposFile := filepath.Join(homeDir, xdgConfigDir, reposFile)
	log.Debug("checking for repo config file", "path", xdgReposFile)
	if _, err := fs.Stat(xdgReposFile); err == nil {
		return xdgReposFile, nil
	}

	xdgConfigFile := filepath.Join(homeDir, xdgConfigDir, globalConfig)
	log.Debug("checking for config file", "path", xdgConfigFile)
	if _, err := fs.Stat(xdgConfigFile); err == nil {
		return xdgConfigFile, nil
	}

	return "", errors.New("no config file found")
}

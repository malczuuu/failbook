// Copyright (c) 2025 Damian Malczewski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// SPDX-License-Identifier: MIT

package problems

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

type Link struct {
	Title string `yaml:"title"`
	Href  string `yaml:"href"`
}

type ProblemConfig struct {
	Version     string `yaml:"version"`
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Title       string `yaml:"title"`
	StatusCode  int    `yaml:"status_code"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
	Links       []Link `yaml:"links"`
}

type ProblemRegistry struct {
	problems map[string]*ProblemConfig
}

func NewProblemRegistry() *ProblemRegistry {
	return &ProblemRegistry{
		problems: make(map[string]*ProblemConfig),
	}
}

func LoadFromDirectory(dirPath string) (*ProblemRegistry, error) {
	registry := NewProblemRegistry()

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("problems directory does not exist: %s", dirPath)
	}

	var loadFailures []error

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			loadFailures = append(loadFailures, fmt.Errorf("access error at %s: %w", path, err))
			return nil
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(d.Name())
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		if err := registry.loadFile(path); err != nil {
			loadFailures = append(loadFailures, fmt.Errorf("failed to load %s: %w", path, err))
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(loadFailures) > 0 {
		errorMsg := "failed to load error configurations:"
		for _, err := range loadFailures {
			errorMsg += fmt.Sprintf("\n  - %s", err.Error())
		}
		return nil, fmt.Errorf("%s", errorMsg)
	}

	log.Info().Int("count", len(registry.problems)).Msg("loaded problem configurations")
	return registry, nil
}

func validateProblemConfig(config *ProblemConfig) error {
	if config.Version != "1" {
		return fmt.Errorf("problem configuration version must be \"1\", got: %s", config.Version)
	}
	if config.ID == "" {
		return fmt.Errorf("problem configuration missing required field: id")
	}
	if config.Title == "" {
		return fmt.Errorf("problem configuration missing required field: title")
	}
	if config.StatusCode == 0 {
		return fmt.Errorf("problem configuration missing required field: status_code")
	}
	if config.Name == "" {
		config.Name = config.Title
	}
	return nil
}

func (r *ProblemRegistry) loadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	docIndex := 0

	for {
		var problem ProblemConfig
		err := decoder.Decode(&problem)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to parse YAML document %d: %w", docIndex, err)
		}

		if err := validateProblemConfig(&problem); err != nil {
			return fmt.Errorf("document %d: %w", docIndex, err)
		}

		if _, exists := r.problems[problem.ID]; exists {
			return fmt.Errorf("document %d: duplicate problem ID found: %s", docIndex, problem.ID)
		}

		r.problems[problem.ID] = &problem
		log.Debug().Str("id", problem.ID).Str("file", filePath).Int("document", docIndex).Msg("loaded problem configuration")
		docIndex++
	}

	if docIndex == 0 {
		return fmt.Errorf("no valid YAML documents found in file")
	}

	return nil
}

func (r *ProblemRegistry) Get(id string) (*ProblemConfig, bool) {
	errConfig, exists := r.problems[id]
	return errConfig, exists
}

func (r *ProblemRegistry) GetAll() map[string]*ProblemConfig {
	return r.problems
}

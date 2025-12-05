package problems

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateProblemConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      ProblemConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: ProblemConfig{
				Version:     "1",
				ID:          "404",
				Title:       "Not Found",
				StatusCode:  404,
				Summary:     "Resource not found",
				Description: "The resource could not be found",
				Links:       []Link{},
			},
			expectError: false,
		},
		{
			name: "missing version",
			config: ProblemConfig{
				ID:         "404",
				Title:      "Not Found",
				StatusCode: 404,
			},
			expectError: true,
			errorMsg:    "problem configuration version must be \"1\", got: ",
		},
		{
			name: "wrong version",
			config: ProblemConfig{
				Version:    "2",
				ID:         "404",
				Title:      "Not Found",
				StatusCode: 404,
			},
			expectError: true,
			errorMsg:    "problem configuration version must be \"1\", got: 2",
		},
		{
			name: "missing id",
			config: ProblemConfig{
				Version:    "1",
				Title:      "Not Found",
				StatusCode: 404,
			},
			expectError: true,
			errorMsg:    "problem configuration missing required field: id",
		},
		{
			name: "missing title",
			config: ProblemConfig{
				Version:    "1",
				ID:         "404",
				StatusCode: 404,
			},
			expectError: true,
			errorMsg:    "problem configuration missing required field: title",
		},
		{
			name: "missing status_code",
			config: ProblemConfig{
				Version: "1",
				ID:      "404",
				Title:   "Not Found",
			},
			expectError: true,
			errorMsg:    "problem configuration missing required field: status_code",
		},
		{
			name: "optional fields can be empty",
			config: ProblemConfig{
				Version:     "1",
				ID:          "404",
				Title:       "Not Found",
				StatusCode:  404,
				Summary:     "",
				Description: "",
				Links:       nil,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProblemConfig(&tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q but got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestProblemRegistry_LoadFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expectError bool
		errorMsg    string
		expectedIDs []string
	}{
		{
			name: "single document",
			fileContent: `version: "1"
id: "404"
title: "Not Found"
status_code: 404
summary: "Not found"
description: "Resource not found"
links: []`,
			expectError: false,
			expectedIDs: []string{"404"},
		},
		{
			name: "multi-document",
			fileContent: `version: "1"
id: "404"
title: "Not Found"
status_code: 404
summary: "Not found"
description: "Resource not found"
links: []
---
version: "1"
id: "500"
title: "Internal Server Error"
status_code: 500
summary: "Server error"
description: "Internal error"
links: []`,
			expectError: false,
			expectedIDs: []string{"404", "500"},
		},
		{
			name: "duplicate in same file",
			fileContent: `version: "1"
id: "404"
title: "Not Found"
status_code: 404
summary: "Not found"
description: "Resource not found"
links: []
---
version: "1"
id: "404"
title: "Duplicate"
status_code: 404
summary: "Duplicate"
description: "Duplicate"
links: []`,
			expectError: true,
			errorMsg:    "document 1: duplicate problem ID found: 404",
		},
		{
			name: "invalid version",
			fileContent: `version: "2"
id: "404"
title: "Not Found"
status_code: 404`,
			expectError: true,
			errorMsg:    "document 0: problem configuration version must be \"1\", got: 2",
		},
		{
			name:        "empty file",
			fileContent: "",
			expectError: true,
			errorMsg:    "no valid YAML documents found in file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.yaml")
			if err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0644); err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			registry := NewProblemRegistry()
			err := registry.loadFile(tmpFile)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q but got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				for _, id := range tt.expectedIDs {
					if _, exists := registry.problems[id]; !exists {
						t.Errorf("expected ID %q to be loaded but it wasn't", id)
					}
				}
				if len(registry.problems) != len(tt.expectedIDs) {
					t.Errorf("expected %d problems but got %d", len(tt.expectedIDs), len(registry.problems))
				}
			}
		})
	}
}

func TestProblemRegistry_GetAndGetAll(t *testing.T) {
	registry := NewProblemRegistry()

	config1 := &ProblemConfig{
		Version:    "1",
		ID:         "404",
		Title:      "Not Found",
		StatusCode: 404,
	}
	config2 := &ProblemConfig{
		Version:    "1",
		ID:         "500",
		Title:      "Internal Server Error",
		StatusCode: 500,
	}

	registry.problems["404"] = config1
	registry.problems["500"] = config2

	t.Run("Get existing", func(t *testing.T) {
		config, exists := registry.Get("404")
		if !exists {
			t.Errorf("expected config to exist")
		}
		if config.ID != "404" {
			t.Errorf("expected ID 404 but got %s", config.ID)
		}
	})

	t.Run("Get non-existing", func(t *testing.T) {
		_, exists := registry.Get("999")
		if exists {
			t.Errorf("expected config not to exist")
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		all := registry.GetAll()
		if len(all) != 2 {
			t.Errorf("expected 2 configs but got %d", len(all))
		}
		if all["404"] == nil || all["500"] == nil {
			t.Errorf("expected both configs to be present")
		}
	})
}

func TestLoadFromDirectory(t *testing.T) {
	t.Run("load multiple files", func(t *testing.T) {
		tmpDir := t.TempDir()

		file1 := filepath.Join(tmpDir, "404.yaml")
		file2 := filepath.Join(tmpDir, "500.yaml")

		content1 := `version: "1"
id: "404"
title: "Not Found"
status_code: 404
summary: "Not found"
description: "Resource not found"
links: []`

		content2 := `version: "1"
id: "500"
title: "Internal Server Error"
status_code: 500
summary: "Server error"
description: "Internal error"
links: []`

		if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		registry, err := LoadFromDirectory(tmpDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(registry.problems) != 2 {
			t.Errorf("expected 2 problems but got %d", len(registry.problems))
		}
	})

	t.Run("duplicate across files", func(t *testing.T) {
		tmpDir := t.TempDir()

		content := `version: "1"
id: "404"
title: "Not Found"
status_code: 404
summary: "Not found"
description: "Resource not found"
links: []`

		file1 := filepath.Join(tmpDir, "file1.yaml")
		file2 := filepath.Join(tmpDir, "file2.yaml")

		if err := os.WriteFile(file1, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		if err := os.WriteFile(file2, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		_, err := LoadFromDirectory(tmpDir)
		if err == nil {
			t.Errorf("expected problem for duplicate IDs across files")
		}

		if err != nil && !containsString(err.Error(), "file2.yaml") {
			t.Errorf("expected problem to mention file2.yaml, got: %v", err)
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		_, err := LoadFromDirectory("/non/existent/path")
		if err == nil {
			t.Errorf("expected problem for non-existent directory")
		}
	})

	t.Run("ignore non-yaml files", func(t *testing.T) {
		tmpDir := t.TempDir()

		yamlContent := `version: "1"
id: "404"
title: "Not Found"
status_code: 404
summary: "Not found"
description: "Resource not found"
links: []`

		yamlFile := filepath.Join(tmpDir, "404.yaml")
		txtFile := filepath.Join(tmpDir, "readme.txt")

		if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("failed to create yaml file: %v", err)
		}
		if err := os.WriteFile(txtFile, []byte("ignore me"), 0644); err != nil {
			t.Fatalf("failed to create txt file: %v", err)
		}

		registry, err := LoadFromDirectory(tmpDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(registry.problems) != 1 {
			t.Errorf("expected 1 error but got %d", len(registry.problems))
		}
	})

	t.Run("multiple invalid files", func(t *testing.T) {
		tmpDir := t.TempDir()

		invalidContent1 := `version: "2"
id: "404"
title: "Not Found"
status_code: 404`

		invalidContent2 := `version: "1"
id: "500"
title: ""
status_code: 500`

		file1 := filepath.Join(tmpDir, "invalid1.yaml")
		file2 := filepath.Join(tmpDir, "invalid2.yaml")

		if err := os.WriteFile(file1, []byte(invalidContent1), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		if err := os.WriteFile(file2, []byte(invalidContent2), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		_, err := LoadFromDirectory(tmpDir)
		if err == nil {
			t.Errorf("expected error for multiple invalid files")
		}

		if err != nil {
			if !containsString(err.Error(), "invalid1.yaml") {
				t.Errorf("expected error to mention invalid1.yaml")
			}
			if !containsString(err.Error(), "invalid2.yaml") {
				t.Errorf("expected error to mention invalid2.yaml")
			}
		}
	})
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && contains(s, substr))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

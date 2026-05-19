package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/andre-0303/devpath-mcp/internal/models"
)

// Load reads progress from path. Returns empty progress if file does not exist.
func Load(path string) (*models.UserProgress, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &models.UserProgress{
				Completed: make(map[string]map[string]bool),
			}, nil
		}
		return nil, err
	}

	var progress models.UserProgress
	if err := json.Unmarshal(data, &progress); err != nil {
		return nil, err
	}

	if progress.Completed == nil {
		progress.Completed = make(map[string]map[string]bool)
	}
	return &progress, nil
}

// Save writes progress to path, creating the directory if needed.
func Save(path string, progress *models.UserProgress) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

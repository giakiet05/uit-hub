package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadFixture(scenario string) (*Fixture, error) {
	if scenario == "" {
		scenario = "default"
	}

	path := filepath.Join("data", "fixtures", scenario+".json")
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read fixture %q: %w", path, err)
	}

	var fixture Fixture
	if err := json.Unmarshal(bytes, &fixture); err != nil {
		return nil, fmt.Errorf("parse fixture %q: %w", path, err)
	}

	return &fixture, nil
}

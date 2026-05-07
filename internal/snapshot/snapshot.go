// Package snapshot provides functionality to capture and persist point-in-time
// scan results, enabling historical comparisons and regression detection.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftctl-report/internal/model"
)

// Entry represents a single stored snapshot with metadata.
type Entry struct {
	ID        string           `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Label     string           `json:"label,omitempty"`
	Result    model.ScanResult `json:"result"`
}

// Store manages snapshot persistence on disk.
type Store struct {
	dir string
}

// NewStore creates a Store backed by the given directory.
// The directory is created if it does not exist.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create directory: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save persists a ScanResult as a new snapshot entry.
// The entry ID is derived from the current UTC timestamp.
func (s *Store) Save(result model.ScanResult, label string) (Entry, error) {
	now := time.Now().UTC()
	entry := Entry{
		ID:        now.Format("20060102T150405Z"),
		CreatedAt: now,
		Label:     label,
		Result:    result,
	}
	path := filepath.Join(s.dir, entry.ID+".json")
	f, err := os.Create(path)
	if err != nil {
		return Entry{}, fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		return Entry{}, fmt.Errorf("snapshot: encode entry: %w", err)
	}
	return entry, nil
}

// Load reads a snapshot entry by its ID.
func (s *Store) Load(id string) (Entry, error) {
	path := filepath.Join(s.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, fmt.Errorf("snapshot: read file: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, fmt.Errorf("snapshot: decode entry: %w", err)
	}
	return entry, nil
}

// List returns all snapshot entries sorted chronologically (oldest first).
func (s *Store) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("snapshot: glob: %w", err)
	}
	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("snapshot: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("snapshot: decode %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

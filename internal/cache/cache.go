// Package cache provides a simple file-backed result cache for driftctl scan
// outputs, keyed by a hash of the input file path and modification time.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/example/driftctl-report/internal/model"
)

// Entry wraps a cached ScanResult with metadata.
type Entry struct {
	CachedAt time.Time        `json:"cached_at"`
	InputHash string          `json:"input_hash"`
	Result   model.ScanResult `json:"result"`
}

// Cache stores and retrieves scan results from a directory on disk.
type Cache struct {
	dir string
}

// New returns a Cache that persists entries under dir.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache: create directory: %w", err)
	}
	return &Cache{dir: dir}, nil
}

// key derives a stable filename from the input file path and its mtime.
func (c *Cache) key(inputPath string) (string, error) {
	info, err := os.Stat(inputPath)
	if err != nil {
		return "", fmt.Errorf("cache: stat input: %w", err)
	}
	raw := fmt.Sprintf("%s|%d", filepath.Clean(inputPath), info.ModTime().UnixNano())
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:]), nil
}

// Get returns a cached ScanResult for inputPath, or (nil, nil) on a miss.
func (c *Cache) Get(inputPath string) (*model.ScanResult, error) {
	k, err := c.key(inputPath)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(c.dir, k+".json"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read entry: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("cache: unmarshal entry: %w", err)
	}
	return &entry.Result, nil
}

// Put writes result to the cache, keyed by inputPath's current mtime.
func (c *Cache) Put(inputPath string, result model.ScanResult) error {
	k, err := c.key(inputPath)
	if err != nil {
		return err
	}
	entry := Entry{
		CachedAt:  time.Now().UTC(),
		InputHash: k,
		Result:    result,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal entry: %w", err)
	}
	dest := filepath.Join(c.dir, k+".json")
	return os.WriteFile(dest, data, 0o644)
}

// Invalidate removes the cache entry for inputPath if it exists.
func (c *Cache) Invalidate(inputPath string) error {
	k, err := c.key(inputPath)
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(c.dir, k+".json"))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

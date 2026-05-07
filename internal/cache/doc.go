// Package cache provides a lightweight, file-backed result cache for driftctl
// scan outputs.
//
// Cache entries are stored as JSON files in a configurable directory. Each
// entry is keyed by a SHA-256 hash derived from the input file path and its
// modification time, so the cache is automatically invalidated whenever the
// source file changes.
//
// Typical usage:
//
//	c, err := cache.New("/var/cache/driftctl-report")
//	if err != nil { ... }
//
//	if result, _ := c.Get(inputPath); result != nil {
//		// use cached result
//	}
//	// ... parse and store
//	_ = c.Put(inputPath, result)
package cache

// Package redact provides a Redactor that scrubs sensitive attribute values
// from drift scan results prior to rendering, exporting, or transmitting
// reports.
//
// Usage:
//
//	r := redact.New(redact.DefaultSensitiveKeys)
//	clean := r.Apply(scanResult)
//
// DefaultSensitiveKeys matches common substrings such as "password",
// "secret", "token", "api_key", and "private_key". A custom slice of
// substrings may be supplied to New for domain-specific scrubbing rules.
//
// Apply never mutates the original ScanResult; it returns a shallow copy
// with a new Attrs map for each affected Resource.
package redact

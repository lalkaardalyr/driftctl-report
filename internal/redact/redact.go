// Package redact provides utilities for scrubbing sensitive field values
// from drift scan results before they are written to reports or exported.
package redact

import (
	"strings"

	"github.com/snyk/driftctl-report/internal/model"
)

// DefaultSensitiveKeys contains field name substrings that are considered
// sensitive and should be redacted from resource attribute maps.
var DefaultSensitiveKeys = []string{
	"secret",
	"password",
	"token",
	"private_key",
	"access_key",
	"api_key",
	"credential",
}

const redactedValue = "[REDACTED]"

// Redactor scrubs sensitive values from scan results.
type Redactor struct {
	keys []string
}

// New returns a Redactor using the provided sensitive key substrings.
// Pass DefaultSensitiveKeys for standard behaviour.
func New(sensitiveKeys []string) *Redactor {
	return &Redactor{keys: sensitiveKeys}
}

// Apply returns a deep copy of the ScanResult with sensitive attribute
// values replaced by the redacted placeholder.
func (r *Redactor) Apply(result model.ScanResult) model.ScanResult {
	out := result
	out.DriftedResources = r.redactList(result.DriftedResources)
	out.UnmanagedResources = r.redactList(result.UnmanagedResources)
	return out
}

func (r *Redactor) redactList(resources []model.Resource) []model.Resource {
	if resources == nil {
		return nil
	}
	out := make([]model.Resource, len(resources))
	for i, res := range resources {
		out[i] = r.redactResource(res)
	}
	return out
}

func (r *Redactor) redactResource(res model.Resource) model.Resource {
	if len(res.Attrs) == 0 {
		return res
	}
	attrs := make(map[string]string, len(res.Attrs))
	for k, v := range res.Attrs {
		if r.isSensitive(k) {
			attrs[k] = redactedValue
		} else {
			attrs[k] = v
		}
	}
	res.Attrs = attrs
	return res
}

func (r *Redactor) isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range r.keys {
		if strings.Contains(lower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}

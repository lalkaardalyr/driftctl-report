package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier sends drift notifications to a generic HTTP webhook endpoint.
type WebhookNotifier struct {
	URL    string
	Client *http.Client
}

// webhookPayload is the JSON body sent to the webhook endpoint.
type webhookPayload struct {
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Severity string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// NewWebhookNotifier creates a WebhookNotifier with a default HTTP client.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL: url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send delivers a notification message to the configured webhook URL.
func (w *WebhookNotifier) Send(msg Message) error {
	severity := "info"
	if msg.Critical {
		severity = "critical"
	} else if msg.Warning {
		severity = "warning"
	}

	payload := webhookPayload{
		Subject:   msg.Subject,
		Body:      msg.Body,
		Severity:  severity,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.Client.Post(w.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("webhook: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}

	return nil
}

package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackNotifier sends notifications to a Slack webhook URL.
type SlackNotifier struct {
	WebhookURL string
	Client     *http.Client
}

// NewSlackNotifier creates a SlackNotifier with the given webhook URL.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
		Client:     &http.Client{},
	}
}

type slackPayload struct {
	Text string `json:"text"`
}

// Send posts the message to Slack as a simple text payload.
func (s *SlackNotifier) Send(msg Message) error {
	payload := slackPayload{
		Text: fmt.Sprintf("*%s*\n%s", msg.Subject, msg.Body),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.Client.Post(s.WebhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("slack: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}

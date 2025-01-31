package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Sends a message to a Slack webhook
func SendSlackMessage(webhookURL string, message SlackMessage) error {
	// Check if the webhook URL is empty
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	// Check that the webhook URL starts with the correct prefix
	if !strings.HasPrefix(webhookURL, "https://hooks.slack.com") {
		return fmt.Errorf("invalid webhook URL")
	}

	// Marshal the message to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Println("Sending message to Slack")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Slack: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Slack: %v", resp.Status)
	}

	return nil
}

// Creates a new Slack message body
func NewSlackMessageBody(event string, message string, url string) SlackMessage {
	msg := SlackMessage{
		Blocks: []SlackBlock{
			{Type: "header", Text: &SlackText{Type: "plain_text", Text: "WG Controller Alert"}},
			{Type: "section", Text: &SlackText{Type: "mrkdwn", Text: "*Event:* " + event}},
			{Type: "section", Text: &SlackText{Type: "mrkdwn", Text: "*Message:* " + message}},
		},
	}

	msg.Blocks = append(msg.Blocks, SlackBlock{
		Type: "actions",
		Elements: []SlackElement{
			{Type: "button", Text: SlackPlainText{Type: "plain_text", Text: "Open Dashboard", Emoji: true}, URL: url},
		},
	})

	return msg
}

// SlackMessage defines the message payload structure
type SlackMessage struct {
	Blocks []SlackBlock `json:"blocks"`
}

// SlackBlock represents a block in the message
type SlackBlock struct {
	Type     string         `json:"type"`
	Text     *SlackText     `json:"text,omitempty"`
	Elements []SlackElement `json:"elements,omitempty"`
}

// SlackText represents a text field in the message
type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SlackElement represents an interactive element, like a button
type SlackElement struct {
	Type string         `json:"type"`
	Text SlackPlainText `json:"text"`
	URL  string         `json:"url"`
}

// SlackPlainText represents plain text inside an interactive element
type SlackPlainText struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

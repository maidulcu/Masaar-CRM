package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewClient(baseURL, model string) *Client {
	return &Client{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	body, _ := json.Marshal(ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ollama: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama: request failed: %w", err)
	}
	defer resp.Body.Close()

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ollama: decode response: %w", err)
	}
	return result.Response, nil
}

// ScoreLead asks the LLM to score a lead from 0–100 based on context.
func (c *Client) ScoreLead(ctx context.Context, contactName, conversationSummary, source string) (string, error) {
	prompt := fmt.Sprintf(`You are a UAE real estate & B2B sales expert.

Score this lead from 0 to 100 and explain why in 2-3 sentences.
Return JSON: {"score": <int>, "reasoning": "<text>"}

Lead details:
- Name: %s
- Source: %s
- Conversation summary: %s

Respond only with the JSON object.`, contactName, source, conversationSummary)

	return c.Generate(ctx, prompt)
}

// DraftReply generates a context-aware WhatsApp reply suggestion.
func (c *Client) DraftReply(ctx context.Context, contactName, language, threadSummary string) (string, error) {
	lang := "English"
	if language == "ar" {
		lang = "Arabic"
	}

	prompt := fmt.Sprintf(`You are a professional UAE sales agent.

Write a warm, concise WhatsApp follow-up message in %s for this contact.
Keep it under 100 words. Be professional but friendly.

Contact: %s
Conversation context: %s

Reply only with the message text, nothing else.`, lang, contactName, threadSummary)

	return c.Generate(ctx, prompt)
}

// SummarizeThread creates a brief summary of a WhatsApp conversation.
func (c *Client) SummarizeThread(ctx context.Context, messages []string) (string, error) {
	if len(messages) == 0 {
		return "", nil
	}

	conversation := ""
	for i, m := range messages {
		conversation += fmt.Sprintf("%d. %s\n", i+1, m)
	}

	prompt := fmt.Sprintf(`Summarize this WhatsApp sales conversation in 2-3 sentences.
Focus on: customer interest, concerns raised, and next action needed.

Conversation:
%s

Summary:`, conversation)

	return c.Generate(ctx, prompt)
}

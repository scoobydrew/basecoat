package claude

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// MiniSuggestion is a single miniature returned by the Claude lookup.
type MiniSuggestion struct {
	Name     string `json:"name"`
	UnitType string `json:"unit_type"`
	Quantity int    `json:"quantity"`
}

// Client wraps the Anthropic API for miniature lookups.
type Client struct {
	client *anthropic.Client
}

func NewClient(apiKey string) *Client {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Client{client: &c}
}

// LookupMinis asks Claude for a list of miniatures in a given game/set.
func (c *Client) LookupMinis(ctx context.Context, game, set string) ([]MiniSuggestion, error) {
	prompt := fmt.Sprintf(
		`You are a tabletop miniature painting expert. List the miniatures included in the following game or set.

Game/System: %s
Set/Box: %s

Respond with ONLY a JSON array. Each element should have:
- "name": the miniature's name or type (string)
- "unit_type": the unit category (e.g. "infantry", "cavalry", "monster", "hero", "vehicle") (string)
- "quantity": typical quantity included in the set (integer, default 1 if unknown)

If you don't recognize the game or set, return an empty array [].
Do not include any text outside the JSON array.`,
		game, set,
	)

	msg, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5HaikuLatest,
		MaxTokens: 2048,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude api: %w", err)
	}

	if len(msg.Content) == 0 {
		return nil, fmt.Errorf("claude returned empty response")
	}

	raw := msg.Content[0].Text

	var minis []MiniSuggestion
	if err := json.Unmarshal([]byte(raw), &minis); err != nil {
		return nil, fmt.Errorf("parse claude response: %w (raw: %s)", err, raw)
	}

	return minis, nil
}

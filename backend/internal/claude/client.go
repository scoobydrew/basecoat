package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

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

// GameMeta holds optional metadata about a game to improve Claude's lookup accuracy.
type GameMeta struct {
	Publisher string
	Year      *int
}

// LookupMinis asks Claude for a list of miniatures in a given game/set.
func (c *Client) LookupMinis(ctx context.Context, game, set string, meta GameMeta) ([]MiniSuggestion, error) {
	var extras strings.Builder
	if meta.Publisher != "" {
		fmt.Fprintf(&extras, "\nPublisher: %s", meta.Publisher)
	}
	if meta.Year != nil {
		fmt.Fprintf(&extras, "\nYear: %d", *meta.Year)
	}

	prompt := fmt.Sprintf(
		`You are a tabletop miniature game expert with detailed knowledge of boardgame and wargame box contents.

List the exact miniatures included in this specific product:

Game/System: %s%s
Box/Set: %s

Rules:
- Only list miniatures you are confident are actually in this box. Do not guess or approximate.
- If you are uncertain about a specific miniature, omit it rather than including it.
- If you do not recognise this product or are not confident in its contents, return an empty array.
- Each entry must have the exact miniature name as it appears in the rulebook or on the box.
- "quantity" is the actual number of that sculpt/unit included, not a default.

Respond with ONLY a JSON array. Each element:
- "name": exact miniature name (string)
- "unit_type": category such as "infantry", "cavalry", "monster", "hero", "vehicle", "terrain" (string)
- "quantity": exact count included in the box (integer)

Do not include any text outside the JSON array.`,
		game, extras.String(), set,
	)

	msg, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-6",
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

	raw := extractJSON(msg.Content[0].Text)

	var minis []MiniSuggestion
	if err := json.Unmarshal([]byte(raw), &minis); err != nil {
		return nil, fmt.Errorf("parse claude response: %w (raw: %s)", err, raw)
	}

	return minis, nil
}

var codeBlockRe = regexp.MustCompile("(?s)```(?:json)?\\s*(\\[.*?\\])\\s*```")

// extractJSON strips markdown code fences if present, otherwise returns the input trimmed.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if m := codeBlockRe.FindStringSubmatch(s); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return s
}

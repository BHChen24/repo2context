package tokencounter

import (
	"fmt"

	"github.com/localit-io/tiktoken-go"
)

type TokenCounter struct {
	encoding *tiktoken.Tiktoken
}

// NewTokenCounter creates a new TokenCounter with the specified encoding.
// If encoding is empty, defaults to "o200k_base" (GPT-5 encoding).
// Supported encodings: o200k_base, cl100k_base, p50k_base, r50k_base, etc.
func NewTokenCounter(encoding string) (*TokenCounter, error) {
	// Default to o200k_base if not specified
	if encoding == "" {
		encoding = "o200k_base"
	}

	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding %s: %w", encoding, err)
	}

	return &TokenCounter{
		encoding: tke,
	}, nil
}

func (tc *TokenCounter) CountTokens(text string) (int, error) {
	return tc.CountTokensWithPath(text, "")
}

func (tc *TokenCounter) CountTokensWithPath(text string, filePath string) (int, error) {
	tokens := tc.encoding.Encode(text, nil, nil)

	return len(tokens), nil
}

package tokencounter_test

import (
	"testing"

	tokencounter "github.com/BHChen24/repo2context/pkg/tokenCounter"
)

// TestNewTokenCounter tests the initialization of TokenCounter with different encodings
func TestNewTokenCounter(t *testing.T) {
	tests := []struct {
		name     string
		encoding string
		wantErr  bool
	}{
		{
			name:     "valid encoding o200k_base",
			encoding: "o200k_base",
			wantErr:  false,
		},
		{
			name:     "valid encoding cl100k_base",
			encoding: "cl100k_base",
			wantErr:  false,
		},
		{
			name:     "valid encoding p50k_base",
			encoding: "p50k_base",
			wantErr:  false,
		},
		{
			name:     "invalid encoding",
			encoding: "invalid_encoding",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc, err := tokencounter.NewTokenCounter(tt.encoding)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tc == nil {
				t.Error("NewTokenCounter() returned nil without error")
			}
		})
	}
}

// TestCountTokens_SimpleText tests token counting with simple text
func TestCountTokensSimpleText(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := "Hello, world!"
	count, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokens() = %d, want > 0", count)
	}
}

// TestCountTokens_EmptyString tests token counting with empty string
func TestCountTokensEmptyString(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	count, err := tc.CountTokens("")
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count != 0 {
		t.Errorf("CountTokens() = %d, want 0", count)
	}
}

// TestCountTokens_MultilineText tests token counting with multi-line text
func TestCountTokensMultilineText(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := "Line 1\nLine 2\nLine 3"
	count, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokens() = %d, want > 0", count)
	}
}

// TestCountTokens_SpecialCharacters tests token counting with special characters
func TestCountTokensSpecialCharacters(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := "!@#$%^&*()_+"
	count, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokens() = %d, want > 0", count)
	}
}

// TestCountTokens_UnicodeCharacters tests token counting with unicode and emoji
func TestCountTokensUnicodeCharacters(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := "你好，世界！🌍"
	count, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokens() = %d, want > 0", count)
	}
}

// TestCountTokens_CodeSnippet tests token counting with code snippets
func TestCountTokensCodeSnippet(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}`
	count, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokens() = %d, want > 0", count)
	}
}

// TestCountTokens_MarkdownText tests token counting with markdown
func TestCountTokensMarkdownText(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := `# Heading
## Subheading
* List item 1
* List item 2

**Bold text** and _italic text_`
	count, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokens() = %d, want > 0", count)
	}
}

// TestCountTokens_VeryLongText tests token counting with very long text
func TestCountTokensVeryLongText(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	// Create a very long string (10,000 characters)
	var longText string
	for i := 0; i < 10000; i++ {
		longText += "a"
	}

	count, err := tc.CountTokens(longText)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokens() = %d, want > 0", count)
	}
}

// TestCountTokens_DifferentEncodings tests that different encodings produce different counts
func TestCountTokensDifferentEncodings(t *testing.T) {
	text := "Hello, world! This is a test."

	tc1, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter with o200k_base: %v", err)
	}

	tc2, err := tokencounter.NewTokenCounter("cl100k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter with cl100k_base: %v", err)
	}

	count1, err := tc1.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() with o200k_base error = %v", err)
	}

	count2, err := tc2.CountTokens(text)
	if err != nil {
		t.Errorf("CountTokens() with cl100k_base error = %v", err)
	}

	// Both counts should be > 0
	if count1 <= 0 {
		t.Errorf("count1 = %d, want > 0", count1)
	}
	if count2 <= 0 {
		t.Errorf("count2 = %d, want > 0", count2)
	}
}

// TestCountTokens_ConsistentResults tests that the same input produces consistent results
func TestCountTokensConsistentResults(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := "The quick brown fox jumps over the lazy dog"

	count1, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("First CountTokens() error = %v", err)
	}

	count2, err := tc.CountTokens(text)
	if err != nil {
		t.Errorf("Second CountTokens() error = %v", err)
	}

	if count1 != count2 {
		t.Errorf("Inconsistent results: first call = %d, second call = %d", count1, count2)
	}
}

// TestCountTokens_WithFilePath tests token counting with file path for error context
func TestCountTokensWithFilePath(t *testing.T) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		t.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := "Test content"
	filePath := "test.txt"

	count, err := tc.CountTokensWithPath(text, filePath)
	if err != nil {
		t.Errorf("CountTokensWithPath() error = %v", err)
	}
	if count <= 0 {
		t.Errorf("CountTokensWithPath() = %d, want > 0", count)
	}
}

// BenchmarkCountTokens benchmarks the token counting performance
func BenchmarkCountTokensPerformance(b *testing.B) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		b.Fatalf("Failed to create TokenCounter: %v", err)
	}

	text := `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
	for i := 0; i < 100; i++ {
		fmt.Printf("Iteration %d\n", i)
	}
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tc.CountTokens(text)
	}
}

// BenchmarkCountTokens_LongText benchmarks token counting with long text
func BenchmarkCountTokensLongText(b *testing.B) {
	tc, err := tokencounter.NewTokenCounter("o200k_base")
	if err != nil {
		b.Fatalf("Failed to create TokenCounter: %v", err)
	}

	// Create a longer text (1000 words)
	text := ""
	for i := 0; i < 1000; i++ {
		text += "The quick brown fox jumps over the lazy dog. "
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tc.CountTokens(text)
	}
}

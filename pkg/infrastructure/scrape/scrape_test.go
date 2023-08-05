package scrape_test

import (
	"strings"
	"testing"

	"notify/pkg/infrastructure/scrape"

	"github.com/PuerkitoBio/goquery"
)

func TestGetFirstNChars(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		selector string
		n        int
		expected string
	}{
		{
			name:     "Simple English text",
			html:     `<div class="test">Hello, world!</div>`,
			selector: ".test",
			n:        11,
			expected: "Hello,world",
		},
		{
			name:     "Two English sentences",
			html:     `<div class="test">Hello, world!</div><div class="test">Goodbye, world!</div>`,
			selector: ".test",
			n:        5,
			expected: "Hello",
		},
		{
			name:     "Simple Japanese text",
			html:     `<div class="test">こんにちは、世界！</div>`,
			selector: ".test",
			n:        5,
			expected: "こんにちは",
		},
		{
			name:     "Mixed Japanese and English text",
			html:     `<div class="test">日向坂46の加藤史帆です。</div>`,
			selector: ".test",
			n:        8,
			expected: "日向坂46の加藤",
		},
		{
			name:     "Nested HTML tags",
			html:     `<div class="test"><p>日向坂46の加藤史帆です。</p></div>`,
			selector: ".test",
			n:        8,
			expected: "日向坂46の加藤",
		},
		{
			name:     "Parallel nested HTML tags",
			html:     `<div class="test"><p>日向坂46の</p><p>加藤史帆です。</p></div>`,
			selector: ".test",
			n:        8,
			expected: "日向坂46の加藤",
		},
		{
			name:     "ID selector",
			html:     `<div id="test">Hello, world!</div>`,
			selector: "#test",
			n:        5,
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatal(err)
			}

			result := scrape.GetFirstNChars(doc, tt.selector, tt.n)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

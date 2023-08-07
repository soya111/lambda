package scrape_test

import (
	"strings"
	"testing"

	"notify/pkg/infrastructure/scrape"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestExtractAndFormatTextFromElement(t *testing.T) {
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
			n:        12,
			expected: "Hello, world",
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
		{
			name:     "Invalid selector",
			html:     `<div class="test">Hello, world!</div>`,
			selector: ".invalid",
			n:        12,
			expected: "",
		},
		{
			name:     "Empty text",
			html:     `<div class="test"></div>`,
			selector: ".test",
			n:        5,
			expected: "",
		},
		{
			name:     "Whitespace only text",
			html:     `<div class="test">     </div>`,
			selector: ".test",
			n:        5,
			expected: "",
		},
		{
			name:     "Text with leading and trailing whitespaces",
			html:     `  <div class="test">   Hello   </div>  `,
			selector: ".test",
			n:        8,
			expected: "Hello",
		},
		{
			name:     "MaxLength set to 0",
			html:     `<div class="test">Hello, world!</div>`,
			selector: ".test",
			n:        0,
			expected: "",
		},
		{
			name:     "Text length same as MaxLength",
			html:     `<div class="test">Hello</div>`,
			selector: ".test",
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
			opt := scrape.TextExtractionOptions{
				MaxLength: tt.n,
			}

			result, err := scrape.ExtractAndFormatTextFromElement(doc, tt.selector, opt)
			assert.NoError(t, err)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

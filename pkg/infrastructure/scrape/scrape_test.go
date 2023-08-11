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
		name      string
		html      string
		selector  string
		n         int
		expected  string
		expectErr bool
	}{
		{
			name:      "Simple English text",
			html:      `<div class="test">Hello, world!</div>`,
			selector:  ".test",
			n:         12,
			expected:  "Hello, world",
			expectErr: false,
		},
		{
			name:      "Two English sentences",
			html:      `<div class="test">Hello, world!</div><div class="test">Goodbye, world!</div>`,
			selector:  ".test",
			n:         5,
			expected:  "Hello",
			expectErr: false,
		},
		{
			name:      "Simple Japanese text",
			html:      `<div class="test">こんにちは、世界！</div>`,
			selector:  ".test",
			n:         5,
			expected:  "こんにちは",
			expectErr: false,
		},
		{
			name:      "Mixed Japanese and English text",
			html:      `<div class="test">日向坂46の加藤史帆です。</div>`,
			selector:  ".test",
			n:         8,
			expected:  "日向坂46の加藤",
			expectErr: false,
		},
		{
			name:      "Nested HTML tags",
			html:      `<div class="test"><p>日向坂46の加藤史帆です。</p></div>`,
			selector:  ".test",
			n:         8,
			expected:  "日向坂46の加藤",
			expectErr: false,
		},
		{
			name:      "Parallel nested HTML tags",
			html:      `<div class="test"><p>日向坂46の</p><p>加藤史帆です。</p></div>`,
			selector:  ".test",
			n:         8,
			expected:  "日向坂46の加藤",
			expectErr: false,
		},
		{
			name:      "ID selector",
			html:      `<div id="test">Hello, world!</div>`,
			selector:  "#test",
			n:         5,
			expected:  "Hello",
			expectErr: false,
		},
		{
			name:      "Invalid selector",
			html:      `<div class="test">Hello, world!</div>`,
			selector:  ".invalid",
			n:         12,
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Empty text",
			html:      `<div class="test"></div>`,
			selector:  ".test",
			n:         5,
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Whitespace only text",
			html:      `<div class="test">     </div>`,
			selector:  ".test",
			n:         5,
			expected:  " ",
			expectErr: false,
		},
		{
			name:      "Text with leading and trailing whitespaces",
			html:      `  <div class="test">   Hello   </div>  `,
			selector:  ".test",
			n:         8,
			expected:  " Hello ",
			expectErr: false,
		},
		{
			name:      "MaxLength set to 0",
			html:      `<div class="test">Hello, world!</div>`,
			selector:  ".test",
			n:         0,
			expected:  "",
			expectErr: false,
		},
		{
			name:      "Text length same as MaxLength",
			html:      `<div class="test">Hello</div>`,
			selector:  ".test",
			n:         5,
			expected:  "Hello",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatal(err)
			}
			opt := scrape.TextExtractionOptions{
				MaxLength: tt.n,
			}

			result, err := scrape.ExtractAndFormatTextFromElement(doc, tt.selector, opt)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result, "Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

package scrape

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

// Document structを返す
// エラー握りつぶし
func GetDocumentFromURL(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d error: %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	// Body内を読み取り
	buffer, _ := io.ReadAll(res.Body)

	// 文字コードを判定
	detector := chardet.NewTextDetector()
	detectResult, _ := detector.DetectBest(buffer)

	// 文字コードの変換
	bufferReader := bytes.NewReader(buffer)
	reader, _ := charset.NewReaderLabel(detectResult.Charset, bufferReader)

	// HTMLをパース
	document, _ := goquery.NewDocumentFromReader(reader)

	return document, nil
}

// TextExtractionOptions defines the set of options for the ExtractAndFormatTextFromElement function.
type TextExtractionOptions struct {
	MaxLength       int  // The maximum number of characters to extract.
	IncludeNewlines bool // Whether to include newlines in the extracted text.
	AppendEllipsis  bool // Whether to append "..." to the text if it exceeds the MaxLength.
}

// ExtractAndFormatTextFromElement extracts text from a specific HTML element defined by the given selector.
// After extracting, it formats the text by removing unnecessary white spaces, preserving meaningful spaces,
// and optionally includes newlines or appends ellipsis based on the provided options.
//
// Parameters:
// - doc: The parsed HTML document.
// - selector: The CSS selector to identify the desired HTML element.
// - opts: The options for extraction and formatting.
//
// Returns:
// - A formatted string extracted from the HTML element.
// - An error if any issues occur during extraction.
func ExtractAndFormatTextFromElement(doc *goquery.Document, selector string, opts TextExtractionOptions) (string, error) {
	text := doc.Find(selector).Text()
	if text == "" {
		return "", fmt.Errorf("no text found for selector: %s", selector)
	}

	// Use regular expression to remove excessive white spaces
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	if !opts.IncludeNewlines {
		text = strings.ReplaceAll(text, "\n", "")
	}

	// Convert to []rune to handle multi-byte characters
	runeText := []rune(text)

	if len(runeText) > opts.MaxLength {
		if opts.AppendEllipsis {
			runeText = runeText[:opts.MaxLength-3] // Reserve 3 characters for "..."
			text = string(runeText) + "..."
		} else {
			runeText = runeText[:opts.MaxLength]
			text = string(runeText)
		}
	}
	return text, nil
}

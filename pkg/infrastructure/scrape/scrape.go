package scrape

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"

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

// GetFirstNChars is a function that extracts text from a specific HTML element,
// removes unnecessary whitespaces and line breaks, and returns the first N characters.
func GetFirstNChars(doc *goquery.Document, selector string, n int) string {
	text := doc.Find(selector).Text()

	// Use regular expression to remove white spaces
	re := regexp.MustCompile(`\s`)
	text = re.ReplaceAllString(text, "")

	// Convert to []rune to handle multi-byte characters
	runeText := []rune(text)

	if len(runeText) > n {
		runeText = runeText[:n]
	}
	return string(runeText)
}

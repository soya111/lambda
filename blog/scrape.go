package blog

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type ScraperInterface interface {
	GetLatestDiaries() []*Diary
}

type Scraper struct {
}

// Document structを返す
func GetDocumentFromURL(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	defer res.Body.Close()

	// Body内を読み取り
	buffer, _ := ioutil.ReadAll(res.Body)

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

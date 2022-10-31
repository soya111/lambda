package blog

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

type NogizakaScraper struct {
	Scraper
}

// 最新の記事を調べる
func (s *NogizakaScraper) GetAndPostLatestDiaries() []*Diary {
	latestDiaries := s.getLatestDiaries()

	res := []*Diary{}
	for _, d := range latestDiaries {
		diary, err := GetDiary("nogizaka_blog", d.MemberName, d.Id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}

		// Dynamoにデータがない場合
		if diary.Id == 0 {
			fmt.Printf("%s %s %s\n%s\n", d.Date, d.Title, d.MemberName, d.Url)
			if err := PostDiary("nogizaka_blog", d); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}
			res = append(res, d)
		}
	}

	return res
}

func (s *NogizakaScraper) getLatestDiaries() []*Diary {
	url := "https://blog.nogizaka46.com/atom.xml"

	loc, _ := time.LoadLocation("Asia/Tokyo")
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	articles := feed.Items

	res := []*Diary{}

	for _, article := range articles {
		reader := strings.NewReader(article.Content)
		document, _ := goquery.NewDocumentFromReader(reader)
		imgs := document.Find("img")
		srcs := []string{}
		imgs.Each(func(i int, s *goquery.Selection) {
			for _, attr := range s.Nodes[0].Attr {
				if attr.Key == "src" {
					srcs = append(srcs, attr.Val)
				}
			}
		})

		url := article.Link
		title := article.Title
		name := article.Author.Name
		diaryId, _ := s.getIdFromURL(url)
		updatedAt, _ := time.ParseInLocation("2006-01-02T15:04:05Z", article.Updated, loc)

		newDiary := NewDiary(url, title, name, updatedAt, diaryId)
		newDiary.Images = srcs
		res = append(res, newDiary)

	}

	return res
}

func (*NogizakaScraper) getIdFromURL(url string) (int, error) {
	id, err := strconv.Atoi(strings.Split(strings.Split(url, "/")[6], ".")[0])
	if err != nil {
		return 0, err
	}
	return id, nil
}

// blog中の全画像を取得
func (s *NogizakaScraper) GetImages(diary *Diary) []string {
	return diary.Images
}

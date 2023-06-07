package blog

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HinatazakaScraper struct {
	Scraper
	repo DiaryRepository
}

func NewHinatazakaScraper(repo DiaryRepository) *HinatazakaScraper {
	return &HinatazakaScraper{repo: repo}
}

// 最新の記事を取得する
func (s *HinatazakaScraper) GetLatestDiaries() ([]*Diary, error) {
	latestDiaries := s.scrapeLatestDiaries()

	res := []*Diary{}
	for _, d := range latestDiaries {
		_, err := s.repo.GetDiary(d.MemberName, d.Id)
		if err != nil {
			// Check if the error is a "not found" error.
			if err == ErrDiaryNotFound {
				// The item is not in the database, so it's a new diary.
				res = append(res, d)
				continue
			}
			// Some other error occurred.
			return nil, err
		}
	}

	return res, nil
}

// 記事を保存する
func (s *HinatazakaScraper) PostDiaries(diaries []*Diary) error {
	for _, d := range diaries {
		fmt.Printf("%s %s %s\n%s\n", d.Date, d.Title, d.MemberName, d.Url)
		if err := s.repo.PostDiary(d); err != nil {
			return err
		}
	}

	return nil
}

func (s *HinatazakaScraper) scrapeLatestDiaries() []*Diary {
	rootURL := "https://www.hinatazaka46.com"
	url := "https://www.hinatazaka46.com/s/official/diary/member/list?ima=0000"
	document, err := GetDocumentFromURL(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	res := []*Diary{}
	articles := document.Find(".p-blog-article")

	articles.Each(func(i int, sl *goquery.Selection) {
		updateAt, err := time.Parse("2006.1.2 15:04 (MST)", strings.TrimSpace(sl.Find(".c-blog-article__date").Text())+" (JST)")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		href := sl.Find(".c-button-blog-detail").Nodes[0].Attr[1].Val
		title := strings.TrimSpace(sl.Find(".c-blog-article__title").Text())
		name := strings.TrimSpace(sl.Find(".c-blog-article__name").Text())
		diaryId := s.getIdFromHref(href)

		newDiary := NewDiary(rootURL+href, title, name, updateAt, diaryId)
		res = append(res, newDiary)
	})

	return reverse(res)
}

func (*HinatazakaScraper) getIdFromHref(href string) int {
	id, err := strconv.Atoi(strings.Split(strings.Split(href, "/")[5], "?")[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return id
}

// blog中の全画像を取得
func (s *HinatazakaScraper) GetImages(diary *Diary) []string {
	document, err := GetDocumentFromURL(diary.Url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	article := document.Find(".c-blog-article__text")
	img := article.Find("img")
	srcs := []string{}
	img.Each(func(i int, s *goquery.Selection) {
		for _, attr := range s.Nodes[0].Attr {
			if attr.Key == "src" {
				srcs = append(srcs, attr.Val)
			}
		}
	})

	return srcs
}

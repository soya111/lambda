package blog

import (
	"fmt"
	"notify/pkg/infrastructure/scrape"
	"notify/pkg/model"
	"notify/pkg/slices"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	RootURL = "https://www.hinatazaka46.com"
	TimeFmt = "2006.1.2 15:04 (MST)"
)

type HinatazakaScraper struct {
	Scraper
	repo model.DiaryRepository
}

func NewHinatazakaScraper(repo model.DiaryRepository) *HinatazakaScraper {
	return &HinatazakaScraper{repo: repo}
}

// 最新の記事を取得する
func (s *HinatazakaScraper) GetLatestDiaries() ([]*ScrapedDiary, error) {
	latestDiaries, err := s.scrapeLatestDiaries()
	if err != nil {
		return nil, err
	}

	res := []*ScrapedDiary{}
	for _, diary := range latestDiaries {
		_, err := s.repo.GetDiary(diary.MemberName, diary.Id)
		if err != nil {
			// Check if the error is a "not found" error.
			if err == model.ErrDiaryNotFound {
				// The item is not in the database, so it's a new diary.
				res = append(res, diary)
				continue
			}
			// Some other error occurred.
			return nil, err
		}
	}

	return res, nil
}

// 記事を保存する
func (s *HinatazakaScraper) PostDiaries(diaries []*ScrapedDiary) error {
	for _, d := range diaries {
		diary := ConvertScrapedDiaryToDiary(d)
		fmt.Printf("%s %s %s\n%s\n", diary.Date, diary.Title, diary.MemberName, diary.Url)
		if err := s.repo.PostDiary(diary); err != nil {
			return err
		}
	}

	return nil
}

func (s *HinatazakaScraper) scrapeLatestDiaries() ([]*ScrapedDiary, error) {
	url := fmt.Sprintf("%s/s/official/diary/member/list?ima=0000", RootURL)
	document, err := scrape.GetDocumentFromURL(url)
	if err != nil {
		return nil, err
	}

	res := []*ScrapedDiary{}
	articles := document.Find(".p-blog-article")

	articles.Each(func(i int, sl *goquery.Selection) {
		diary, err := s.parseDiaryFromSelection(sl)
		if err != nil {
			fmt.Printf("error parsing diary from selection: %v\n", err)
			return
		}
		res = append(res, diary)
	})

	return slices.Reverse(res), nil
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
func (s *HinatazakaScraper) GetImages(document *goquery.Document) []string {
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

func (s *HinatazakaScraper) GetMemberIcon(document *goquery.Document) string {
	var iconUrl = "https://natalie.mu/music/news/472084"
	// Find the div with class .c-blog-member__icon
	document.Find(".c-blog-member__icon").Each(func(i int, s *goquery.Selection) {
		// Get the style attribute
		style, exists := s.Attr("style")
		if exists {
			// Split the style string into 2 parts: "background-image:url(" and the url with ");" at the end
			splittedStyle := strings.Split(style, "(")
			if len(splittedStyle) == 2 {
				// Remove the ");" from the end of the second part of splittedStyle to get the url
				iconUrl = strings.TrimSuffix(splittedStyle[1], ");")
			}
		}
	})
	return iconUrl
}

// 各メンバーごとの最新記事を取得する
func (s *HinatazakaScraper) GetLatestDiaryByMember(memberName string) (*ScrapedDiary, error) {
	memberNumber, err := model.GetMemberNumber(memberName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/s/official/diary/member/list?ima=0000&ct=%s", RootURL, memberNumber)

	document, err := scrape.GetDocumentFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get document from url %s: %w", url, err)
	}

	article := document.Find(".p-blog-article").First()

	diary, err := s.parseDiaryFromSelection(article)

	if err != nil {
		return nil, fmt.Errorf("failed to parse diary from selection: %w", err)
	}

	return diary, nil
}

func (s *HinatazakaScraper) parseDiaryFromSelection(sl *goquery.Selection) (*ScrapedDiary, error) {
	href := sl.Find(".c-button-blog-detail").Nodes[0].Attr[1].Val
	title := strings.TrimSpace(sl.Find(".c-blog-article__title").Text())
	name := strings.TrimSpace(sl.Find(".c-blog-article__name").Text())
	diaryId := s.getIdFromHref(href)

	date, err := time.Parse(TimeFmt, strings.TrimSpace(sl.Find(".c-blog-article__date").Text())+" (JST)")
	if err != nil {
		fmt.Println(err)
	}

	images := s.GetImages(&goquery.Document{Selection: sl})
	lead := scrape.GetFirstNChars(&goquery.Document{Selection: sl}, ".c-blog-article__text", 50)
	iconUrl := s.GetMemberIcon(&goquery.Document{Selection: sl})

	diary := NewScrapedDiary(RootURL+href, title, name, date, diaryId, images, lead, iconUrl)

	return diary, nil
}

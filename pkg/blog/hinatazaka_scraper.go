package blog

import (
	"context"
	"fmt"
	"notify/pkg/infrastructure/scrape"
	"notify/pkg/model"

	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	RootURL = "https://www.hinatazaka46.com"
	TimeFmt = "2006.1.2 15:04 (MST)"
)

// HinatazakaScraper scrapes Hinatazaka46's blog
type HinatazakaScraper struct {
	Scraper
}

// NewHinatazakaScraper returns a new HinatazakaScraper
func NewHinatazakaScraper() *HinatazakaScraper {
	return &HinatazakaScraper{}
}

// ScrapeLatestDiaries scrapes the latest diaries in order of old
func (s *HinatazakaScraper) ScrapeLatestDiaries(ctx context.Context) ([]*ScrapedDiary, error) {
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
		memberId, err := model.GetMemberId(diary.MemberName)
		if err != nil {
			fmt.Printf("error getting member id: %v\n", err)
		}
		diary.SetMemberIcon(s.GetIconURLByID(document, memberId))
		res = append(res, diary)
	})

	slices.Reverse(res)

	return res, nil
}

func (*HinatazakaScraper) getIdFromHref(href string) int {
	id, err := strconv.Atoi(strings.Split(strings.Split(href, "/")[5], "?")[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return id
}

// GetImages returns the list of image URLs in the blog
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

// GetIconURLByID returns the icon URL of the specified member
func (s *HinatazakaScraper) GetIconURLByID(document *goquery.Document, memberID string) string {
	var iconUrl = "https://cdn.hinatazaka46.com/images/14/14d/a9bac831ed1e6a4fdd93c4271aa8a.jpg"

	query := fmt.Sprintf(`.p-blog-face__list[data-member="%s"]`, memberID)
	div := document.Find(query).First().Find(".c-blog-face__item")

	// Get the style attribute
	style, exists := div.Attr("style")
	if exists {
		// Split the style string to extract the URL
		splittedStyle := strings.Split(style, "(")
		if len(splittedStyle) == 2 {
			// Remove the ");" from the end of the second part of splittedStyle to get the url
			iconUrl = strings.TrimSuffix(splittedStyle[1], ");")
		}
	}

	return iconUrl
}

// GetLatestDiaryByMember returns the latest diary of the specified member
func (s *HinatazakaScraper) GetLatestDiaryByMember(memberName string) (*ScrapedDiary, error) {
	memberId, err := model.GetMemberId(memberName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/s/official/diary/member/list?ima=0000&ct=%s", RootURL, memberId)

	document, err := scrape.GetDocumentFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get document from url %s: %w", url, err)
	}

	article := document.Find(".p-blog-article").First()

	diary, err := s.parseDiaryFromSelection(article)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diary from selection: %w", err)
	}

	diary.SetMemberIcon(s.GetIconURLByID(document, memberId))

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

	opt := scrape.TextExtractionOptions{
		MaxLength:       50,
		IncludeNewlines: false,
		AppendEllipsis:  true,
	}

	lead, err := scrape.ExtractAndFormatTextFromElement(&goquery.Document{Selection: sl}, ".c-blog-article__text", opt)
	if err != nil {
		fmt.Printf("error extracting text from element: %v\n", err)
	}

	diary := NewScrapedDiary(RootURL+href, title, name, date, diaryId, images, lead)

	return diary, nil
}

package blog

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SakurazakaScraper struct {
	Scraper
}

// 最新の記事を調べる
func (s *SakurazakaScraper) GetAndPostLatestDiaries() []*Diary {
	latestDiaries := s.getLatestDiaries()

	res := []*Diary{}
	for _, d := range latestDiaries {
		diary, err := GetDiary("sakurazaka_blog", d.MemberName, d.Id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}

		// Dynamoにデータがない場合
		if diary.Id == 0 {
			fmt.Printf("%s %s %s\n%s\n", d.Date, d.Title, d.MemberName, d.Url)
			if err := PostDiary("sakurazaka_blog", d); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}
			res = append(res, d)
		}
	}

	return res
}

func (s *SakurazakaScraper) getLatestDiaries() []*Diary {
	rootURL := "https://sakurazaka46.com"
	url := "https://sakurazaka46.com/s/s46/diary/blog/list?ima=0501"
	document, err := GetDocumentFromURL(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	res := []*Diary{}
	articles := document.Find(".com-blog-part > li")

	articles.Each(func(i int, sl *goquery.Selection) {
		href := sl.Find("a").Nodes[0].Attr[0].Val
		title := sl.Find(".title").Text()
		name := sl.Find("p.name").Text()
		diaryId := s.getIdFromHref(href)

		// 更新日時が取れないので現在時刻に
		newDiary := NewDiary(rootURL+href, title, name, time.Now(), diaryId)
		res = append(res, newDiary)
	})

	return reverse(res)
}

func (*SakurazakaScraper) getIdFromHref(href string) int {
	id, err := strconv.Atoi(strings.Split(strings.Split(href, "/")[5], "?")[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return id
}

// blog中の全画像を取得
func (s *SakurazakaScraper) GetImages(diary *Diary) []string {
	rootURL := "https://sakurazaka46.com"
	document, err := GetDocumentFromURL(diary.Url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	date, _ := time.Parse("2006/1/2 15:04 (MST)", document.Find(".wf-a").Text()+" (JST)")
	diary.Date = date.Format("2006.1.2 15:04 (MST)")

	article := document.Find(".box-article")
	img := article.Find("img")
	srcs := []string{}
	img.Each(func(i int, s *goquery.Selection) {
		for _, attr := range s.Nodes[0].Attr {
			if attr.Key == "src" {
				srcs = append(srcs, rootURL+attr.Val)
			}
		}
	})

	return srcs
}

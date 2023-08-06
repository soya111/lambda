package blog

import (
	"notify/pkg/model"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Scraper interface {
	GetLatestDiaries() ([]*ScrapedDiary, error)
	PostDiaries([]*ScrapedDiary) error
	GetImages(*goquery.Document) []string
	GetMemberIcon(*goquery.Document) string
}

type ScrapedDiary struct {
	Url        string   `json:"url"`
	Title      string   `json:"title"`
	MemberName string   `json:"member_name"`
	Date       string   `json:"date"`
	Id         int      `json:"diary_id"`
	Images     []string `json:"images"`
	Lead       string   `json:"lead"`
	MemberIcon string   `json:"member_icon"`
}

func NewScrapedDiary(url, title, memberName string, date time.Time, id int, images []string, lead string, memberIcon string) *ScrapedDiary {
	return &ScrapedDiary{url, title, memberName, date.Format(TimeFmt), id, images, lead, memberIcon}
}

// ScrapedDiaryをDiaryに変換する
func ConvertScrapedDiaryToDiary(s *ScrapedDiary) *model.Diary {
	return &model.Diary{
		Url:        s.Url,
		Title:      s.Title,
		MemberName: s.MemberName,
		Date:       s.Date,
		Id:         s.Id,
		Images:     s.Images,
	}
}

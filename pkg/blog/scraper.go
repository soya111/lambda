package blog

import (
	"notify/pkg/model"
	"time"
)

type Scraper interface {
	ScrapeLatestDiaries() ([]*ScrapedDiary, error)
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

func NewScrapedDiary(url, title, memberName string, date time.Time, id int, images []string, lead string) *ScrapedDiary {
	return &ScrapedDiary{url, title, memberName, date.Format(TimeFmt), id, images, lead, ""}
}

// SetMemberIconはScrapedDiaryのMemberIconを設定します。
func (sd *ScrapedDiary) SetMemberIcon(iconUrl string) {
	sd.MemberIcon = iconUrl
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

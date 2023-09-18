package blog

import (
	"context"
	"notify/pkg/model"
	"time"
)

// Scraper is an interface for scraping blogs.
type Scraper interface {
	ScrapeLatestDiaries(context.Context) ([]*ScrapedDiary, error)
}

// ScrapedDiary is a struct that represents a scraped diary.
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

// NewScrapedDiary creates a new ScrapedDiary.
func NewScrapedDiary(url, title, memberName string, date time.Time, id int, images []string, lead string) *ScrapedDiary {
	return &ScrapedDiary{url, title, memberName, date.Format(model.TimeFmt), id, images, lead, ""}
}

// SetMemberIcon sets the member icon url.
func (sd *ScrapedDiary) SetMemberIcon(iconUrl string) {
	sd.MemberIcon = iconUrl
}

// IsNewDiary returns true if Date is within 24h.
func (sd *ScrapedDiary) IsNew() bool {
	timeTypeDate, _ := time.Parse(model.TimeFmt, sd.Date)
	timeDifference := time.Since(timeTypeDate)
	judgment := 24 * time.Hour

	return timeDifference <= judgment
}

// ConvertScrapedDiaryToDiary converts ScrapedDiary to Diary.
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

package blog

import (
	"notify/pkg/model"

	"github.com/PuerkitoBio/goquery"
)

type Scraper interface {
	GetLatestDiaries() ([]*model.Diary, error)
	PostDiaries([]*model.Diary) error
	GetImages(*goquery.Document) []string
	GetMemberIcon(*goquery.Document) string
}

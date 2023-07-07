package blog

import "notify/pkg/model"

type Scraper interface {
	GetLatestDiaries() ([]*model.Diary, error)
	PostDiaries([]*model.Diary) error
	GetImages(*model.Diary) []string
}

package blog

import "notify/pkg/model"

func reverse[T any](a []T) []T {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - i - 1
		a[i], a[j] = a[j], a[i]
	}
	return a
}

type Scraper interface {
	GetLatestDiaries() ([]*model.Diary, error)
	PostDiaries([]*model.Diary) error
	GetImages(*model.Diary) []string
}

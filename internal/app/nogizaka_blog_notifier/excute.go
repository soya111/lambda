package notifier

import (
	"notify/internal/pkg/blog"
	"notify/internal/pkg/database"
	"notify/internal/pkg/line"
	"notify/internal/pkg/notifier"
)

func ExcuteFunction() {
	db, err := database.NewDynamo()
	if err != nil {
		panic(err)
	}
	notifier.Excute(&blog.NogizakaScraper{}, line.NewLinebot(), db)
}

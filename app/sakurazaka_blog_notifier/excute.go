package notifier

import (
	"notify/pkg/blog"
	"notify/pkg/database"
	"notify/pkg/line"
	"notify/pkg/notifier"
)

func ExcuteFunction() {
	db, err := database.NewDynamo()
	if err != nil {
		panic(err)
	}
	notifier.Excute(&blog.SakurazakaScraper{}, line.NewLinebot(), db)
}

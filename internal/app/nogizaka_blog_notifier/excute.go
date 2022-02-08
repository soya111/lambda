package notifier

import (
	"notify/internal/pkg/blog"
	"notify/internal/pkg/notifier"
)

func ExcuteFunction() {
	notifier.Excute(&blog.NogizakaScraper{})
}

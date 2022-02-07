package notifier

import (
	"fmt"
	"os"

	"notify/internal/pkg/blog"
	"notify/internal/pkg/database"
	"notify/internal/pkg/notifier"
)

func ExcuteFunction() {
	to, err := database.GetDestination()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	notifier.Excute(&blog.NogizakaScraper{}, to)
}

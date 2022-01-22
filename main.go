package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"main/blog"
	"main/line"
)

func excute(s blog.ScraperInterface) {
	to := []string{os.Getenv("ME")}
	latestDiaries := s.GetLatestDiaries()
	for _, diary := range latestDiaries {
		images := diary.GetImages()
		text := fmt.Sprintf("%s %s %s\n%s", diary.Date, diary.Title, diary.MemberName, diary.Url)
		line.PushTextMessages(to, text)
		line.PushFlexImagesMessage(to, images)
	}
}

func excuteFunction() {
	excute(&blog.Scraper{})
}

func init() {
	// set timezone
	time.Local = time.FixedZone("JST", 9*60*60)
}

func main() {
	lambda.Start(excuteFunction)
}

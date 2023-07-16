package main

import (
	"context"
	"fmt"
	"log"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/dynamodb"
	"notify/pkg/infrastructure/line"
	"notify/pkg/notifier"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var bot *line.Linebot
var sess *session.Session

func init() {
	// set timezone
	time.Local = time.FixedZone("JST", 9*60*60)

	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")
	var err error
	bot, err = line.NewLinebot(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	sess, err = session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	lambda.Start(func(ctx context.Context) error {
		diary := dynamodb.NewDynamoDiaryRepository(sess, "hinatazaka_blog")
		scraper := blog.NewHinatazakaScraper(diary)
		subscriber := dynamodb.NewDynamoSubscriberRepository(sess)

		err := notifier.Excute(ctx, scraper, bot, subscriber)
		if err != nil {
			return fmt.Errorf("ApplicationError in Excute function: %v", err)
		}

		return nil
	})
}

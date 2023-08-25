package main

import (
	"context"
	"fmt"
	"log"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/dynamodb"
	"notify/pkg/infrastructure/line"
	"notify/pkg/logging"
	"notify/pkg/notifier"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"go.uber.org/zap"
)

var (
	bot    *line.Linebot
	sess   *session.Session
	logger *zap.Logger
)

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

	sess = session.Must(session.NewSession())

	logger = logging.InitializeLogger()
}

func main() {
	lambda.Start(func(ctx context.Context) error {
		ctx = logging.ContextWithLogger(ctx, logger)
		db := dynamo.New(sess)
		diary := dynamodb.NewDiaryRepository(db, "hinatazaka_blog")
		scraper := blog.NewHinatazakaScraper()
		subscriber := dynamodb.NewSubscriberRepository(sess)

		notifier := notifier.NewNotifier(scraper, bot, subscriber, diary)
		err := notifier.Execute(ctx)
		if err != nil {
			// TODO: spelling error
			return fmt.Errorf("ApplicationError in Excute function: %v", err)
		}

		return nil
	})
}

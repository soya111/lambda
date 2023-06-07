package main

import (
	"context"
	"log"
	"notify/pkg/blog"
	"notify/pkg/database"
	"notify/pkg/line"
	"notify/pkg/notifier"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var bot *line.Linebot
var sess *session.Session
var db *database.Dynamo

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

	db, err = database.NewDynamo()
	if err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(func() {
		ctx := context.Background()
		repo := blog.NewDynamoDiaryRepository(sess, "hinatazaka_blog")
		scraper := blog.NewHinatazakaScraper(repo)
		notifier.Excute(ctx, scraper, bot, db)
	})
}

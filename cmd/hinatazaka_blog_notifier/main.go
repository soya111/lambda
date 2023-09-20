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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"go.uber.org/zap"
)

var (
	mainFunc func()
	bot      *line.Linebot
	sess     *session.Session
	logger   *zap.Logger
	db       *dynamo.DB
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

	if runningLocally() {
		const (
			// local dynamodb settings
			AWS_REGION      = "ap-northeast-1"
			DYNAMO_ENDPOINT = "http://dynamodb-local:8000"
		)
		db = dynamo.New(sess, &aws.Config{
			Region:      aws.String(AWS_REGION),
			Endpoint:    aws.String(DYNAMO_ENDPOINT),
			Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
		})
		mainFunc = runAsLocal
	} else {
		db = dynamo.New(sess)
		mainFunc = runAsLambda
	}
}

// Execute does the main logic
func Execute(ctx context.Context) error {
	ctx = logging.ContextWithLogger(ctx, logger)
	diary := dynamodb.NewDiaryRepository(db, "hinatazaka_blog")
	scraper := blog.NewHinatazakaScraper()
	subscriber := dynamodb.NewSubscriberRepository(db)

	notifier := notifier.NewNotifier(scraper, bot, subscriber, diary)
	err := notifier.Execute(ctx)
	if err != nil {
		return fmt.Errorf("ApplicationError in Execute function: %v", err)
	}

	return nil
}

// lambdaHandler is the AWS Lambda handler
func lambdaHandler(ctx context.Context) error {
	return Execute(ctx)
}

// runAsLambda runs the application as an AWS Lambda function
func runAsLambda() {
	lambda.Start(lambdaHandler)
}

// runAsLocal runs the application as a local server
func runAsLocal() {
	ctx := context.Background()
	err := Execute(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mainFunc()
}

// runningLocally checks if the application is running locally or in AWS Lambda
func runningLocally() bool {
	_, isLocal := os.LookupEnv("IS_LOCAL")
	return isLocal
}

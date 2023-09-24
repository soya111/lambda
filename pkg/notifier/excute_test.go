package notifier

import (
	"context"
	"os"
	"testing"

	"zephyr/pkg/blog"
	"zephyr/pkg/infrastructure/dynamodb"
	"zephyr/pkg/infrastructure/line"
	"zephyr/pkg/logging"
	"zephyr/pkg/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Skip("skipping this test for now")
	err := godotenv.Load("../.env")
	assert.NoError(t, err)

	// LINE settings
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")
	me := os.Getenv("ME")
	bot, err := line.NewLinebot(channelSecret, channelToken)
	assert.NoError(t, err)

	// local dynamodb settings
	DYNAMO_ENDPOINT := "http://localhost:8000"
	AWS_REGION := "ap-northeast-1"

	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(DYNAMO_ENDPOINT),
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	})
	assert.NoError(t, err)
	db := dynamo.New(sess)

	subscriber := dynamodb.NewSubscriberRepository(db)
	name := "高瀬愛奈"
	err = subscriber.Subscribe(model.Subscriber{MemberName: name, UserId: me})
	assert.NoError(t, err)
	defer func() {
		err := subscriber.Unsubscribe(name, me)
		assert.NoError(t, err)
	}()

	diary := dynamodb.NewDiaryRepository(db, "hinatazaka_blog")

	// Scraper settings
	scraper := blog.NewHinatazakaScraper()

	logger := logging.InitializeLogger()
	ctx := logging.ContextWithLogger(context.Background(), logger)

	notifier := NewNotifier(scraper, bot, subscriber, diary)
	err = notifier.Execute(ctx)
	assert.NoError(t, err)
}

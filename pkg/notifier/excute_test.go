package notifier

import (
	"context"
	"os"
	"testing"

	"notify/pkg/blog"
	"notify/pkg/infrastructure/dynamodb"
	"notify/pkg/infrastructure/line"
	"notify/pkg/logging"
	"notify/pkg/model"

	"github.com/aws/aws-sdk-go/aws"
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
	DYNAMO_REGION := "ap-northeast-1"
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(DYNAMO_ENDPOINT),
		Region:   aws.String(DYNAMO_REGION),
	})
	assert.NoError(t, err)

	subscriber := dynamodb.NewSubscriberRepository(sess)
	name := "高瀬愛奈"
	err = subscriber.Subscribe(model.Subscriber{MemberName: name, UserId: me})
	assert.NoError(t, err)
	defer func() {
		err := subscriber.Unsubscribe(name, me)
		assert.NoError(t, err)
	}()

	diary := dynamodb.NewDiaryRepository(dynamo.New(sess), "hinatazaka_blog")

	// Scraper settings
	scraper := blog.NewHinatazakaScraper()

	logger := logging.InitializeLogger()
	ctx := logging.ContextWithLogger(context.Background(), logger)

	notifier := NewNotifier(scraper, bot, subscriber, diary)
	err = notifier.Execute(ctx)
	assert.NoError(t, err)
}

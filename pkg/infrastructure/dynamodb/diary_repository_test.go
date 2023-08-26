package dynamodb

import (
	"fmt"
	"notify/pkg/model"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
)

// docker compose up -dをしてからテストを実行する
func TestDiaryRepository(t *testing.T) {
	t.Skip("skipping this test for now")
	AWS_REGION := "ap-northeast-1"
	DYNAMO_ENDPOINT := "http://localhost:8000"

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Endpoint:    aws.String(DYNAMO_ENDPOINT),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	})
	if err != nil {
		panic(err)
	}

	tableName := "hinatazaka_blog"

	db := dynamo.New(sess)
	repo := NewDiaryRepository(db, tableName)

	diary := model.NewDiary("https://www.hinatazaka46.com/s/official/diary/detail/34467", "タイトル", "小坂菜緒", time.Now(), 1)
	err = repo.PostDiary(diary)
	if err != nil {
		t.Fatal(err)
	}

	diary, err = repo.GetDiary("小坂菜緒", 1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", diary)

	assert.Equal(t, "タイトル", diary.Title)
	assert.Equal(t, "小坂菜緒", diary.MemberName)
}

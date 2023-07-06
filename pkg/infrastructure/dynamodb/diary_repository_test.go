package dynamodb

import (
	"fmt"
	"notify/pkg/model"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// docker compose up -dをしてからテストを実行する
func TestDynamoDiaryRepository(t *testing.T) {
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

	repo := NewDynamoDiaryRepository(sess, tableName)
	db := repo.db

	err = db.Table(tableName).DeleteTable().Run()
	if err != nil {
		// テーブルが存在しない場合はエラーになるので無視する
		fmt.Println(err)
	}

	err = db.CreateTable(tableName, model.Diary{}).Run()
	if err != nil {
		t.Fatal(err)
	}

	diary := model.NewDiary("https://www.hinatazaka46.com/s/official/diary/detail/34467", "タイトル", "小坂菜緒", time.Now(), 1)
	err = repo.PostDiary(diary)
	if err != nil {
		t.Fatal(err)
	}

	diary, err = repo.GetDiary("小坂菜緒", 1)
	if err != nil {
		t.Fatal(err)
	}

	if diary.Title != "タイトル" {
		t.Errorf("diary.Title = %s, want %s", diary.Title, "タイトル")
	}
}

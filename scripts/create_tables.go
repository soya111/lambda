package main

import (
	"fmt"
	"notify/pkg/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

const (
	AWS_REGION      = "ap-northeast-1"
	DYNAMO_ENDPOINT = "http://dynamodb-local:8000"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Endpoint:    aws.String(DYNAMO_ENDPOINT),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	}))
	db := dynamo.New(sess)

	tableName := "hinatazaka_blog"
	// テーブルが存在するかチェック
	if _, err := db.Table(tableName).Describe().Run(); err == nil {
		// テーブルが存在する場合、削除
		fmt.Printf("Table %s exists, deleting...\n", tableName)
		if err := db.Table(tableName).DeleteTable().Run(); err != nil {
			panic(err)
		}
		fmt.Printf("Table %s deleted successfully.\n", tableName)
	}

	table := db.CreateTable(tableName, model.Diary{})
	if err := table.Run(); err != nil {
		panic(err)
	}
	fmt.Printf("Table %s created successfully.\n", tableName)

	// テーブル一覧を表示
	tables, err := db.ListTables().All()
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables:")
	for _, table := range tables {
		fmt.Println(table)
	}
}

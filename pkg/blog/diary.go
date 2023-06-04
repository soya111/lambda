package blog

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	// メンバーに紐づいた番号
	artistNumbers = []int{2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
)

type Diary struct {
	Url        string `dynamodbav:"url" json:"url"`
	Title      string `dynamodbav:"title" json:"title"`
	MemberName string `dynamodbav:"member_name" json:"member_name"`
	Date       string `dynamodbav:"date" json:"date"`
	Id         int    `dynamodbav:"diary_id" json:"diary_id"`
	Images     []string
}

func NewDiary(url string, title string, memberName string, date time.Time, id int) *Diary {
	return &Diary{url, title, memberName, date.Format("2006.1.2 15:04 (MST)"), id, []string{}}
}

func reverse(a []*Diary) []*Diary {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - i - 1
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// DynamoからGET
func GetDiary(tableName string, memberName string, diaryId int) (*Diary, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)

	// 検索条件を用意
	getParam := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"member_name": {
				S: aws.String(memberName),
			},
			"diary_id": {
				N: aws.String(strconv.Itoa(diaryId)),
			},
		},
	}

	// 検索
	result, err := db.GetItem(getParam)
	if err != nil {
		return nil, err
	}

	diary := Diary{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &diary)
	if err != nil {
		return nil, err
	}

	return &diary, nil
}

// DynamoにPOST
func PostDiary(tableName string, diary *Diary) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	db := dynamodb.New(sess)

	// attribute value作成
	inputAV, err := dynamodbattribute.MarshalMap(diary)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      inputAV,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

type Scraper interface {
	GetLatestDiaries() ([]*Diary, error)
	PostDiaries([]*Diary) error
	GetImages(*Diary) []string
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

var (
	// メンバーに紐づいた番号
	artistNumbers = []int{2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	rootURL       = "https://www.hinatazaka46.com"
)

type Diary struct {
	Url        string `dynamodbav:"url" json:"url"`
	Title      string `dynamodbav:"title" json:"title"`
	MemberName string `dynamodbav:"member_name" json:"member_name"`
	Date       string `dynamodbav:"date" json:"date"`
	Id         int    `dynamodbav:"diary_id" json:"diary_id"`
}

func NewDiary(url string, title string, memberName string, date time.Time, id int) *Diary {
	return &Diary{url, title, memberName, date.Format("2006.1.2 15:04 (MST)"), id}
}

// line送信
func pushImageMessages(to []string, messages ...*linebot.ImageMessage) {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range messages {
		for _, to := range to {
			if _, err := bot.PushMessage(to, message).Do(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func pushTextMessages(to []string, messages ...*linebot.TextMessage) {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range messages {
		for _, to := range to {
			if _, err := bot.PushMessage(to, message).Do(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// urlをImageMessageに変換
func makeImageMessages(urls []string) []*linebot.ImageMessage {
	var messages []*linebot.ImageMessage
	for _, url := range urls {
		messages = append(messages, linebot.NewImageMessage(url, url))
	}
	return messages
}

// blog中の全画像を取得
func getImagesFromDiary(diary *Diary) []string {
	document, err := getDocumentFromURL(diary.Url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	article := document.Find(".c-blog-article__text")
	img := article.Find("img")
	srcs := []string{}
	img.Each(func(i int, s *goquery.Selection) {
		for _, attr := range s.Nodes[0].Attr {
			if attr.Key == "src" {
				srcs = append(srcs, attr.Val)
			}
		}
	})

	return srcs
}

// Document structを返す
func getDocumentFromURL(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	defer res.Body.Close()

	// Body内を読み取り
	buffer, _ := ioutil.ReadAll(res.Body)

	// 文字コードを判定
	detector := chardet.NewTextDetector()
	detectResult, _ := detector.DetectBest(buffer)

	// 文字コードの変換
	bufferReader := bytes.NewReader(buffer)
	reader, _ := charset.NewReaderLabel(detectResult.Charset, bufferReader)

	// HTMLをパース
	document, _ := goquery.NewDocumentFromReader(reader)

	return document, nil
}

func getIdFromHref(href string) int {
	id, err := strconv.Atoi(strings.Split(strings.Split(href, "/")[5], "?")[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return id
}

// DynamoからGET
func getDiaryFromDynamoDB(memberName string, diaryId int) (*Diary, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)

	// 検索条件を用意
	getParam := &dynamodb.GetItemInput{
		TableName: aws.String("hinatazaka_blog"),
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
func postDiaryFromDynamoDB(diary *Diary) error {
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
		TableName: aws.String("hinatazaka_blog"),
		Item:      inputAV,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// 最新のblogを調べる
func getLatestDiaries() []*Diary {
	url := "https://www.hinatazaka46.com/s/official/diary/member/list?ima=0000"
	document, err := getDocumentFromURL(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	res := []*Diary{}
	articles := document.Find(".p-blog-article")
	// blogひとつづつ更新確認
	articles.Each(func(i int, s *goquery.Selection) {
		updateAt, err := time.Parse("2006.1.2 15:04 (MST)", strings.TrimSpace(s.Find(".c-blog-article__date").Text())+" (JST)")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		href := s.Find(".c-button-blog-detail").Nodes[0].Attr[1].Val
		title := strings.TrimSpace(s.Find(".c-blog-article__title").Text())
		name := strings.TrimSpace(s.Find(".c-blog-article__name").Text())
		diaryId := getIdFromHref(href)

		diary, err := getDiaryFromDynamoDB(name, diaryId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Dynamoにデータがない場合
		if diary.Id == 0 {
			fmt.Printf("%s %s %s\n%s\n", updateAt.Format("2006.1.2 15:04 (MST)"), title, name, rootURL+href)
			newDiary := NewDiary(rootURL+href, title, name, updateAt, diaryId)
			if err := postDiaryFromDynamoDB(newDiary); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			res = append(res, newDiary)
		}
	})

	return res
}

func reverse(a []*Diary) []*Diary {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - i - 1
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func excuteFunction() {
	to := []string{os.Getenv("ME"), os.Getenv("KNPT")}
	latestDiaries := getLatestDiaries()
	for _, diary := range reverse(latestDiaries) {
		images := getImagesFromDiary(diary)
		text := fmt.Sprintf("%s %s %s\n%s", diary.Date, diary.Title, diary.MemberName, diary.Url)
		pushTextMessages(to, []*linebot.TextMessage{linebot.NewTextMessage(text)}...)
		messages := makeImageMessages(images)
		pushImageMessages(to, messages...)
	}
}

func init() {
	// set timezone
	time.Local = time.FixedZone("JST", 9*60*60)
}

func main() {
	lambda.Start(excuteFunction)
}

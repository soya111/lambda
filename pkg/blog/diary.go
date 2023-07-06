package blog

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	// メンバーに紐づいた番号
	artistNumbers = []int{2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
)

type Diary struct {
	Url        string   `dynamo:"url" json:"url"`
	Title      string   `dynamo:"title" json:"title"`
	MemberName string   `dynamo:"member_name,hash" json:"member_name"`
	Date       string   `dynamo:"date" json:"date"`
	Id         int      `dynamo:"diary_id,range" json:"diary_id"`
	Images     []string `dynamo:"images,omitempty" json:"images"`
}

func NewDiary(url string, title string, memberName string, date time.Time, id int) *Diary {
	return &Diary{url, title, memberName, date.Format("2006.1.2 15:04 (MST)"), id, []string{}}
}

func reverse[T any](a []T) []T {
	for i := 0; i < len(a)/2; i++ {
		j := len(a) - i - 1
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// DiaryRepository provides an interface for database operations on Diaries
type DiaryRepository interface {
	GetDiary(memberName string, diaryId int) (*Diary, error)
	PostDiary(diary *Diary) error
}

type DynamoDiaryRepository struct {
	db    *dynamo.DB
	table dynamo.Table
}

func NewDynamoDiaryRepository(sess *session.Session, tableName string) *DynamoDiaryRepository {
	db := dynamo.New(sess)
	table := db.Table(tableName)
	return &DynamoDiaryRepository{
		db:    db,
		table: table,
	}
}

var ErrDiaryNotFound = errors.New("diary not found")

func (r *DynamoDiaryRepository) GetDiary(memberName string, diaryId int) (*Diary, error) {
	diary := new(Diary)
	err := r.table.Get("member_name", memberName).Range("diary_id", dynamo.Equal, diaryId).One(diary)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, ErrDiaryNotFound
		}
		return nil, err
	}
	return diary, nil
}

func (r *DynamoDiaryRepository) PostDiary(diary *Diary) error {
	return r.table.Put(diary).Run()
}

type Scraper interface {
	GetLatestDiaries() ([]*Diary, error)
	PostDiaries([]*Diary) error
	GetImages(*Diary) []string
}

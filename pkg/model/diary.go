package model

import (
	"errors"
	"time"
)

var (
	// メンバーに紐づいた番号
	ArtistNumbers = []int{2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
)

// Diary represents a diary
type Diary struct {
	Url        string   `dynamo:"url" json:"url"`
	Title      string   `dynamo:"title" json:"title"`
	MemberName string   `dynamo:"member_name,hash" json:"member_name"`
	Date       string   `dynamo:"date" json:"date"`
	Id         int      `dynamo:"diary_id,range" json:"diary_id"`
	Images     []string `dynamo:"images,omitempty" json:"images"`
}

// NewDiary creates a new Diary
func NewDiary(url string, title string, memberName string, date time.Time, id int) *Diary {
	return &Diary{url, title, memberName, date.Format("2006.1.2 15:04 (MST)"), id, []string{}}
}

// DiaryRepository provides an interface for database operations on Diaries
type DiaryRepository interface {
	GetDiary(memberName string, diaryId int) (*Diary, error)
	PostDiary(diary *Diary) error
}

// ErrDiaryNotFound is returned when the requested diary is not found
var ErrDiaryNotFound = errors.New("diary not found")

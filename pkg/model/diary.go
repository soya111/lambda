package model

import (
	"errors"
	"time"
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

package main

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type ScraperMock struct {
}

func (s *ScraperMock) getLatestDiaries() []*Diary {
	return []*Diary{
		NewDiary("https://www.hinatazaka46.com/s/official/diary/detail/20317", "ニャー0( =^ ・_・^)= 〇", "加藤 史帆", time.Now(), 20317),
	}
}

func TestExcute(t *testing.T) {
	godotenv.Load(".env")
	excute(&ScraperMock{})
}

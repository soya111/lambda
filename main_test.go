package main

import (
	"main/blog"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type ExecutorMock struct {
}

func (s *ExecutorMock) GetLatestDiaries() []*blog.Diary {
	return []*blog.Diary{
		blog.NewDiary("https://www.hinatazaka46.com/s/official/diary/detail/20317", "ニャー0( =^ ・_・^)= 〇", "加藤 史帆", time.Now(), 20317),
	}
}

func TestExcute(t *testing.T) {
	godotenv.Load(".env")
	excute(&ExecutorMock{}, []string{os.Getenv("ME")})
}

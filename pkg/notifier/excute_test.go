package notifier

import (
	"fmt"
	"testing"
	"time"

	"notify/pkg/blog"

	"github.com/joho/godotenv"
)

type ScraperMock struct{}

func (*ScraperMock) GetLatestDiaries() ([]*blog.Diary, error) {
	return []*blog.Diary{
		blog.NewDiary("https://www.hinatazaka46.com/s/official/diary/detail/20317", "ニャー0( =^ ・_・^)= 〇", "加藤 史帆", time.Now(), 20317),
	}, nil
}

func (*ScraperMock) PostDiaries(diaries []*blog.Diary) error {
	// モックなので何もしない
	return nil
}

func (*ScraperMock) GetImages(diary *blog.Diary) []string {
	var s = &blog.HinatazakaScraper{}
	return s.GetImages(diary)
}

type BotMock struct{}

func (*BotMock) PushTextMessages(to []string, messages ...string) {
	fmt.Println("PushTextMessages")
	fmt.Println(to, messages)
}

func (*BotMock) PushFlexImagesMessage(to []string, urls []string) {
	fmt.Println("PushFlexImagesMessage")
	fmt.Println(to, urls)
}

func (*BotMock) ReplyTextMessages(token string, message string) error {
	return nil
}

type DBMock struct{}

func (*DBMock) GetDestination(memberName string) ([]string, error) {
	return []string{"kosakana"}, nil
}

func TestExcute(t *testing.T) {
	godotenv.Load("../.env")
	Excute(&ScraperMock{}, &BotMock{}, &DBMock{})
}

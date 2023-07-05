package notifier

import (
	"context"
	"fmt"
	"testing"
	"time"

	"notify/pkg/blog"
	"notify/pkg/database"

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

func (*BotMock) PushTextMessages(ctx context.Context, to []string, messages ...string) error {
	fmt.Println("PushTextMessages")
	fmt.Println(to, messages)
	return nil
}

func (*BotMock) PushFlexImagesMessage(ctx context.Context, to []string, urls []string) error {
	fmt.Println("PushFlexImagesMessage")
	fmt.Println(to, urls)
	return nil
}

func (*BotMock) ReplyTextMessages(ctx context.Context, token string, message string) error {
	return nil
}

type MockSubscriberRepository struct{}

func (*MockSubscriberRepository) GetAllByMemberName(memberName string) ([]string, error) {
	return []string{"こさかな"}, nil
}

func (*MockSubscriberRepository) Subscribe(subscriber database.Subscriber) error {
	return nil
}

func (*MockSubscriberRepository) Unsubscribe(memberName, userId string) error {
	return nil
}

func (*MockSubscriberRepository) GetAllById(id string) ([]database.Subscriber, error) {
	return []database.Subscriber{}, nil
}

func TestExcute(t *testing.T) {
	godotenv.Load("../.env")
	Excute(context.Background(), &ScraperMock{}, &BotMock{}, &MockSubscriberRepository{})
}

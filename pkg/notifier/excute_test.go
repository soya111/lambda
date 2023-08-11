package notifier

import (
	"context"
	"os"
	"testing"
	"time"

	"notify/pkg/blog"
	"notify/pkg/infrastructure/line"
	"notify/pkg/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

type ScraperMock struct{}

func (*ScraperMock) ScrapeLatestDiaries() ([]*blog.ScrapedDiary, error) {
	// return []*model.Diary{
	// 	model.NewDiary("https://www.hinatazaka46.com/s/official/diary/detail/20317", "ニャー0( =^ ・_・^)= 〇", "加藤 史帆", time.Now(), 20317),
	// }, nil
	return []*blog.ScrapedDiary{
		blog.NewScrapedDiary("https://www.hinatazaka46.com/s/official/diary/detail/20317", "ニャー0( =^ ・_・^)= 〇", "加藤 史帆", time.Now(), 20317, []string{}, ""),
	}, nil
}

func (*ScraperMock) PostDiaries(diaries []*blog.ScrapedDiary) error {
	// モックなので何もしない
	return nil
}

func (*ScraperMock) GetImages(document *goquery.Document) []string {
	var s = &blog.HinatazakaScraper{}
	return s.GetImages(document)
}

func (*ScraperMock) GetMemberIcon(document *goquery.Document) string {
	return "https://cdn.hinatazaka46.com/images/14/0a0/472f1b55902a03c7b685fd958e085/400_320_102400.jpg"
}

type MockSubscriberRepository struct {
	to string
}

func NewMockSubscriberRepository(to string) *MockSubscriberRepository {
	return &MockSubscriberRepository{to}
}

func (s *MockSubscriberRepository) GetAllByMemberName(memberName string) ([]string, error) {
	return []string{s.to}, nil
}

func (*MockSubscriberRepository) Subscribe(subscriber model.Subscriber) error {
	return nil
}

func (*MockSubscriberRepository) Unsubscribe(memberName, userId string) error {
	return nil
}

func (*MockSubscriberRepository) GetAllById(id string) ([]model.Subscriber, error) {
	return []model.Subscriber{}, nil
}

func TestExecute(t *testing.T) {
	t.Skip("skipping this test for now")
	_ = godotenv.Load("../.env")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")
	me := os.Getenv("ME")
	bot, err := line.NewLinebot(channelSecret, channelToken)
	if err != nil {
		t.Fatal(err)
	}
	notifier := NewNotifier(&ScraperMock{}, bot, NewMockSubscriberRepository(me), nil)
	err = notifier.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

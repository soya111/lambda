package notifier

import (
	"os"
	"testing"
	"time"

	"notify/internal/pkg/blog"

	"github.com/joho/godotenv"
)

type ScraperMock struct {
}

func (*ScraperMock) GetAndPostLatestDiaries() []*blog.Diary {
	return []*blog.Diary{
		blog.NewDiary("https://www.hinatazaka46.com/s/official/diary/detail/20317", "ニャー0( =^ ・_・^)= 〇", "加藤 史帆", time.Now(), 20317),
	}
}

func (*ScraperMock) GetImages(diary *blog.Diary) []string {
	var s = &blog.HinatazakaScraper{}
	return s.GetImages(diary)
}

func TestExcute(t *testing.T) {
	godotenv.Load("../.env")
	Excute(&ScraperMock{}, []string{os.Getenv("ME")})
}

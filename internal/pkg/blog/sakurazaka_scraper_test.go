package blog

import (
	"fmt"
	"testing"
	"time"
)

func TestGetLatestDiaries(t *testing.T) {
	s := &SakurazakaScraper{}
	d := s.getLatestDiaries()
	fmt.Printf("%#v\n", d[0])
}

func TestGetImages(t *testing.T) {
	s := &SakurazakaScraper{}
	d := NewDiary("https://sakurazaka46.com/s/s46/diary/detail/42564?ima=2759&cd=blog", "おしらせ〜", "原田 葵", time.Now(), 42564)
	imgs := s.GetImages(d)
	fmt.Printf("%#v\n", imgs)
}

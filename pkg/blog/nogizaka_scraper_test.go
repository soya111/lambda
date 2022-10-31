package blog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIdFromURL(t *testing.T) {
	s := &NogizakaScraper{}
	for _, tc := range []struct {
		title string
		url   string
		want  int
		isErr bool
	}{
		{
			title: "No error",
			url:   "https://blog.nogizaka46.com/tamami.sakaguchi/2022/01/065171.php",
			want:  65171,
			isErr: false,
		},
		{
			title: "Error",
			url:   "https://blog.nogizaka46.com/tamami.sakaguchi/2022/01/a.php",
			want:  0,
			isErr: true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			id, err := s.getIdFromURL(tc.url)
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, id, tc.want)
			}
		})
	}
}

func TestGetLatestDiaries(t *testing.T) {
	s := &NogizakaScraper{}
	d := s.getLatestDiaries()
	fmt.Printf("%#v\n", d[0])
}

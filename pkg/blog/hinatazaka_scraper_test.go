package blog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestDiaryByMember(t *testing.T) {
	s := NewHinatazakaScraper()
	tests := []struct {
		name       string
		memberName string
		expectErr  bool
	}{
		{"加藤史帆", "加藤史帆", false},
		{"正源司陽子", "正源司陽子", false},
		{"ポカ", "ポカ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diary, err := s.GetLatestDiaryByMember(tt.memberName)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, diary)
			}
			if diary != nil {
				fmt.Printf("%+v\n", diary)
			}
		})
	}
}

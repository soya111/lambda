package profile

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetProfileSelection(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectederror error
	}{
		{"ExistentMember", "潮紗理菜", nil},
		{"NonExistentMember", "白石麻衣", ErrNonExistentMember},
		{"ポカ", "ポカ", ErrNoUrl},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, actual := GetProfileSelection(tt.input)
			assert.Equal(t, tt.expectederror, actual)
		})
	}
}

func TestScrapeProfile(t *testing.T) {
	var ushiosarinaProfile = Profile{
		"1997年12月26日",
		calcAge(time.Date(1997, 12, 26, 0, 0, 0, 0, time.Local), time.Now()),
		"やぎ座",
		"157.5cm",
		"神奈川県",
		"O型",
		"https://cdn.hinatazaka46.com/images/14/9d4/dc3eef1e11944f0ee69459463a4cb/1000_1000_102400.jpg",
	}

	tests := []struct {
		name     string
		input    string
		expected Profile
	}{
		{"Nornal", "潮紗理菜", ushiosarinaProfile},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			selection, _ := GetProfileSelection(tt.input)
			actual := ScrapeProfile(selection)
			assert.Equal(t, tt.expected, *actual)
		})
	}
}

func Test_normalizeDate(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expecteddate time.Time
		expectederr  bool
	}{
		{"YYYY年MM月DD日", "2000年12月12日", time.Date(2000, 12, 12, 0, 0, 0, 0, time.UTC), false},
		{"YYYY年M月DD日", "2000年2月12日", time.Date(2000, 2, 12, 0, 0, 0, 0, time.UTC), false},
		{"YYYY年MM月D日", "2000年12月2日", time.Date(2000, 12, 2, 0, 0, 0, 0, time.UTC), false},
		{"YYYY年M月D日", "2000年2月2日", time.Date(2000, 2, 2, 0, 0, 0, 0, time.UTC), false},
		{"YYYY/M/D日", "2000/2/2", time.Time{}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := normalizeDate(tt.input)
			if tt.expectederr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expecteddate, actual)
			}
		})
	}
}

func Test_calcAge(t *testing.T) {
	tests := []struct {
		name          string
		inputbirthday time.Time
		inputnow      time.Time
		expected      string
	}{
		{"BeforeBirthday", time.Date(2000, 12, 15, 0, 0, 0, 0, time.Local), time.Date(2020, 8, 15, 0, 0, 0, 0, time.Local), "19"},
		{"AfterBirthday", time.Date(2000, 6, 15, 0, 0, 0, 0, time.Local), time.Date(2020, 8, 15, 0, 0, 0, 0, time.Local), "20"},
		{"TodayIsBirthday", time.Date(2000, 8, 15, 0, 0, 0, 0, time.Local), time.Date(2020, 8, 15, 0, 0, 0, 0, time.Local), "20"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := calcAge(tt.inputbirthday, tt.inputnow)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

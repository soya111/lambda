package profile

import (
	"notify/pkg/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getProfileSelection(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"ExistentMember", "潮紗理菜", nil},
		{"NonExistentMember", "白石麻衣", model.ErrNonExistentMember},
		{"ポカ", "ポカ", ErrNoUrl},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			selection, err := getProfileSelection(tt.input)
			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NotNil(t, selection)
			}
		})
	}
}

func TestScrapeProfile(t *testing.T) {
	var ushiosarinaProfile = Profile{
		"潮紗理菜",
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
			profile, err := ScrapeProfile(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *profile)
		})
	}
}

func Test_normalizeDate(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedDate time.Time
		expectedErr  bool
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
			date, err := normalizeDate(tt.input)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDate, date)
			}
		})
	}
}

func Test_calcAge(t *testing.T) {
	tests := []struct {
		name          string
		inputBirthday time.Time
		inputNow      time.Time
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
			actual := calcAge(tt.inputBirthday, tt.inputNow)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCreateProfileMessage(t *testing.T) {
	var ushiosarinaProfile = &Profile{
		"潮紗理菜",
		"1997年12月26日",
		calcAge(time.Date(1997, 12, 26, 0, 0, 0, 0, time.Local), time.Now()),
		"やぎ座",
		"157.5cm",
		"神奈川県",
		"O型",
		"https://cdn.hinatazaka46.com/images/14/9d4/dc3eef1e11944f0ee69459463a4cb/1000_1000_102400.jpg",
	}

	tests := []struct {
		name         string
		inputProfile *Profile
		expected     string
	}{
		{"Nornal", ushiosarinaProfile, "潮紗理菜\n生年月日:1997年12月26日\n年齢:25歳\n星座:やぎ座\n身長:157.5cm\n出身地:神奈川県\n血液型:O型"},
		{"ポカ", PokaProfile, "ポカ\n生年月日:2019年12月25日\n年齢:3歳\n星座:やぎ座\n身長:???\n出身地:???\n血液型:???"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := CreateProfileMessage(tt.inputProfile)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

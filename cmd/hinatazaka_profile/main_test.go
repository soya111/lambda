package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_inputName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Nomal", "潮紗理菜", "潮紗理菜"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := flag.CommandLine.Set("name", tt.input)
			if err != nil {
				fmt.Printf("Error setting command-line arguments: %v\n", err)
			}
			t.Parallel()
			actual := inputName()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_getProfileSelection(t *testing.T) {
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
			_, actual := getProfileSelection(tt.input)
			assert.Equal(t, tt.expectederror, actual)
		})
	}
}

func Test_scrapeProfile(t *testing.T) {
	var ushiosarinaProfile = profile{
		"1997年12月26日",
		calcAge(time.Date(1997, 12, 26, 0, 0, 0, 0, time.Local), time.Now()),
		"やぎ座",
		"157.5cm",
		"神奈川県",
		"O型",
	}

	tests := []struct {
		name     string
		input    string
		expected profile
	}{
		{"Nornal", "潮紗理菜", ushiosarinaProfile},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			selection, _ := getProfileSelection(tt.input)
			actual := scrapeProfile(selection)
			assert.Equal(t, tt.expected, *actual)
		})
	}
}

func Test_normalizeDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{"YYYY年MM月DD日", "2000年12月12日", time.Date(2000, 12, 12, 0, 0, 0, 0, time.UTC)},
		{"YYYY年M月DD日", "2000年2月12日", time.Date(2000, 2, 12, 0, 0, 0, 0, time.UTC)},
		{"YYYY年MM月D日", "2000年12月2日", time.Date(2000, 12, 2, 0, 0, 0, 0, time.UTC)},
		{"YYYY年M月D日", "2000年2月2日", time.Date(2000, 2, 2, 0, 0, 0, 0, time.UTC)},
		{"YYYY/M/D日", "2000/2/2", time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, _ := normalizeDate(tt.input)
			assert.Equal(t, tt.expected, actual)
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

func Test_isTodayBirthday(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected bool
	}{
		{"Birthday", time.Now(), true},
		{"NotBirthday", time.Now().AddDate(0, 0, 1), false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := isTodayBirthday(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_outputProfile(t *testing.T) {
	var ushiosarinaProfile = profile{
		"1997年12月26日",
		calcAge(time.Date(1997, 12, 26, 0, 0, 0, 0, time.Local), time.Now()),
		"やぎ座",
		"157.5cm",
		"神奈川県",
		"O型",
	}

	tests := []struct {
		name         string
		inputname    string
		inputprofile profile
		expected     string
	}{
		{"Nornal", "潮紗理菜", ushiosarinaProfile, "潮紗理菜\n生年月日:1997年12月26日, 年齢:25歳, 星座:やぎ座, 身長:157.5cm, 出身地:神奈川県, 血液型:O型\n"},
		{"ポカ", "ポカ", pokaProfile, "ポカ\n生年月日:2019年12月25日, 年齢:3歳, 星座:やぎ座, 身長:???, 出身地:???, 血液型:???\n"},
	}

	for _, tt := range tests {
		// 標準出力をキャプチャ
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// テスト対象の関数呼び出し
		outputProfile(tt.inputname, tt.inputprofile)

		// 標準出力を元に戻す
		w.Close()
		os.Stdout = old

		// キャプチャした出力を読み取る
		var capturedOutput string
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r)
		capturedOutput = buf.String()

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := capturedOutput
			assert.Equal(t, tt.expected, actual)
		})
	}
}

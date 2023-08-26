package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

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
			_, actual := getProfileSelection(tt.input)
			assert.Equal(t, tt.expectederror, actual)
		})
	}
}

func Test_scrapeProfile(t *testing.T) {
	var usiosarinaProfile = profile{
		"1997年12月26日",
		"",
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
		{"Nornal", "潮紗理菜", usiosarinaProfile},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			selection, _ := getProfileSelection(tt.input)
			actual := scrapeProfile(selection)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_outputProfile(t *testing.T) {
	var usiosarinaProfile = profile{
		"1997年12月26日",
		"",
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
		{"Nornal", "潮紗理菜", usiosarinaProfile, "潮紗理菜\n生年月日:1997年12月26日, 年齢:25歳, 星座:やぎ座, 身長:157.5cm, 出身地:神奈川県, 血液型:O型"},
		{"ポカ", "ポカ", pokaProfile, "ポカ\n生年月日:2019年12月25日, 年齢:3歳, 星座:やぎ座, 身長:???, 出身地:???, 血液型:???"},
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
			actual := capturedOutput
			assert.Equal(t, tt.expected, actual)
		})
	}
}

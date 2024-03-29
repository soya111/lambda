package model_test

import (
	"testing"
	"zephyr/pkg/model"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"No Spaces", "name", "name"},
		{"Leading Spaces", "  name", "name"},
		{"Trailing Spaces", "name  ", "name"},
		{"Middle Spaces", "n am e", "name"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := model.NormalizeName(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestTranslateNicknametoMember(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"nomal", "なっちょ", "潮紗理菜"},
		{"NoExist", "たなかくん", "たなかくん"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := model.TranslateNicknametoMember(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestIsMember(t *testing.T) {
	assert.True(t, model.IsMember("潮紗理菜"))
	assert.False(t, model.IsMember("非メンバー"))
}

func TestIsGrad(t *testing.T) {
	assert.True(t, model.IsGrad("柿崎芽実"))
	assert.False(t, model.IsGrad("非卒業メンバー"))
}

func TestGetMemberId(t *testing.T) {
	id, err := model.GetMemberId("潮紗理菜")
	assert.NoError(t, err)
	assert.Equal(t, "2", id)

	_, err = model.GetMemberId("非メンバー")
	assert.Error(t, err)
}

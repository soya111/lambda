package line_test

import (
	"testing"
	"zephyr/pkg/infrastructure/line"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
)

func TestCreateTextMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []string
		want     []linebot.SendingMessage
	}{
		{
			name:     "No Messages",
			messages: []string{},
			want:     []linebot.SendingMessage{},
		},
		{
			name:     "Single Message",
			messages: []string{"Hello, World!"},
			want:     []linebot.SendingMessage{linebot.NewTextMessage("Hello, World!")},
		},
		{
			name:     "Multiple Messages",
			messages: []string{"Hello, World!", "How are you?"},
			want:     []linebot.SendingMessage{linebot.NewTextMessage("Hello, World!"), linebot.NewTextMessage("How are you?")},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := line.CreateTextMessages(tt.messages...)
			assert.Equal(t, tt.want, got)
		})
	}
}

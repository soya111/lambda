package line_test

import (
	"testing"
	"zephyr/pkg/infrastructure/line"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
)

func TestExtractEventSourceIdentifier(t *testing.T) {
	tests := []struct {
		name   string
		source linebot.EventSource
		want   string
	}{
		{
			name:   "User Source",
			source: linebot.EventSource{Type: linebot.EventSourceTypeUser, UserID: "userID123"},
			want:   "userID123",
		},
		{
			name:   "Group Source",
			source: linebot.EventSource{Type: linebot.EventSourceTypeGroup, GroupID: "groupID123"},
			want:   "groupID123",
		},
		{
			name:   "Room Source",
			source: linebot.EventSource{Type: linebot.EventSourceTypeRoom, RoomID: "roomID123"},
			want:   "roomID123",
		},
		{
			name:   "Unknown Source",
			source: linebot.EventSource{Type: "unknown"},
			want:   "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			event := &linebot.Event{Source: &tt.source}
			got := line.ExtractEventSourceIdentifier(event)
			assert.Equal(t, tt.want, got)
		})
	}
}

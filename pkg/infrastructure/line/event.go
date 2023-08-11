package line

import "github.com/line/line-bot-sdk-go/v7/linebot"

// ExtractEventSourceIdentifier returns the event source identifier.
func ExtractEventSourceIdentifier(event *linebot.Event) string {
	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		return event.Source.UserID
	case linebot.EventSourceTypeGroup:
		return event.Source.GroupID
	case linebot.EventSourceTypeRoom:
		return event.Source.RoomID
	}

	return ""
}

package line

import (
	"context"
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// ReplyTextMessages reply text messages to LINE
func (b *Linebot) ReplyTextMessages(ctx context.Context, token string, message string) error {
	if _, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(message)).WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}

// ReplyTextandImageMessages reply text and image messages to LINE
func (b *Linebot) ReplyTextandImageMessages(ctx context.Context, token string, message, imgurl string) error {
	if _, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(message), linebot.NewImageMessage(imgurl, imgurl)).WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}

// ReplyMessage reply messages to LINE
func (b *Linebot) ReplyMessage(ctx context.Context, token string, messages ...linebot.SendingMessage) error {
	if _, err := b.Client.ReplyMessage(token, messages...).WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}

// ReplyWithError reply message to LINE and return error
func (b *Linebot) ReplyWithError(ctx context.Context, token, replyMessage string, err error) error {
	if err := b.ReplyTextMessages(ctx, token, replyMessage); err != nil {
		return fmt.Errorf("replyWithError: %w", err)
	}
	return err
}

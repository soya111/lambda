package line

import (
	"context"
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (b *Linebot) ReplyTextMessages(ctx context.Context, token string, message string) error {
	if _, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(message)).WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}

func (b *Linebot) ReplyMessage(ctx context.Context, token string, messages ...linebot.SendingMessage) error {
	if _, err := b.Client.ReplyMessage(token, messages...).WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}

func (b *Linebot) ReplyWithError(ctx context.Context, token, replyMessage string, err error) error {
	if err := b.ReplyTextMessages(ctx, token, replyMessage); err != nil {
		return fmt.Errorf("replyWithError: %w", err)
	}
	return err
}

package infrastructure

import "context"

type Client interface {
	PushTextMessages(ctx context.Context, to []string, messages ...string) error
	PushFlexImagesMessage(ctx context.Context, to []string, urls []string) error
	ReplyTextMessages(ctx context.Context, token string, message string) error
}

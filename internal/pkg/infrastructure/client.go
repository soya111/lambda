package infrastructure

type Client interface {
	PushTextMessages(to []string, messages ...string)
	PushFlexImagesMessage(to []string, urls []string)
	ReplyTextMessages(token string, message string) error
}

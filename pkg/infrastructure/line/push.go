package line

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// PushMessages push messages to LINE
func (b *Linebot) PushMessages(ctx context.Context, to []string, messages ...linebot.SendingMessage) error {
	var result *multierror.Error

	for _, to := range to {
		response, err := b.Client.PushMessage(to, messages...).WithContext(ctx).Do()
		if err != nil {
			result = multierror.Append(result, err)
			fmt.Printf("Push Messages to %s, error: %v\n", to, err)
		} else {
			fmt.Printf("Push Messages to %s, success, response: %v\n", to, response)
		}
	}
	return result.ErrorOrNil()
}

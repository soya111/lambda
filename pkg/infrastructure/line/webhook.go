package line

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// ParseRequest without signature validation
func ParseRequestWithoutSignatureValidation(r *http.Request) ([]*linebot.Event, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err = json.Unmarshal(body, request); err != nil {
		return nil, err
	}
	return request.Events, nil
}

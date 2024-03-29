package line

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// RequestParser is an interface for parsing request
type RequestParser interface {
	ParseRequest(req *http.Request) ([]*linebot.Event, error)
}

// LocalParser is a parser for local development
type LocalParser struct{}

// ParseRequest without signature validation
// This function is used for testing
// Dangerous: This function does not validate the signature
// If you want to validate the signature, use ParseRequest instead
// When you get dangerous message, you should exit the program
// TODO: remove this function
func (p *LocalParser) ParseRequest(r *http.Request) ([]*linebot.Event, error) {
	fmt.Println("Dangerous: This function does not validate the signature")

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

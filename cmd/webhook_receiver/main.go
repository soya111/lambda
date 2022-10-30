package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notify/internal/app/webhook"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func init() {
	// localで実行するとき用
	err := godotenv.Load(".env")
	bot, err = linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path
	body := request.Body
	method := request.HTTPMethod

	lambdaCtx, _ := lambdacontext.FromContext(ctx)
	requestId := lambdaCtx.AwsRequestID

	switch path {
	case "/webhook":
		// LINEのsdkがHTTPを前提にParseしているのでHttpRequestに戻す
		r := &core.RequestAccessor{}
		httpRequest, err := r.EventToRequest(request)
		if err != nil {
			return newResponse(http.StatusInternalServerError), err
		}

		events, err := bot.ParseRequest(httpRequest)
		if err != nil {
			fmt.Printf("RequestId: %s, Method: %s, Path: %s, Body: %s\n", requestId, method, path, body)
			if err == linebot.ErrInvalidSignature {
				return newResponse(http.StatusBadRequest), err
			} else {
				return newResponse(http.StatusInternalServerError), err
			}
		}

		for _, event := range events {
			// 解析用ログ出力
			fmt.Println(marshal(event))

			var wg sync.WaitGroup
			wg.Add(3)
			go func() {
				defer wg.Done()
				res, _ := bot.GetProfile(event.Source.UserID).Do()
				fmt.Println(marshal(res))
			}()
			go func() {
				defer wg.Done()
				res2, _ := bot.GetGroupSummary(event.Source.GroupID).Do()
				fmt.Println(marshal(res2))
			}()
			go func() {
				defer wg.Done()
				res3, _ := bot.GetGroupMemberProfile(event.Source.GroupID, event.Source.UserID).Do()
				fmt.Println(marshal(res3))
			}()
			wg.Wait()

			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					webhook.HandleTextMessage(message.Text, event)
				}
			case linebot.EventTypeLeave:
				// webhook.HandleEventLeave(event)
			}

		}
		return newResponse(http.StatusOK), nil
	default:
		return newResponse(http.StatusBadRequest), nil
	}
}

func newAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	var headers = make(map[string]string)
	var mHeaders = make(map[string][]string)
	return events.APIGatewayProxyResponse{Headers: headers, MultiValueHeaders: mHeaders}
}

func newResponse(statusCode int) events.APIGatewayProxyResponse {
	res := newAPIGatewayProxyResponse()
	res.StatusCode = statusCode
	return res
}

func marshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("marshal: %s\n", err)
	}
	return string(b)
}

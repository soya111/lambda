package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notify/app/webhook"
	"notify/pkg/infrastructure/dynamodb"
	"notify/pkg/infrastructure/line"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/guregu/dynamo"
	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *line.Linebot
var db *dynamo.DB
var sess *session.Session

func init() {
	// localで実行するとき用
	err := godotenv.Load(".env")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")
	bot, err = line.NewLinebot(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	sess = session.Must(session.NewSession())
	db = dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
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
			return newResponse(http.StatusInternalServerError), fmt.Errorf("RequestId: %s, Method: %s, Path: %s, Body: %s, Error: %v", requestId, method, path, body, err)
		}

		events, err := bot.ParseRequest(httpRequest)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				return newResponse(http.StatusBadRequest), fmt.Errorf("invalid signature: %v", err)
			} else {
				return newResponse(http.StatusInternalServerError), fmt.Errorf("failed to parse request: %v", err)
			}
		}

		var result *multierror.Error
		var wg sync.WaitGroup

		// ここから正常系の処理をやる
		repo := dynamodb.NewDynamoSubscriberRepository(sess)
		handler := webhook.NewHandler(bot, db, repo)

		for _, event := range events {
			// 解析用ログ出力
			fmt.Println(marshal(event))

			// いつか見たい時が来た時ように出しとく
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

			err := handler.HandleEvent(ctx, event)
			if err != nil {
				result = multierror.Append(result, fmt.Errorf("RequestId: %s, Error: %v", requestId, err))
			}
		}
		wg.Wait()

		return newResponse(http.StatusOK), result.ErrorOrNil()
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

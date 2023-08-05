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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *line.Linebot
var sess *session.Session

func init() {
	var err error
	_ = godotenv.Load(".env")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")
	bot, err = line.NewLinebot(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	sess = session.Must(session.NewSession())
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path

	switch path {
	case "/webhook":
		return handleWebhook(ctx, request)
	default:
		return newResponse(http.StatusBadRequest), nil
	}
}

func handleWebhook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	method := request.HTTPMethod
	ip := request.RequestContext.Identity.SourceIP

	lambdaCtx, _ := lambdacontext.FromContext(ctx)
	requestId := lambdaCtx.AwsRequestID
	fmt.Printf("RequestId: %s, IP: %s, Method: %s, Path: /webhook, Body: %s\n", requestId, ip, method, body)

	r := &core.RequestAccessor{}
	httpRequest, err := r.EventToRequest(request)
	if err != nil {
		return handleError(err, "Failed to convert request", http.StatusInternalServerError)
	}

	events, err := bot.ParseRequest(httpRequest)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return handleError(err, "Invalid signature", http.StatusBadRequest)
		} else {
			return handleError(err, "Failed to parse request", http.StatusInternalServerError)
		}
	}

	var result *multierror.Error
	var wg sync.WaitGroup

	repo := dynamodb.NewSubscriberRepository(sess)
	handler := webhook.NewHandler(bot, repo)

	for _, event := range events {
		wg.Add(3)
		handleEventWithProfile(event, &wg)
		err := handler.HandleEvent(ctx, event)
		result = multierror.Append(result, err)
	}
	wg.Wait()

	return newResponse(http.StatusOK), result.ErrorOrNil()
}

func handleEventWithProfile(event *linebot.Event, wg *sync.WaitGroup) {
	// 解析用ログ出力
	fmt.Println(marshal(event))
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
}

func handleError(err error, msg string, statusCode int) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("%s: %v\n", msg, err)
	return newResponse(statusCode), fmt.Errorf("%s: %v", msg, err)
}

func newResponse(statusCode int) events.APIGatewayProxyResponse {
	res := events.APIGatewayProxyResponse{Headers: make(map[string]string), MultiValueHeaders: make(map[string][]string), Body: ""}
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

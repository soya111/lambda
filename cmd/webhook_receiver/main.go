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
	"notify/pkg/logging"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/guregu/dynamo"
	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

var (
	bot    *line.Linebot
	db     *dynamo.DB
	logger *zap.Logger
)

func init() {
	var err error
	_ = godotenv.Load(".env")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")
	bot, err = line.NewLinebot(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	sess := session.Must(session.NewSession())
	db = dynamo.New(sess)

	logger = logging.InitializeLogger()
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path

	ctx = logging.ContextWithLogger(ctx, logger)

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

	logger := logging.LoggerFromContext(ctx)

	lambdaCtx, _ := lambdacontext.FromContext(ctx)
	requestId := lambdaCtx.AwsRequestID
	logger.Info("request", zap.String("requestId", requestId), zap.String("ip", ip), zap.String("method", method), zap.String("path", "/webhook"), zap.String("body", body))

	r := &core.RequestAccessor{}
	httpRequest, err := r.EventToRequest(request)
	if err != nil {
		return handleError(ctx, err, "Failed to convert request", http.StatusInternalServerError)
	}

	events, err := bot.ParseRequest(httpRequest)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return handleError(ctx, err, "Invalid signature", http.StatusBadRequest)
		} else {
			return handleError(ctx, err, "Failed to parse request", http.StatusInternalServerError)
		}
	}

	var result *multierror.Error
	var wg sync.WaitGroup

	repo := dynamodb.NewSubscriberRepository(db)
	handler := webhook.NewHandler(bot, repo)

	for _, event := range events {
		wg.Add(3)
		handleEventWithProfile(event, &wg)
		logger.Info("Handling event", zap.String("eventType", string(event.Type)))
		if err := handler.HandleEvent(ctx, event); err != nil {
			logger.Error("Failed to handle event", zap.String("eventType", string(event.Type)), zap.Error(err))
		} else {
			logger.Info("Successfully handled event", zap.String("eventType", string(event.Type)))
		}
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

func handleError(ctx context.Context, err error, msg string, statusCode int) (events.APIGatewayProxyResponse, error) {
	logger := logging.LoggerFromContext(ctx)
	logger.Error(msg, zap.Error(err))
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
		logger.Warn("marshal", zap.Error(err))
	}
	return string(b)
}

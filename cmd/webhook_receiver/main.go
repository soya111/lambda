package main

import (
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

	"github.com/aws/aws-sdk-go/aws/session"
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

func handler(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	method := r.Method
	ip := r.RemoteAddr

	logger.Info("request", zap.String("ip", ip), zap.String("method", method), zap.String("path", "/webhook"), zap.String("body", string(body)))

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			logger.Error("Invalid signature", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid signature"))
		} else {
			logger.Error("Failed to parse request", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to parse request"))
		}
		return
	}

	var result *multierror.Error
	var wg sync.WaitGroup

	repo := dynamodb.NewSubscriberRepository(db)
	handler := webhook.NewHandler(bot, repo)

	ctx := logging.ContextWithLogger(r.Context(), logger)

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

	if result.ErrorOrNil() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to handle event"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully handled event"))
	}
}

func main() {
	http.HandleFunc("/webhook", handler)
	http.ListenAndServe(":8888", nil)
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

func marshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Warn("marshal", zap.Error(err))
	}
	return string(b)
}

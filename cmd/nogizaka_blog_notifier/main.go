package main

import (
	"time"

	n "notify/app/nogizaka_blog_notifier"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	// set timezone
	time.Local = time.FixedZone("JST", 9*60*60)
}

func main() {
	lambda.Start(n.ExcuteFunction)
}

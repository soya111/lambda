#!/bin/sh

while true; do
    go run cmd/hinatazaka_blog_notifier/main.go
    sleep 600
done

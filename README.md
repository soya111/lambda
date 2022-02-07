## build
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o handler ./cmd/hinatazaka_blog_notifier/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o handler ./cmd/sakurazaka_blog_notifier/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o handler ./cmd/nogizaka_blog_notifier/main.go

## zip
build-lambda-zip.exe -output handler.zip handler
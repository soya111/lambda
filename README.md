# LINE Messenger

## Build

Github Actionsで自動化した

localでビルドする場合

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o handler main.go
```

## Zip

Localでzipしたいときはこれ

インストール

```bash
go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest
```

```bash
cp handler.zip handler.zip.bk
build-lambda-zip.exe -output handler.zip handler
```

## Deploy

```bash
aws lambda update-function-code --function-name {{name}} --zip-file fileb://handler.zip 
```

## Coverage

```bash
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

[`coverage.html`](./coverage.html) をブラウザで開く

## About Hinatazaka46

Hinatazaka46 is a Japanese idol group.

Click ['here'](https://www.hinatazaka46.com) for more information.
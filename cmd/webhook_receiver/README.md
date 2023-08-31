# Localでwebhookを受け取るサーバーを立てる方法

## 1. ngrokをインストール

<https://ngrok.com/download>

## 2. ngrokを起動

```bash
ngrok http 8080
```

## 3. ngrokのURLを設定

URLをコピーして、LINE DevelopersのWebhook URLに設定する。

## 4. サーバーを起動

```bash
go run ./cmd/webhook_receiver/main.go
```

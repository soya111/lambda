# Localでwebhookを受け取るサーバーを立てる方法

## 1. Port forwarding local servicesの設定

vscodeのパネルからポートに移動し、"Forward a Port"を選択する。  
ポート番号:8080を追加後、右クリックでポートの表示範囲をPublicに変更する。

## 2. ローカルアドレスを設定

ローカルアドレスをコピーして、LINE DevelopersのWebhook URLに設定する。

## 3. サーバーを起動

```bash
go run ./cmd/webhook_receiver/main.go
```

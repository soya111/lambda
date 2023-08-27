# ローカルでの Notifier の使い方

ローカルでNotifierを使用するための手順を以下に示します。

## 準備

### 1. DynamoDB Localの準備

- **Dockerのインストール**: まず、Dockerをシステムにインストールします。
- **DynamoDB Localの起動**: 以下のコマンドを使用してDynamoDB Localを起動します。

    ```bash
    docker compose up --build -d
    ```

### 2. DynamoDB用のNoSQL Workbenchのダウンロード

- **ダウンロード**: [こちら](https://docs.aws.amazon.com/ja_jp/amazondynamodb/latest/developerguide/workbench.settingup.html)からダウンロードできます。
- **接続**: Add connectionをクリックし、`localhost:8000`へ接続します。

### 3. 個人用LINEbotアカウントを作成

- **Botの作成**: [このコンソール](https://developers.line.biz/console/)からbotを作成します。
- **情報の取得**: チャネルアクセストークン、チャネルシークレット、ユーザーIDを取得します。

### 4. envファイル変更

- **ファイルの編集**: 前の手順で取得した値に書き換えるためのコマンドは次のとおりです。

    ```bash
    cd pkg
    cp .env.sample .env
    code .env
    ```

## 実行

テストを実行する場合、以下のコマンドを使用します。

```bash
go test -timeout 30s -run ^TestExecute$ notify/pkg/notifier
```

t.Skip のコメントアウトを忘れずに!

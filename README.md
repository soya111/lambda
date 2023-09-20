# LINE Messenger Integration

This guide will take you through the steps to automate, build, zip, and deploy a LINE Messenger integration using Go, as well as executing it using Docker Compose.

---

## 🔨 Build

Use **GitHub Actions** for automation. If you wish to build locally:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o handler cmd/{{app_name}}/main.go
```

---

## 🗜️ Zip

To zip files locally:

1. **Installation**:

    ```bash
    go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest
    ```

1. **Zip Creation**:

    ```bash
    cp handler.zip handler.zip.bk
    build-lambda-zip.exe -output handler.zip handler
    ```

---

## 🚀 Deploy

Deploy to AWS Lambda:

```bash
aws lambda update-function-code --function-name {{name}} --zip-file fileb://handler.zip 
```

---

## 📊 Coverage

To obtain coverage statistics:

```bash
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

View the coverage report by opening [`coverage.html`](./coverage.html) in a browser.

---

## 🐳 Docker

To execute the integration:

1. **Environment File**: Create a .env file in the root directory to store environment variables such as access keys. An example could be:

    ```env
    CHANNEL_SECRET=your_channel_secret
    CHANNEL_TOKEN=your_channel_token
    ME=your_user_id
    IS_LOCAL=1
    ```

1. **Docker Compose**: Assuming you have a docker-compose.yml file, you can run the integration via:

    ```bash
    docker compose up
    ```

    For **subsequent runs**, if there are containers you prefer not to start, specifically start only the necessary ones:

    ```bash
    docker compose up [service-name]
    ```

Ensure that the Docker Compose file references the .env file for environment variables.

---

## 🎤 About Hinatazaka46

Hinatazaka46 is a captivating Japanese idol group. Dive deep into their world and learn about their incredible journey.

![Hinatazaka46](https://www.thefirsttimes.jp/admin/wp-content/uploads/5023/06/20230623-dm-100001.jpg)

[Discover More](https://www.hinatazaka46.com)

---

💡 Note: Replace placeholders like {{name}} and {{app_name}} with the appropriate values before execution.

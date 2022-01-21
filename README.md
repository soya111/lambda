## build
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o handler main.go

## zip
build-lambda-zip.exe handler handler
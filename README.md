## Build

Github Actionsで自動化した

## Things to improve on

- dynamoの操作をgureguで

## zip

Localでzipしたいときはこれ

Windowsでzipするとうごかないこともある

```bash
cp handler.zip handler.zip.bk
build-lambda-zip.exe -output handler.zip handler
aws lambda update-function-code --function-name {{name}} --zip-file fileb://handler.zip 
```

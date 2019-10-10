env GOOS=linux GOARCH=amd64 go build -o /tmp/main github.com/joel-ezell/serverless-hasher
zip -j /tmp/main.zip /tmp/main
aws lambda update-function-code --function-name hasher --zip-file fileb:///tmp/main.zip

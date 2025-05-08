package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/myarik/aws-samples/cdklambdatemplate/demo/handler"
	"github.com/myarik/aws-samples/cdklambdatemplate/demo/pkg/logger"
)

func main() {
	if err := logger.Init("info", "development"); err != nil {
		fmt.Printf("init logger error: %v\n", err)
	}
	defer logger.Sync()
	lambda.Start(handler.Run)
}

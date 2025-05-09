package handler

import (
	"context"
	"encoding/json"
	"github.com/myarik/aws-samples/cdklambdatemplate/demo/pkg/logger"
	"go.uber.org/zap"
)

// Run is the main handler function for the Lambda function.
func Run(ctx context.Context, event json.RawMessage) (map[string]interface{}, error) {
	// Log the raw JSON data
	log := logger.L()

	// Log structured data if needed
	log.Info("Receive a data", zap.Any("event", event))

	// Return a successful response
	return map[string]interface{}{
		"status":  "success",
		"message": "Data received and logged successfully",
	}, nil
}

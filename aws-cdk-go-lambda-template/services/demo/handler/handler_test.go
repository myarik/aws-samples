package handler

import (
	"context"
	"encoding/json"
	"github.com/myarik/aws-samples/cdklambdatemplate/demo/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	// Create test cases
	tests := []struct {
		name        string
		event       json.RawMessage
		wantStatus  string
		wantMessage string
		wantErr     bool
	}{
		{
			name:        "valid event",
			event:       json.RawMessage(`{"key": "value"}`),
			wantStatus:  "success",
			wantMessage: "Data received and logged successfully",
			wantErr:     false,
		},
		{
			name:        "empty event",
			event:       json.RawMessage(`{}`),
			wantStatus:  "success",
			wantMessage: "Data received and logged successfully",
			wantErr:     false,
		},
	}

	// Initialize logger for testing
	err := logger.Init("info", "development")
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function with test data
			got, err := Run(context.Background(), tt.event)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check response
			assert.Equal(t, tt.wantStatus, got["status"])
			assert.Equal(t, tt.wantMessage, got["message"])
		})
	}
}

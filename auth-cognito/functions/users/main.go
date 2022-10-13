package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type user struct {
	Key       int      `json:"key"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	Address   string   `json:"address"`
	Tags      []string `json:"tags"`
}

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.WithFields(log.Fields{
		"request": request.RequestContext.Authorizer,
	}).Info("Request received")

	respBody, err := json.Marshal([]user{
		{
			Key:       1,
			FirstName: "John",
			LastName:  "Doe",
			Age:       32,
			Address:   "1234 Main St",
			Tags:      []string{"nice", "developer"},
		},
		{
			Key:       2,
			FirstName: "Jane",
			LastName:  "Green",
			Age:       25,
			Address:   "5678 1st Ave W",
			Tags:      []string{"cool", "designer"},
		},
		{
			Key:       3,
			FirstName: "Jim",
			LastName:  "Carrey",
			Age:       55,
			Address:   "9012 Hollywood Blvd",
			Tags:      []string{"legend"},
		},
	})
	if err != nil {
		log.WithError(err).Error("cannot marshalling response")
		return events.APIGatewayProxyResponse{
			StatusCode:      http.StatusInternalServerError,
			Body:            string(respBody),
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization," +
					"X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            string(respBody),
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization," +
				"X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}

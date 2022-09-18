package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	rerror "github.com/myarik/vp/aws_image_thumbnail/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type media struct {
	MediaId      string `json:"media_id"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	ThumbnailULR string `json:"thumbnail_url"`
}

var (
	dynamodbClient *dynamodb.Client
	tableName      string
	staticURL      string
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.WithError(err).Fatal("Error loading AWS config")
	}
	dynamodbClient = dynamodb.NewFromConfig(awsConfig)
	tableName = os.Getenv("TABLE_NAME")

	staticURL = os.Getenv("STATIC_URL")
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	productId := request.PathParameters["id"]

	queryProjection := expression.NamesList(
		expression.Name("MediaId"),
		expression.Name("Type"),
		expression.Name("URL"),
		expression.Name("ThumbnailULR"),
	)

	queryExpr, err := expression.NewBuilder().WithKeyCondition(
		expression.Key("ProductId").Equal(expression.Value(productId))).WithProjection(
		queryProjection).Build()
	if err != nil {
		log.WithError(err).WithField("productId", productId).Error("cannot building query expression")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "")
	}

	queryOutput, err := dynamodbClient.Query(ctx, &dynamodb.QueryInput{
		KeyConditionExpression:    queryExpr.KeyCondition(),
		ProjectionExpression:      queryExpr.Projection(),
		ExpressionAttributeNames:  queryExpr.Names(),
		ExpressionAttributeValues: queryExpr.Values(),
		TableName:                 aws.String(tableName),
	})
	if err != nil {
		log.WithError(err).WithField("productId", productId).Error("cannot query dynamodb")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "")
	}

	response := make([]media, 0)
	for _, item := range queryOutput.Items {
		var m media
		if err = attributevalue.UnmarshalMap(item, &m); err != nil {
			log.WithError(err).WithField("productId", productId).Error("cannot unmarshal dynamodb item")
			return rerror.WrapStatusCode(http.StatusInternalServerError, "")
		}

		m.URL = fmt.Sprintf("%s%s", staticURL, m.URL)
		if m.ThumbnailULR != "" {
			m.ThumbnailULR = fmt.Sprintf("%s%s", staticURL, m.ThumbnailULR)
		}
		response = append(response, m)
	}

	respBody, err := json.Marshal(response)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
		}).Error("cannot marshalling response")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "cannot marshal response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            string(respBody),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

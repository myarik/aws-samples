package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/models"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

var (
	dynamodbClient *dynamodb.Client
	tableName      string
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
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		thumbnail := models.Thumbnail{}
		if err := json.Unmarshal([]byte(snsRecord.Message), &thumbnail); err != nil {
			log.WithError(err).WithField("event", record.EventSource).Error("cannot unmarshal event")
			continue
		}

		updateExpr, err := expression.NewBuilder().WithUpdate(expression.Set(
			expression.Name("ThumbnailULR"),
			expression.Value(thumbnail.URL))).Build()
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"productId": thumbnail.ProductId,
				"mediaId":   thumbnail.MediaId,
			}).Error("cannot update dynamodb expression")
		}

		_, err = dynamodbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"MediaId":   &types.AttributeValueMemberS{Value: thumbnail.MediaId},
				"ProductId": &types.AttributeValueMemberS{Value: thumbnail.ProductId},
			},
			UpdateExpression:          updateExpr.Update(),
			ExpressionAttributeNames:  updateExpr.Names(),
			ExpressionAttributeValues: updateExpr.Values(),
		})
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"productId":        thumbnail.ProductId,
				"mediaId":          thumbnail.MediaId,
				"URL":              thumbnail.URL,
				"UpdateExpression": updateExpr.Update(),
			}).Error("cannot update the item")
			continue
		}
		log.WithFields(log.Fields{
			"productId": thumbnail.ProductId,
			"mediaId":   thumbnail.MediaId,
		}).Info("thumbnail saved")
	}

}

func main() {
	lambda.Start(Handler)
}

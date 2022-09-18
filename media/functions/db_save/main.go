package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/models"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
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

		media := models.Media{}
		if err := json.Unmarshal([]byte(snsRecord.Message), &media); err != nil {
			log.WithError(err).Error("cannot unmarshalling message")
			continue
		}

		_, err := dynamodbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"ProductId": &types.AttributeValueMemberS{Value: media.ProductId},
				"MediaId":   &types.AttributeValueMemberS{Value: media.Id},
				"Type":      &types.AttributeValueMemberS{Value: string(media.Type)},
				"URL":       &types.AttributeValueMemberS{Value: media.URL},
				"CreatedAt": &types.AttributeValueMemberN{Value: strconv.FormatInt(media.CreatedAt, 10)},
			},
		})
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"mediaId":   media.Id,
				"productId": media.ProductId,
			}).Error("cannot insert an item")
			continue
		}
		log.WithFields(log.Fields{
			"mediaId":   media.Id,
			"productId": media.ProductId,
		}).Info("media created")
	}

}

func main() {
	lambda.Start(Handler)
}

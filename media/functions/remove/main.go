package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	typesDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	rerror "github.com/myarik/vp/aws_image_thumbnail/pkg/errors"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/event"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	dynamodbClient      *dynamodb.Client
	snsClient           *sns.Client
	topicARN, tableName string
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
	snsClient = sns.NewFromConfig(awsConfig)
	topicARN = os.Getenv("SNS_TOPIC_ARN")

	dynamodbClient = dynamodb.NewFromConfig(awsConfig)
	tableName = os.Getenv("TABLE_NAME")
}

// publishEvent sends a message to the SNS topic
func publishEvent(ctx context.Context, s3Key string) error {
	// Publish message to SNS topic
	input := &sns.PublishInput{
		Message:  aws.String(s3Key),
		TopicArn: aws.String(topicARN),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"event": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(event.RemoveObject)),
			},
		},
	}
	if _, err := snsClient.Publish(ctx, input); err != nil {
		return errors.Wrap(err, "cannot publish message to SNS topic")
	}
	log.WithField("Event", event.RemoveObject).Info("message published")
	return nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	productId := request.PathParameters["id"]
	mediaId := request.PathParameters["mediaId"]

	projectionExpr, err := expression.NewBuilder().WithProjection(
		expression.NamesList(expression.Name("URL"), expression.Name("ThumbnailULR"))).Build()
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
			"mediaId":   mediaId,
		}).Error("Error building projection expression")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "")
	}

	// TODO check how to use dynamodb.GetItemInput
	output, err := dynamodbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]typesDynamodb.AttributeValue{
			"MediaId":   &typesDynamodb.AttributeValueMemberS{Value: mediaId},
			"ProductId": &typesDynamodb.AttributeValueMemberS{Value: productId},
		},
		ExpressionAttributeNames: projectionExpr.Names(),
		ProjectionExpression:     projectionExpr.Projection(),
		TableName:                aws.String(tableName),
	})
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
			"mediaId":   mediaId,
		}).Error("Error getting item from DynamoDB")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "")
	}
	if output.Item == nil {
		log.WithFields(log.Fields{
			"productId": productId,
			"mediaId":   mediaId,
		}).Error("Item not found")
		return rerror.WrapStatusCode(http.StatusNotFound, "")
	}
	// decode item
	s3Keys := []string{
		output.Item["URL"].(*typesDynamodb.AttributeValueMemberS).Value,
		output.Item["ThumbnailULR"].(*typesDynamodb.AttributeValueMemberS).Value,
	}

	// send message to SNS topic
	for _, s3Key := range s3Keys {
		if err = publishEvent(ctx, s3Key); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"productId": productId,
				"mediaId":   mediaId,
			}).Error("cannot publish message to SNS topic")
		}
	}

	// delete item from DynamoDB
	_, err = dynamodbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]typesDynamodb.AttributeValue{
			"MediaId":   &typesDynamodb.AttributeValueMemberS{Value: mediaId},
			"ProductId": &typesDynamodb.AttributeValueMemberS{Value: productId},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
			"mediaId":   mediaId,
		}).Error("Error deleting item from DynamoDB")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "")
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusNoContent,
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

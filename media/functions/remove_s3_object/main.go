package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

var (
	s3Client *s3.Client
	bucket   string
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
	s3Client = s3.NewFromConfig(awsConfig)
	bucket = os.Getenv("MEDIA_BUCKET")
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(snsRecord.Message),
		})
		if err != nil {
			log.WithError(err).WithField("key", snsRecord.Message).Error("Error deleting object")
			return
		}
		log.WithField("key", snsRecord.Message).Info("Object deleted")
	}
}

func main() {
	lambda.Start(Handler)
}

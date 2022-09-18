package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/disintegration/imaging"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/event"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	snsClient *sns.Client
	s3Client  *s3.Client
	bucket    string

	imageWidth, imageHeight int
	topicARN                string
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
	// setup SNS client
	snsClient = sns.NewFromConfig(awsConfig)

	topicARN = os.Getenv("SNS_TOPIC_ARN")
	bucket = os.Getenv("MEDIA_BUCKET")
}

// downloadFile downloads the file from S3 and decodes it into an image.Image
func downloadFile(ctx context.Context, bucket string, key string) (image.Image, error) {
	// download the file from S3
	output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	// Decode the image from the buffer and return it
	return imaging.Decode(output.Body, imaging.AutoOrientation(true))
}

func uploadFile(ctx context.Context, key string, img image.Image) error {
	// encode the image into a buffer
	buf := new(bytes.Buffer)
	err := imaging.Encode(buf, img, imaging.PNG)
	if err != nil {
		return err
	}

	// upload the file to S3
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		ACL:         s3Types.ObjectCannedACLPublicRead,
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return err
	}
	return nil
}

func createThumbnail(ctx context.Context, media *models.Media) (*models.Thumbnail, error) {
	originImage, err := downloadFile(ctx, bucket, media.URL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot downloading file")
	}
	thumbnail := imaging.Thumbnail(originImage, imageWidth, imageHeight, imaging.CatmullRom)

	url := strings.Replace(
		media.URL, filepath.Ext(media.URL), fmt.Sprintf("-%dx%d_thumbnail.png", imageWidth, imageHeight), 1)

	if err = uploadFile(ctx, url, thumbnail); err != nil {
		return nil, errors.Wrap(err, "cannot uploading file")
	}
	return &models.Thumbnail{
		ProductId: media.ProductId,
		MediaId:   media.Id,
		URL:       url,
		Height:    imageHeight,
		Width:     imageWidth,
	}, nil
}

// publishEvent sends a message to the SNS topic
func publishEvent(ctx context.Context, item *models.Thumbnail) error {
	body, err := json.Marshal(item)
	if err != nil {
		return errors.Wrap(err, "cannot marshal message")
	}

	snsInput := &sns.PublishInput{
		Message:  aws.String(string(body)),
		TopicArn: aws.String(topicARN),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"event": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(event.ThumbnailCreated)),
			},
		},
	}

	_, err = snsClient.Publish(ctx, snsInput)
	if err != nil {
		return errors.Wrap(err, "cannot publishing message")
	}
	log.WithField("Event", event.ThumbnailCreated).Info("message published")

	return nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		media := models.Media{}
		if err := json.Unmarshal([]byte(snsRecord.Message), &media); err != nil {
			log.WithError(err).WithField("event", record.EventSource).Error("cannot unmarshal event")
			continue
		}

		thumbnail, err := createThumbnail(ctx, &media)
		if err != nil {
			log.WithFields(log.Fields{
				"mediaId":   media.Id,
				"productId": media.ProductId,
			}).WithError(err).Error("cannot creating thumbnail")
			continue
		}

		if err = publishEvent(ctx, thumbnail); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"event":     event.ThumbnailCreated,
				"mediaId":   media.Id,
				"productId": media.ProductId,
			}).Error("cannot send message to SNS")
			continue
		}
		log.WithFields(log.Fields{
			"URL":       thumbnail.URL,
			"mediaId":   media.Id,
			"productId": media.ProductId,
		}).Info("thumbnail created")
	}

}

func main() {
	var err error

	imageWidth, err = strconv.Atoi(os.Getenv("THUMBNAIL_WIDTH"))
	if err != nil {
		log.Fatal("Invalid THUMBNAIL_WIDTH parameter")
	}

	imageHeight, err = strconv.Atoi(os.Getenv("THUMBNAIL_HEIGHT"))
	if err != nil {
		log.Fatal("Invalid THUMBNAIL_HEIGHT parameter")
	}

	lambda.Start(Handler)
}

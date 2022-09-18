package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/grokify/go-awslambda"
	rerror "github.com/myarik/vp/aws_image_thumbnail/pkg/errors"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/event"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/models"
	"github.com/myarik/vp/aws_image_thumbnail/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	snsClient *sns.Client
	topicARN  string

	s3Client            *s3.Client
	bucket, mediaPrefix string

	supportedFormats = map[string]bool{
		"image/gif":  true, // GIF
		"image/jpeg": true, // JPEG
		"image/jpg":  true, // JPG
		"image/png":  true, // PNG
		"video/mp4":  true, // MP4
	}

	mediaTypes = map[string]models.MediaType{
		"image/gif":  models.Image, // GIF
		"image/jpeg": models.Image, // JPEG
		"image/jpg":  models.Image, // JPG
		"image/png":  models.Image, // PNG
		"video/mp4":  models.Video, // MP4
	}
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
	bucket, mediaPrefix = os.Getenv("MEDIA_BUCKET"), os.Getenv("MEDIA_PREFIX")

}

// publishEvent sends a message to the SNS topic
func publishEvent(ctx context.Context, event event.Event, body string) error {
	snsInput := &sns.PublishInput{
		Message:  aws.String(body),
		TopicArn: aws.String(topicARN),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"event": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(event)),
			},
		},
	}

	_, err := snsClient.Publish(ctx, snsInput)
	if err != nil {
		return errors.Wrap(err, "cannot publishing message")
	}
	log.WithField("Event", event).Info("message published")
	return nil
}

// uploadMedia uploads a media file to S3
func uploadMedia(ctx context.Context, bucket, key string, body []byte, contentType string) error {
	// upload the file to S3
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		ACL:         s3Types.ObjectCannedACLPublicRead,
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return errors.Wrap(err, "cannot upload media")
	}
	return nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	productId := request.PathParameters["id"]

	multipartReader, err := awslambda.NewReaderMultipart(request)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
		}).Error("cannot creating multipart reader")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "")
	}

	part, err := multipartReader.NextPart()

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
		}).Error("cannot getting next part")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "")
	}
	defer func(part *multipart.Part) {
		err = part.Close()
		if err != nil {
			log.Warn("cannot closing part")
		}
	}(part)

	// Check file format
	if !supportedFormats[part.Header.Get("Content-Type")] {
		log.WithFields(log.Fields{
			"Content-Type": part.Header.Get("Content-Type"),
			"productId":    productId,
		}).Warn("unsupported content type")

		return rerror.WrapCode(rerror.CodeInvalidRequestError, "unsupported content type")
	}

	body, err := ioutil.ReadAll(part)

	mediaId := utils.RandString()

	fileName := fmt.Sprintf("%s%s", mediaId, filepath.Ext(part.FileName()))
	url := fmt.Sprintf(
		"%s/%s/%s",
		mediaPrefix,
		productId,
		fileName,
	)

	if err = uploadMedia(ctx, bucket, url, body, part.Header.Get("Content-Type")); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
		}).Error("cannot uploading media")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "db_save returns an error")
	}

	media := models.Media{
		Id:        mediaId,
		ProductId: productId,
		Type:      mediaTypes[part.Header.Get("Content-Type")],
		URL:       url,
		CreatedAt: time.Now().Unix(),
	}

	respBody, err := json.Marshal(media)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
			"mediaId":   mediaId,
		}).Error("cannot marshalling response")
		return rerror.WrapStatusCode(http.StatusInternalServerError, "cannot marshal response")
	}

	if err = publishEvent(ctx, event.MediaUploaded, string(respBody)); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"productId": productId,
			"mediaId":   mediaId,
			"event":     event.MediaUploaded,
		}).Error("cannot send message to SNS")
	}

	if media.Type == models.Image {
		if err = publishEvent(ctx, event.CreateThumbnail, string(respBody)); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"productId": productId,
				"mediaId":   mediaId,
				"event":     event.CreateThumbnail,
			}).Error("cannot send message to SNS")
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            string(respBody),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

package util

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hillview.tv/coreAPI/env"
	"github.com/hillview.tv/coreAPI/media"
)

var (
	contentClient = func() *s3.Client {
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(env.ContentS3Region),
		)
		if err != nil {
			panic(fmt.Errorf("unable to load aws credentials: %w", err))
		}

		// An S3-compatible store (OpenBucket) replaces AWS when configured.
		if env.ContentS3Endpoint == "" {
			return s3.NewFromConfig(cfg)
		}
		cfg.Credentials = credentials.NewStaticCredentialsProvider(env.ContentS3AccessKeyID, env.ContentS3SecretAccessKey, "")
		return s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(env.ContentS3Endpoint)
			o.UsePathStyle = true
			// The SDK's default (checksums "when supported") sends PutObject
			// bodies as aws-chunked with a CRC32 trailer. OpenBucket doesn't
			// decode that and stores the chunk framing inside the object, so
			// every upload came out corrupt. Only checksum when an operation
			// requires it.
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		})
	}()
)

type UploadImageRequest struct {
	Image []byte
	ID    int
	Route string
}

func UploadImage(ctx context.Context, req UploadImageRequest) (string, error) {
	contentType, ext, err := media.ImageType(req.Image)
	if err != nil {
		return "", err
	}

	now := time.Now()
	sec := now.Unix()
	generated := strconv.Itoa(req.ID) + "-" + RandomString(10, LetterNumberBytes) + "-" + strconv.FormatInt(sec, 10)
	key := req.Route + generated + ext
	input := &s3.PutObjectInput{
		Bucket:      aws.String(env.ContentS3Bucket),
		Key:         aws.String(key),
		ACL:         types.ObjectCannedACLPublicRead,
		Body:        bytes.NewReader(req.Image),
		ContentType: aws.String(contentType),
	}

	_, err = contentClient.PutObject(ctx, input)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error putting object in s3: %w", err)
	}

	return env.ContentPublicBaseURL + "/" + key, nil
}

func ValidS3Route(route string) bool {
	switch route {
	case "images/assets/":
		return true
	case "thumbnails/":
		return true
	default:
		return false
	}
}

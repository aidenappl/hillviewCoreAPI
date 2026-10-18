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
		})
	}()
)

type UploadImageRequest struct {
	Image []byte
	ID    int
	Route string
}

func UploadImage(ctx context.Context, req UploadImageRequest) (string, error) {
	now := time.Now()
	sec := now.Unix()
	generated := strconv.Itoa(req.ID) + "-" + RandomString(10, LetterNumberBytes) + "-" + strconv.FormatInt(sec, 10)
	input := &s3.PutObjectInput{
		Bucket:      aws.String(env.ContentS3Bucket),
		Key:         aws.String(req.Route + *aws.String(generated) + ".jpeg"),
		ACL:         types.ObjectCannedACLPublicRead,
		Body:        bytes.NewReader(req.Image),
		ContentType: aws.String("image/jpeg"),
	}

	_, err := contentClient.PutObject(ctx, input)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error putting object in s3: %w", err)
	}

	return env.ContentPublicBaseURL + "/" + req.Route + generated + ".jpeg", nil
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

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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	s3ClientUSWest1 = s3.NewFromConfig(awsConfig)

	awsConfig = func() aws.Config {
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
		)
		if err != nil {
			panic(fmt.Errorf("unable to load aws credentials: %w", err))
		}

		return cfg
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
		Bucket:      aws.String("content.hillview.tv"),
		Key:         aws.String(req.Route + *aws.String(generated) + ".jpeg"),
		ACL:         types.ObjectCannedACLPublicRead,
		Body:        bytes.NewReader(req.Image),
		ContentType: aws.String("image/jpeg"),
	}

	_, err := s3ClientUSWest1.PutObject(ctx, input)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error putting object in s3: %w", err)
	}

	return "https://content.hillview.tv/" + req.Route + generated + ".jpeg", nil
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

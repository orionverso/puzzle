package receiver

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3ReceiverProps struct {
	putObjectInput s3.PutObjectInput
}

type s3Receiver struct {
	s3.Client
	destination string
}

func (rv s3Receiver) Write(ctx context.Context, st string) error {
	var sprops s3ReceiverProps = s3ReceiverProps_DEFAULT

	body := []byte(st)

	sprops.putObjectInput.Bucket = aws.String(rv.destination)
	sprops.putObjectInput.Key = aws.String(randstr(10))
	sprops.putObjectInput.Body = bytes.NewReader(body)

	_, err := rv.PutObject(ctx, &sprops.putObjectInput)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("data delivered")

	return nil
}

// SETTINGS
var s3ReceiverProps_DEFAULT s3ReceiverProps = s3ReceiverProps{
	putObjectInput: s3.PutObjectInput{
		ContentType: aws.String("application/json"),
	},
}

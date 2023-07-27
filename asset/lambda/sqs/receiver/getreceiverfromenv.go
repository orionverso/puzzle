package receiver

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Receiver interface {
	Write(ctx context.Context, st string) error
}

func GetReceiverFromEnv(ctx context.Context, cfg aws.Config) Receiver {
	storage := os.Getenv("STORAGE_SERVICE")

	if storage == "" {
		log.Panicln("The storage service service is empty")
	}

	destination := os.Getenv("DESTINATION")

	if storage == "" {
		log.Panicln("The destination is empty")
	}

	if storage == "DYNAMODB" {
		return dynamoDbReceiver{
			*dynamodb.NewFromConfig(cfg),
			destination,
		}
	}

	if storage == "S3" {
		return s3Receiver{
			*s3.NewFromConfig(cfg),
			destination,
		}
	}

	return nil
}

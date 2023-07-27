package main

import (
	"context"
	"echo/receiver"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

func handler(ctx context.Context, snsEvent events.SNSEvent) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Panicln("Configuration has not been loaded")
	}

	rec := receiver.GetReceiverFromEnv(ctx, cfg) // Change service destination if change env vars
	if err != nil {
		log.Panicln("Receive has not been loaded")
	}

	for _, record := range snsEvent.Records {
		err = rec.Write(ctx, record.SNS.Message)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

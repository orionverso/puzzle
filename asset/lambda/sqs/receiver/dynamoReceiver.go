package receiver

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type itemProps struct {
	index             string
	attributeValue    types.AttributeValue
	attributeValueKey types.AttributeValue
}

type putItemInputProps struct {
	item map[string]types.AttributeValue
	itemProps
}

type dynamoDbReceiverProps struct {
	putItemInput dynamodb.PutItemInput
	putItemInputProps
}

type dynamoDbReceiver struct {
	dynamodb.Client
	destination string
}

func (rv dynamoDbReceiver) Write(ctx context.Context, st string) error {
	var sprops dynamoDbReceiverProps = dynamoDbReceiverProps_DEFAULT

	sprops.putItemInput.TableName = aws.String(rv.destination)
	sprops.attributeValue = &types.AttributeValueMemberS{Value: st}

	sprops.putItemInput.Item[sprops.index] = sprops.attributeValue

	_, err := rv.PutItem(ctx, &sprops.putItemInput)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// SETTINGS
// ---DEFAULT---
var itemProps_DEFAULT itemProps = itemProps{
	index:             "id",
	attributeValueKey: &types.AttributeValueMemberS{Value: "primaryKey"},
	attributeValue:    &types.AttributeValueMemberS{Value: "default-value"},
}

var putItemInputProps_DEFAULT putItemInputProps = putItemInputProps{
	itemProps: itemProps_DEFAULT,
	item: map[string]types.AttributeValue{
		itemProps_DEFAULT.index: itemProps_DEFAULT.attributeValueKey,
	},
}

var dynamoDbReceiverProps_DEFAULT dynamoDbReceiverProps = dynamoDbReceiverProps{
	putItemInputProps: putItemInputProps_DEFAULT,
	putItemInput: dynamodb.PutItemInput{
		Item: putItemInputProps_DEFAULT.item,
	},
}

// ---

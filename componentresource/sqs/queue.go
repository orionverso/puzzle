package sqs

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CompanyQueue struct {
	pulumi.ResourceState
	QueueArn  pulumi.StringOutput
	QueueName pulumi.StringOutput
	Queue     pulumi.IDOutput
}

type CompanyQueueArgs struct {
	sqs.QueueArgs
}

func NewCompanyQueue(ctx *pulumi.Context, name string, args *CompanyQueueArgs, opts ...pulumi.ResourceOption) (*CompanyQueue, error) {
	componentResource := &CompanyQueue{}

	if args == nil {
		args = &CompanyQueueArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:sqs:CompanyQueue", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	qu, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-queue", name), &args.QueueArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{
		"QueueName": qu.Name,
		"QueueArn":  qu.Arn,
		"Queue":     qu.ID(),
	})

	ctx.Export("QueueName", qu.Name)
	ctx.Export("QueueArn", qu.Arn)
	ctx.Export("Queue", qu.ID())

	return componentResource, nil
}

package pieces

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CompanyQueue struct {
	pulumi.ResourceState
	QueueArn  pulumi.StringOutput
	QueueName pulumi.StringOutput
	QueueID   pulumi.IDOutput
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
	err := ctx.RegisterComponentResource("puzzle:pieces:CompanyQueue", name, componentResource, opts...)
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
		"QueueID":   qu.ID(),
	})

	ctx.Export("QueueName", qu.Name)
	ctx.Export("QueueArn", qu.Arn)
	ctx.Export("QueueID", qu.ID())

	componentResource.QueueArn = qu.Arn
	componentResource.QueueName = qu.Name
	componentResource.QueueID = qu.ID()

	return componentResource, nil
}

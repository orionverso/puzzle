package board

import (
	"fmt"
	"puzzle/pieces"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type QueueFuncBucket struct {
	pulumi.ResourceState
	FuncBucket
	pieces.CompanyQueue
}

type QueueFuncBucketArgs struct {
	FuncBucketArgs
	pieces.CompanyQueueArgs
	lambda.EventSourceMappingArgs
}

func NewQueueFuncBucket(ctx *pulumi.Context, name string, args *QueueFuncBucketArgs, opts ...pulumi.ResourceOption) (*QueueFuncBucket, error) {
	componentResource := &QueueFuncBucket{}

	if args == nil {
		args = &QueueFuncBucketArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:board:MessageTriggerFuncTable", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	qu, err := pieces.NewCompanyQueue(ctx, "companyqueue", &args.CompanyQueueArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	inlineargs := iam.RoleInlinePolicyArgs{
		Name: pulumi.String("LambdaQueueBasicPolicy"),
		Policy: pulumi.Sprintf(`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action":       [
                "sqs:ReceiveMessage",
                "sqs:DeleteMessage",
                "sqs:GetQueueAttributes",
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "%s"
        }
    ]
}`, qu.QueueArn),
	}

	args.FuncBucketArgs.CompanyFuncArgs.AppendPolicyToInlinePolicies(inlineargs)

	fnbk, err := NewFuncBucket(ctx, "funcbucket", &args.FuncBucketArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	args.EventSourceMappingArgs.EventSourceArn = qu.QueueArn
	args.EventSourceMappingArgs.FunctionName = fnbk.FunctionName

	_, err = lambda.NewEventSourceMapping(ctx, fmt.Sprintf("%s-eventsourcemapping", name), &args.EventSourceMappingArgs, pulumi.Parent(componentResource), pulumi.DependsOn([]pulumi.Resource{qu, fnbk}))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

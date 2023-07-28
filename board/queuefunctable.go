package board

import (
	"fmt"
	"puzzle/pieces"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type QueueFuncTable struct {
	pulumi.ResourceState
	FuncTable
	pieces.CompanyQueue
}

type QueueFuncTableArgs struct {
	FuncTableArgs
	pieces.CompanyQueueArgs
	lambda.EventSourceMappingArgs
}

func NewQueueFuncTable(ctx *pulumi.Context, name string, args *QueueFuncTableArgs, opts ...pulumi.ResourceOption) (*QueueFuncTable, error) {
	componentResource := &QueueFuncTable{}

	if args == nil {
		args = &QueueFuncTableArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource(fmt.Sprintf("puzzle:board:%s", triggerwriteandsave), name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	qu, err := pieces.NewCompanyQueue(ctx, trigger, &args.CompanyQueueArgs, pulumi.Parent(componentResource))

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

	args.FuncTableArgs.CompanyFuncArgs.AppendPolicyToInlinePolicies(inlineargs)

	fntb, err := NewFuncTable(ctx, writeandsave, &args.FuncTableArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	args.EventSourceMappingArgs.EventSourceArn = qu.QueueArn
	args.EventSourceMappingArgs.FunctionName = fntb.FunctionName

	_, err = lambda.NewEventSourceMapping(ctx, fmt.Sprintf("%s-eventsourcemapping", name), &args.EventSourceMappingArgs, pulumi.Parent(componentResource), pulumi.DependsOn([]pulumi.Resource{qu, fntb}))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

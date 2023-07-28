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

	args.FuncBucketArgs.CompanyFuncArgs.AppendPolicyToInlinePolicies(inlineargs)

	fnbk, err := NewFuncBucket(ctx, writeandsave, &args.FuncBucketArgs, pulumi.Parent(componentResource))

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

//This is an example of static configuration to put in main()
// args := board.QueueFuncBucketArgs{
// 	CompanyQueueArgs: pieces.CompanyQueueArgs{
// 		QueueArgs: sqs.QueueArgs{},
// 	},
//
// 	FuncBucketArgs: board.FuncBucketArgs{
// 		CompanyFuncArgs: pieces.CompanyFuncArgs{
// 			RoleArgs: iam.RoleArgs{
// 				AssumeRolePolicy: pulumi.String(`{
//     "Version": "2012-10-17",
//     "Statement": [
//         {
//             "Effect": "Allow",
//             "Action": [
//                 "sts:AssumeRole"
//             ],
//             "Principal": {
//                 "Service": [
//                     "lambda.amazonaws.com"
//                 ]
//             }
//         }
//     ]
// }`),
// 			},
// 			FunctionArgs: lambda.FunctionArgs{
// 				Runtime:     pulumi.StringPtr("go1.x"),
// 				Code:        pulumi.NewFileArchive("./asset/lambda/sqs/handler.zip"),
// 				Handler:     pulumi.StringPtr("handler"),
// 				Description: pulumi.StringPtr("This function goes to write to bucket"),
// 				Timeout:     pulumi.IntPtr(5),
// 			},
// 		},
// 		CompanyBucketArgs: pieces.CompanyBucketArgs{
//
// 			BucketV2Args: s3.BucketV2Args{
// 				ForceDestroy: pulumi.BoolPtr(true),
// 			},
// 			BucketServerSideEncryptionConfigurationV2Args: s3.BucketServerSideEncryptionConfigurationV2Args{
// 				Rules: s3.BucketServerSideEncryptionConfigurationV2RuleArray{
// 					s3.BucketServerSideEncryptionConfigurationV2RuleArgs{
// 						ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationV2RuleApplyServerSideEncryptionByDefaultArgs{
// 							KmsMasterKeyId: pulumi.StringPtr("alias/aws/s3"),
// 							SseAlgorithm:   pulumi.String("aws:kms"),
// 						},
// 						BucketKeyEnabled: pulumi.BoolPtr(true),
// 					},
// 				},
// 			},
// 			BucketVersioningV2Args: s3.BucketVersioningV2Args{
// 				VersioningConfiguration: s3.BucketVersioningV2VersioningConfigurationArgs{
// 					Status: pulumi.String("Enabled"),
// 				},
// 			},
// 			BucketLoggingV2Args: s3.BucketLoggingV2Args{},
// 		},
// 	},
// }
//
// board.NewQueueFuncBucket(ctx, "QueueTriggerLambdaWriteToStorage", &args)

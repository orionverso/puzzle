package board

import (
	"fmt"
	"puzzle/pieces"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type TopicFuncBucket struct {
	pulumi.ResourceState
	FuncBucket
	pieces.CompanyTopic
}

type TopicFuncBucketArgs struct {
	FuncBucketArgs
	pieces.CompanyTopicArgs
	lambda.EventSourceMappingArgs
}

func NewTopicFuncBucket(ctx *pulumi.Context, name string, args *TopicFuncBucketArgs, opts ...pulumi.ResourceOption) (*TopicFuncBucket, error) {
	componentResource := &TopicFuncBucket{}

	if args == nil {
		args = &TopicFuncBucketArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource(fmt.Sprintf("puzzle:board:%s", triggerwriteandsave), name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	tp, err := pieces.NewCompanyTopic(ctx, trigger, &args.CompanyTopicArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	inlineargs := iam.RoleInlinePolicyArgs{
		Name: pulumi.String("LambdaTopicBasicPolicy"),
		Policy: pulumi.Sprintf(`{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "Statement1",
			"Effect": "Allow",
			"Action": "sns:GetTopicAttributes",
			"Resource": "%s"
		}
	]
}`, tp.TopicArn),
	}

	args.FuncBucketArgs.CompanyFuncArgs.AppendPolicyToInlinePolicies(inlineargs)

	fnbk, err := NewFuncBucket(ctx, writeandsave, &args.FuncBucketArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	permargs := lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  fnbk.FunctionArn,
		Principal: pulumi.String("sns.amazonaws.com"),
		SourceArn: tp.TopicArn,
	}

	_, err = lambda.NewPermission(ctx, fmt.Sprintf("%s-sns-trigger-permision", name), &permargs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	_, err = sns.NewTopicSubscription(ctx, fmt.Sprintf("%s-lambda-subscription", name), &sns.TopicSubscriptionArgs{
		Topic:    tp.TopicArn,
		Protocol: pulumi.String("lambda"),
		Endpoint: fnbk.FunctionArn,
	}, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

//This is an example of static configuration to put in main()
// args := board.TopicFuncBucketArgs{
// 	CompanyTopicArgs: pieces.CompanyTopicArgs{
// 		TopicArgs: sqs.TopicArgs{},
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
// board.NewTopicFuncBucket(ctx, "TopicTriggerLambdaWriteToStorage", &args)

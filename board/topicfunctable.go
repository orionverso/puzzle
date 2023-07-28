package board

import (
	"fmt"
	"puzzle/pieces"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type TopicFuncTable struct {
	pulumi.ResourceState
	FuncTable
	pieces.CompanyTopic
}

type TopicFuncTableArgs struct {
	FuncTableArgs
	pieces.CompanyTopicArgs
	lambda.EventSourceMappingArgs
}

func NewTopicFuncTable(ctx *pulumi.Context, name string, args *TopicFuncTableArgs, opts ...pulumi.ResourceOption) (*TopicFuncTable, error) {
	componentResource := &TopicFuncTable{}

	if args == nil {
		args = &TopicFuncTableArgs{}
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

	args.FuncTableArgs.CompanyFuncArgs.AppendPolicyToInlinePolicies(inlineargs)

	fntb, err := NewFuncTable(ctx, writeandsave, &args.FuncTableArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	permargs := lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  fntb.FunctionArn,
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
		Endpoint: fntb.FunctionArn,
	}, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

//This is an example of static configuration to put in main()
// args := board.TopicFuncTableArgs{
// 	CompanyTopicArgs: pieces.CompanyTopicArgs{
// 		TopicArgs: sns.TopicArgs{},
// 	},
// 	FuncTableArgs: board.FuncTableArgs{
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
// 				Code:        pulumi.NewFileArchive("./asset/lambda/sns/handler.zip"),
// 				Handler:     pulumi.StringPtr("handler"),
// 				Description: pulumi.StringPtr("This function goes to write to table"),
// 				Timeout:     pulumi.IntPtr(5),
// 			},
// 		},
// 		CompanyTableArgs: pieces.CompanyTableArgs{
// 			TableArgs: dynamodb.TableArgs{
// 				Attributes: dynamodb.TableAttributeArray{dynamodb.TableAttributeArgs{
// 					Name: pulumi.String("id"),
// 					Type: pulumi.String("S"),
// 				},
// 				},
// 				HashKey:       pulumi.StringPtr("id"),
// 				ReadCapacity:  pulumi.IntPtr(5),
// 				WriteCapacity: pulumi.IntPtr(5),
// 			},
// 		},
// 	},
// }
//
// board.NewTopicFuncTable(ctx, "TopicTriggerLambdaWriteToStorage", &args)

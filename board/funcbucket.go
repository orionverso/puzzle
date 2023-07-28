package board

import (
	"fmt"
	"puzzle/pieces"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type FuncBucket struct {
	pulumi.ResourceState
}

type FuncBucketArgs struct {
	pieces.CompanyFuncArgs
	pieces.CompanyBucketArgs
}

func NewFuncBucket(ctx *pulumi.Context, name string, args *FuncBucketArgs, opts ...pulumi.ResourceOption) (*FuncBucket, error) {
	componentResource := &FuncBucket{}

	if args == nil {
		args = &FuncBucketArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:board:LambdaStorage", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	bk, err := pieces.NewCompanyBucket(ctx, fmt.Sprintf("%s-companybucket", name), &args.CompanyBucketArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	//dynamyc configurations
	args.CompanyFuncArgs.RoleArgs.InlinePolicies = iam.RoleInlinePolicyArray{iam.RoleInlinePolicyArgs{
		Name: pulumi.String("WriteToDynamoDb"),
		Policy: pulumi.Sprintf(`{ 
                    "Version": "2012-10-17",
                    "Statement": [{
                        "Effect": "Allow",
                        "Action": "s3:PutObject",
                        "Resource": "%s/*"
                    }]
                }`, bk.BucketArn), //asynchronous value
	},
	}

	args.CompanyFuncArgs.Environment = lambda.FunctionEnvironmentArgs{
		Variables: pulumi.ToStringMapOutput(map[string]pulumi.StringOutput{
			"STORAGE_SERVICE": pulumi.Sprintf("S3"),
			"DESTINATION":     pulumi.Sprintf("%s", bk.BucketName), //asynchronous value
		}),
	}

	_, err = pieces.NewCompanyFunc(ctx, fmt.Sprintf("%s-companyfunc", name), &args.CompanyFuncArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

//This is an example of static configuration to put in main()

// 		args := board.FuncBucketArgs{
// 			CompanyFuncArgs: pieces.CompanyFuncArgs{
// 				RoleArgs: iam.RoleArgs{
// 					AssumeRolePolicy: pulumi.String(`{
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
// 				},
// 				FunctionArgs: lambda.FunctionArgs{
// 					Runtime:     pulumi.StringPtr("go1.x"),
// 					Code:        pulumi.NewFileArchive("./asset/lambda/sqs/handler.zip"),
// 					Handler:     pulumi.StringPtr("handler"),
// 					Description: pulumi.StringPtr("This function goes to write to table"),
// 					Timeout:     pulumi.IntPtr(5),
// 				},
// 			},
// 			CompanyBucketArgs: pieces.CompanyBucketArgs{
//
// 				BucketV2Args: s3.BucketV2Args{
// 					ForceDestroy: pulumi.BoolPtr(true),
// 				},
// 				BucketServerSideEncryptionConfigurationV2Args: s3.BucketServerSideEncryptionConfigurationV2Args{
// 					Rules: s3.BucketServerSideEncryptionConfigurationV2RuleArray{
// 						s3.BucketServerSideEncryptionConfigurationV2RuleArgs{
// 							ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationV2RuleApplyServerSideEncryptionByDefaultArgs{
// 								KmsMasterKeyId: pulumi.StringPtr("alias/aws/s3"),
// 								SseAlgorithm:   pulumi.String("aws:kms"),
// 							},
// 							BucketKeyEnabled: pulumi.BoolPtr(true),
// 						},
// 					},
// 				},
// 				BucketVersioningV2Args: s3.BucketVersioningV2Args{
// 					VersioningConfiguration: s3.BucketVersioningV2VersioningConfigurationArgs{
// 						Status: pulumi.String("Enabled"),
// 					},
// 				},
// 				BucketLoggingV2Args: s3.BucketLoggingV2Args{},
// 			},
// 		}
//
//

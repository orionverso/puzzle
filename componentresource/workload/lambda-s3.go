package workload

import (
	"fmt"
	mylb "puzzle/componentresource/lambda"
	mys3 "puzzle/componentresource/s3"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LambdaS3 struct {
	pulumi.ResourceState
}

type LambdaS3Args struct {
	mylb.CompanyFuncArgs
	mys3.CompanyBucketArgs
}

func NewLambdaS3(ctx *pulumi.Context, name string, args *LambdaS3Args, opts ...pulumi.ResourceOption) (*LambdaS3, error) {
	componentResource := &LambdaS3{}

	if args == nil {
		args = &LambdaS3Args{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:workload:LambdaStorage", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	bk, err := mys3.NewCompanyBucket(ctx, fmt.Sprintf("%s-companybucket", name), &args.CompanyBucketArgs, pulumi.Parent(componentResource))

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

	_, err = mylb.NewCompanyFunc(ctx, fmt.Sprintf("%s-companyfunc", name), &args.CompanyFuncArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

//This is an example of static configuration to put in main()

// 		args := workload.LambdaS3Args{
// 			CompanyFuncArgs: mylb.CompanyFuncArgs{
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
// 			CompanyBucketArgs: mys3.CompanyBucketArgs{
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

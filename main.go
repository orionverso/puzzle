package main

import (
	"puzzle/board"
	"puzzle/pieces"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		//static configuration
		args := board.QueueFuncBucketArgs{
			CompanyQueueArgs: pieces.CompanyQueueArgs{
				QueueArgs: sqs.QueueArgs{},
			},

			FuncBucketArgs: board.FuncBucketArgs{
				CompanyFuncArgs: pieces.CompanyFuncArgs{
					RoleArgs: iam.RoleArgs{
						AssumeRolePolicy: pulumi.String(`{
		    "Version": "2012-10-17",
		    "Statement": [
		        {
		            "Effect": "Allow",
		            "Action": [
		                "sts:AssumeRole"
		            ],
		            "Principal": {
		                "Service": [
		                    "lambda.amazonaws.com"
		                ]
		            }
		        }
		    ]
		}`),
					},
					FunctionArgs: lambda.FunctionArgs{
						Runtime:     pulumi.StringPtr("go1.x"),
						Code:        pulumi.NewFileArchive("./asset/lambda/sqs/handler.zip"),
						Handler:     pulumi.StringPtr("handler"),
						Description: pulumi.StringPtr("This function goes to write to bucket"),
						Timeout:     pulumi.IntPtr(5),
					},
				},
				CompanyBucketArgs: pieces.CompanyBucketArgs{

					BucketV2Args: s3.BucketV2Args{
						ForceDestroy: pulumi.BoolPtr(true),
					},
					BucketServerSideEncryptionConfigurationV2Args: s3.BucketServerSideEncryptionConfigurationV2Args{
						Rules: s3.BucketServerSideEncryptionConfigurationV2RuleArray{
							s3.BucketServerSideEncryptionConfigurationV2RuleArgs{
								ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationV2RuleApplyServerSideEncryptionByDefaultArgs{
									KmsMasterKeyId: pulumi.StringPtr("alias/aws/s3"),
									SseAlgorithm:   pulumi.String("aws:kms"),
								},
								BucketKeyEnabled: pulumi.BoolPtr(true),
							},
						},
					},
					BucketVersioningV2Args: s3.BucketVersioningV2Args{
						VersioningConfiguration: s3.BucketVersioningV2VersioningConfigurationArgs{
							Status: pulumi.String("Enabled"),
						},
					},
					BucketLoggingV2Args: s3.BucketLoggingV2Args{},
				},
			},
		}

		board.NewQueueFuncBucket(ctx, "QueueTriggerLambdaWriteToStorage", &args)

		return nil
	})
}

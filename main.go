package main

import (
	mys3 "puzzle/componentresource/s3"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		args := mys3.CompanyBucketArgs{

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
		}

		mys3.NewCompanyBucket(ctx, "Storage", &args)

		return nil
	})
}

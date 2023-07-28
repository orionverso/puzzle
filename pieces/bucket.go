package pieces

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CompanyBucket struct {
	pulumi.ResourceState
	BucketArn  pulumi.StringOutput
	BucketName pulumi.StringOutput
	BucketID   pulumi.IDOutput
}

type CompanyBucketArgs struct {
	s3.BucketV2Args
	s3.BucketServerSideEncryptionConfigurationV2Args
	s3.BucketVersioningV2Args
	s3.BucketLoggingV2Args
}

func NewCompanyBucket(ctx *pulumi.Context, name string, args *CompanyBucketArgs, opts ...pulumi.ResourceOption) (*CompanyBucket, error) {
	componentResource := &CompanyBucket{}

	if args == nil {
		args = &CompanyBucketArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:pieces:CompanyBucket", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	bk, err := s3.NewBucketV2(ctx, fmt.Sprintf("%s-bucket", name), &args.BucketV2Args, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	args.BucketServerSideEncryptionConfigurationV2Args.Bucket = bk.ID()

	_, err = s3.NewBucketServerSideEncryptionConfigurationV2(ctx, fmt.Sprintf("%s-server-side-encrytion", name), &args.BucketServerSideEncryptionConfigurationV2Args, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	args.BucketVersioningV2Args.Bucket = bk.ID()

	_, err = s3.NewBucketVersioningV2(ctx, fmt.Sprintf("%s-versioning", name), &args.BucketVersioningV2Args, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{
		"BucketArn":  bk.Arn,
		"BucketName": bk.Bucket,
		"BucketID":   bk.ID(),
	})

	ctx.Export("BucketArn", bk.Arn)
	ctx.Export("BucketName", bk.Bucket)
	ctx.Export("BucketID", bk.ID())

	componentResource.BucketArn = bk.Arn
	componentResource.BucketName = bk.Bucket
	componentResource.BucketID = bk.ID()

	return componentResource, nil
}

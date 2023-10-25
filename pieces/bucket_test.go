package pieces

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestCompanyBucket(t *testing.T) {
	t.Parallel()

	clean := func(ctx context.Context, stackname string, stack auto.Stack) {

		_, err := stack.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))

		if err != nil {
			fmt.Printf("Failed to resources: %v", err)
			fmt.Println("You must remove the reources manually")
			t.FailNow()
		}

		fmt.Println("Stack resources successfully destroyed")

		err = stack.Workspace().RemoveStack(ctx, stackname)

		if err != nil {
			fmt.Printf("Failed to destroy stack: %v", err)
			fmt.Println("You must remove the stack manually")
			t.FailNow()
		}

		fmt.Println("Stack successfully destroyed")
	}

	resources := func(ctx *pulumi.Context) error {
		args := CompanyBucketArgs{

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

		_, err := NewCompanyBucket(ctx, "MyCompanyBucketTesting", &args)

		if err != nil {
			t.FailNow()
		}
		return nil
	}

	ctx := context.Background()
	org := "orionverso"
	projectName := "Test"
	stackName := auto.FullyQualifiedStackName(org, projectName, "ephemeral-bucket")

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, resources)

	fmt.Printf("Created/Selected stack %q\n", stackName)

	w := stack.Workspace()

	defer clean(ctx, stackName, stack) // PASS or FAIL: Clean resources!!

	//Install dependecies
	fmt.Println("Installing the AWS plugin")

	err = w.InstallPlugin(ctx, "aws", "v5.0.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully installed AWS plugin")

	//Set stack configuration
	stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "us-west-2"})

	fmt.Println("Successfully set config")

	//Deploy Resources
	fmt.Println("Start stack up")

	_, err = stack.Up(ctx, optup.ProgressStreams(os.Stdout))

	if err != nil {
		fmt.Printf("Failed to deploy stack resources: %v\n", err)
		t.FailNow()
	}

	fmt.Println("Successfully stack resources deployment")

	/*
		You can run some integration testing here with aws-sdk-go.
		For example, put object to bucket and check
	*/

}

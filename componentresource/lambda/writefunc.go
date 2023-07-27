package lambda

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CompanyFunc struct {
	pulumi.ResourceState
	FunctionArn  pulumi.StringOutput
	FunctionName pulumi.StringOutput
	FunctionID   pulumi.IDOutput
}

type CompanyFuncArgs struct {
	lambda.FunctionArgs
	iam.RoleArgs
}

func NewCompanyFunc(ctx *pulumi.Context, name string, args *CompanyFuncArgs, opts ...pulumi.ResourceOption) (*CompanyFunc, error) {
	componentResource := &CompanyFunc{}

	if args == nil {
		args = &CompanyFuncArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:lamdba:CompanyFunc", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	rl, err := iam.NewRole(ctx, fmt.Sprintf("%s-role", name), &args.RoleArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	args.FunctionArgs.Role = rl.Arn

	fn, err := lambda.NewFunction(ctx, fmt.Sprintf("%s-function", name), &args.FunctionArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{
		"FunctionArn":  fn.Arn,
		"FunctionName": fn.Name,
		"FunctionID":   fn.ID(),
	})

	ctx.Export("FunctionArn", fn.Arn)
	ctx.Export("FunctionName", fn.Name)
	ctx.Export("FunctionID", fn.ID())

	componentResource.FunctionArn = fn.Arn
	componentResource.FunctionName = fn.Name
	componentResource.FunctionID = fn.ID()

	return componentResource, nil
}

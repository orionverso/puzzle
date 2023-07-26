package lambda

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type WriteFunc struct {
	pulumi.ResourceState
	Function pulumi.IDOutput
}

type WriteFuncArgs struct {
	lambda.FunctionArgs
	iam.RoleArgs
}

func NewWriteFunc(ctx *pulumi.Context, name string, args *WriteFuncArgs, opts ...pulumi.ResourceOption) (*WriteFunc, error) {
	componentResource := &WriteFunc{}

	if args == nil {
		args = &WriteFuncArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:lamdba:WriteFunc", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	rl, err := iam.NewRole(ctx, fmt.Sprintf("%s-role", name), &args.RoleArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	args.FunctionArgs.Role = rl.ID()

	fn, err := lambda.NewFunction(ctx, fmt.Sprintf("%s-function", name), &args.FunctionArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	componentResource.Function = fn.ID()

	return componentResource, nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		//Resource Declaration

		if err != nil {
			return err
		}

		return nil
	})
}

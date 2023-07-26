package dynamodb

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CompanyTable struct {
	pulumi.ResourceState
}

type CompanyTableArgs struct {
	dynamodb.TableArgs
}

func NewCompanyTable(ctx *pulumi.Context, name string, args *CompanyTableArgs, opts ...pulumi.ResourceOption) (*CompanyTable, error) {
	componentResource := &CompanyTable{}

	if args == nil {
		args = &CompanyTableArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:dynamodb:CompanyTable", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	_, err = dynamodb.NewTable(ctx, fmt.Sprintf("%s-table", name), &args.TableArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

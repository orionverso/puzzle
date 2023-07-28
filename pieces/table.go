package pieces

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CompanyTable struct {
	pulumi.ResourceState
	TableArn  pulumi.StringOutput
	TableName pulumi.StringOutput
	TableID   pulumi.IDOutput
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
	err := ctx.RegisterComponentResource("puzzle:pieces:CompanyTable", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	tb, err := dynamodb.NewTable(ctx, fmt.Sprintf("%s-table", name), &args.TableArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{
		"TableArn":  tb.Arn,
		"TableName": tb.Name,
		"TableID":   tb.ID(),
	})

	ctx.Export("TableArn", tb.Arn)
	ctx.Export("TableName", tb.Name)
	ctx.Export("TableID", tb.ID())

	componentResource.TableArn = tb.Arn
	componentResource.TableName = tb.Name
	componentResource.TableID = tb.ID()

	return componentResource, nil
}

package board

import (
	"fmt"
	"puzzle/pieces"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type FuncTable struct {
	pulumi.ResourceState
	*pieces.CompanyFunc
	*pieces.CompanyTable
}

type FuncTableArgs struct {
	pieces.CompanyFuncArgs
	pieces.CompanyTableArgs
}

func NewFuncTable(ctx *pulumi.Context, name string, args *FuncTableArgs, opts ...pulumi.ResourceOption) (*FuncTable, error) {
	componentResource := &FuncTable{}
	// awsconf := config.New(ctx, "aws")
	// region := awsconf.Get("region")

	if args == nil {
		args = &FuncTableArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:board:LambdaStorage", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	tb, err := pieces.NewCompanyTable(ctx, fmt.Sprintf("%s-companytable", name), &args.CompanyTableArgs, pulumi.Parent(componentResource))

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
                        "Action": "dynamodb:PutItem",
                        "Resource": "%s"
                    }]
                }`, tb.TableArn), //asynchronous value
	},
	}

	args.CompanyFuncArgs.Environment = lambda.FunctionEnvironmentArgs{
		Variables: pulumi.ToStringMapOutput(map[string]pulumi.StringOutput{
			"STORAGE_SERVICE": pulumi.Sprintf("DYNAMODB"),
			"DESTINATION":     pulumi.Sprintf("%s", tb.TableName), //asynchronous value
		}),
	}

	fn, err := pieces.NewCompanyFunc(ctx, fmt.Sprintf("%s-companyfunc", name), &args.CompanyFuncArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	componentResource.CompanyFunc = fn
	componentResource.CompanyTable = tb

	return componentResource, nil
}

//This is an example of static configuration to put in main()

// 		args := board.FuncTableArgs{
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
// 			CompanyTableArgs: pieces.CompanyTableArgs{
// 				TableArgs: dynamodb.TableArgs{
// 					Attributes: dynamodb.TableAttributeArray{dynamodb.TableAttributeArgs{
// 						Name: pulumi.String("id"),
// 						Type: pulumi.String("S"),
// 					},
// 					},
// 					HashKey:       pulumi.StringPtr("id"),
// 					ReadCapacity:  pulumi.IntPtr(5),
// 					WriteCapacity: pulumi.IntPtr(5),
// 				},
// 			},
// 		}

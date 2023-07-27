package workload

import (
	"fmt"
	mydb "puzzle/componentresource/dynamodb"
	mylb "puzzle/componentresource/lambda"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LambdaDb struct {
	pulumi.ResourceState
}

type LambdaDbArgs struct {
	mylb.CompanyFuncArgs
	mydb.CompanyTableArgs
}

func NewLambdaDb(ctx *pulumi.Context, name string, args *LambdaDbArgs, opts ...pulumi.ResourceOption) (*LambdaDb, error) {
	componentResource := &LambdaDb{}
	// awsconf := config.New(ctx, "aws")
	// region := awsconf.Get("region")

	if args == nil {
		args = &LambdaDbArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:workload:LambdaStorage", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	tb, err := mydb.NewCompanyTable(ctx, fmt.Sprintf("%s-companytable", name), &args.CompanyTableArgs, pulumi.Parent(componentResource))

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

	_, err = mylb.NewCompanyFunc(ctx, fmt.Sprintf("%s-companyfunc", name), &args.CompanyFuncArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{})

	return componentResource, nil
}

//This is an example of static configuration to put in main()

// 		args := workload.LambdaDbArgs{
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
// 			CompanyTableArgs: mydb.CompanyTableArgs{
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

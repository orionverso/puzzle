package main

import (
	wk "puzzle/componentresource/workload"

	mydb "puzzle/componentresource/dynamodb"
	mylb "puzzle/componentresource/lambda"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		//static configurations
		args := wk.LambdaDbArgs{
			CompanyFuncArgs: mylb.CompanyFuncArgs{
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
					Description: pulumi.StringPtr("This function goes to write to table"),
					Timeout:     pulumi.IntPtr(5),
				},
			},
			CompanyTableArgs: mydb.CompanyTableArgs{
				TableArgs: dynamodb.TableArgs{
					Attributes: dynamodb.TableAttributeArray{dynamodb.TableAttributeArgs{
						Name: pulumi.String("id"),
						Type: pulumi.String("S"),
					},
					},
					HashKey:       pulumi.StringPtr("id"),
					ReadCapacity:  pulumi.IntPtr(5),
					WriteCapacity: pulumi.IntPtr(5),
				},
			},
		}

		wk.NewLambdaDb(ctx, "MyLambdaDb", &args)

		return nil
	})
}

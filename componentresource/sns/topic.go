package sns

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CompanyTopic struct {
	pulumi.ResourceState
}

type CompanyTopicArgs struct {
	sns.TopicArgs
	sns.TopicSubscriptionArgs
}

func NewCompanyTopic(ctx *pulumi.Context, name string, args *CompanyTopicArgs, opts ...pulumi.ResourceOption) (*CompanyTopic, error) {
	componentResource := &CompanyTopic{}

	if args == nil {
		args = &CompanyTopicArgs{}
	}

	// <package>:<module>:<type>
	err := ctx.RegisterComponentResource("puzzle:sns:CompanyTopic", name, componentResource, opts...)
	if err != nil {
		return nil, err
	}

	tp, err := sns.NewTopic(ctx, fmt.Sprintf("%s-topic", name), &args.TopicArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	args.TopicSubscriptionArgs.Topic = tp.ID()

	_, err = sns.NewTopicSubscription(ctx, fmt.Sprintf("%s-subscription", name), &args.TopicSubscriptionArgs, pulumi.Parent(componentResource))

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(componentResource, pulumi.Map{
		"TopicName": tp.Name,
		"TopicArn":  tp.Name,
		"Topic":     tp.ID(),
	})

	ctx.Export("TopicName", tp.Name)
	ctx.Export("TopicArn", tp.Name)
	ctx.Export("Topic", tp.ID())

	return componentResource, nil
}

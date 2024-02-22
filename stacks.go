package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StackProps struct {
	awscdk.StackProps
}

func NewStack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, jsii.String("AsanaPolice"), &sprops)

	f := awslambda.NewDockerImageFunction(stack, jsii.String("Lambda"), &awslambda.DockerImageFunctionProps{
		Code: awslambda.DockerImageCode_FromImageAsset(
			jsii.String("lambda/"),
			&awslambda.AssetImageCodeProps{},
		),
		Architecture: awslambda.Architecture_ARM_64(),
		Timeout:      awscdk.Duration_Minutes(jsii.Number(1)),
		MemorySize:   jsii.Number(256),
		LogRetention: awslogs.RetentionDays_TWO_WEEKS,
	})
	f.Role().AddManagedPolicy(
		awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMReadOnlyAccess")),
	)

	rule := awsevents.NewRule(stack, jsii.String("ScheduleRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Minute:  jsii.String("0"),
			Hour:    jsii.String("0"),
			WeekDay: jsii.String("TUE,THU"),
		}),
	})
	rule.AddTarget(
		awseventstargets.NewLambdaFunction(f, &awseventstargets.LambdaFunctionProps{}),
	)

	f2 := awslambda.NewDockerImageFunction(stack, jsii.String("Lambda2"), &awslambda.DockerImageFunctionProps{
		Code: awslambda.DockerImageCode_FromImageAsset(
			jsii.String("lambda2/"),
			&awslambda.AssetImageCodeProps{},
		),
		Architecture: awslambda.Architecture_ARM_64(),
		Timeout:      awscdk.Duration_Minutes(jsii.Number(1)),
		MemorySize:   jsii.Number(256),
		LogRetention: awslogs.RetentionDays_TWO_WEEKS,
	})
	f2.Role().AddManagedPolicy(
		awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMReadOnlyAccess")),
	)

	rule2 := awsevents.NewRule(stack, jsii.String("ScheduleRule2"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Minute:  jsii.String("0"),
			Hour:    jsii.String("0"),
			WeekDay: jsii.String("WED"),
		}),
	})
	rule2.AddTarget(
		awseventstargets.NewLambdaFunction(f2, &awseventstargets.LambdaFunctionProps{}),
	)

	return stack
}

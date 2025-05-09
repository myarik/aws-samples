package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CalendarSyncStackProps struct {
	awscdk.StackProps
}

func NewTStack(scope constructs.Construct, id string, props *CalendarSyncStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	//// Create a Lambda function that prints "Hello World"
	helloWorldLambda := awslambda.NewFunction(stack, jsii.Sprintf("%s-Demo", id), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		// Runtime requires you to use bootstrap as the executable name
		Handler:      jsii.String("bootstrap"),
		Architecture: awslambda.Architecture_ARM_64(),
		//AWS CDK will automatically zip the contents of this directory during deployment.
		Code:         awslambda.Code_FromAsset(jsii.String("../.build/lambdas/demo"), nil),
		FunctionName: jsii.Sprintf("%s-Demo", id),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(10)),

		Environment: &map[string]*string{
			"Environment": jsii.String("staging"),
			"LOG_LEVEL":   jsii.String("info"),
			"APP_VERSION": jsii.String("0.0.1"),
		},
		LogRetention:          awslogs.RetentionDays_THREE_DAYS,
		LoggingFormat:         awslambda.LoggingFormat_JSON,
		SystemLogLevelV2:      awslambda.SystemLogLevel_INFO,
		ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO,
	})

	//// Output the Lambda function name
	awscdk.NewCfnOutput(stack, jsii.String("FunctionName"), &awscdk.CfnOutputProps{
		Value: helloWorldLambda.FunctionName(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewTStack(app, "LambdaStack", &CalendarSyncStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil
}

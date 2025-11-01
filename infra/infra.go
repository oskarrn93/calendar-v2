package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awss3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
	validator "github.com/oskarrn93/calendar-v2/internal/validation"
)

type AppConfig struct {
	RapidApiKey string `validate:"required"`
}

func (a *AppConfig) Validate() error {
	return validator.ValidateStruct(a)
}

func ReadRequiredEnvironmentVariables() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file was provided")
	}

	rapidApiKey := os.Getenv("RAPIDAPI_KEY")

	appConfig := AppConfig{
		RapidApiKey: rapidApiKey,
	}

	if err := appConfig.Validate(); err != nil {
		log.Fatalf("Failed to initialize app config due to missing variables: %v", err)
	}

	return appConfig
}

type InfraStackProps struct {
	awscdk.StackProps
}

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	envConfig := ReadRequiredEnvironmentVariables()

	stack := awscdk.NewStack(scope, &id, &sprops)

	s3Bucket := awss3.NewBucket(stack, jsii.String("calendar-v2-bucket"), &awss3.BucketProps{
		BucketName:    aws.String("calendar-oskarrosen-io"),
		AccessControl: awss3.BucketAccessControl_PRIVATE,
	})

	lambda := awslambda.NewFunction(stack, jsii.String("calendar-v2-bucket-lambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("bootstrap"),
		Code: awslambda.Code_FromDockerBuild(aws.String("./"), &awslambda.DockerBuildAssetOptions{
			File: aws.String("Dockerfile"),
		}),
		Environment: &map[string]*string{
			"RAPIDAPI_KEY":   &envConfig.RapidApiKey,
			"S3_BUCKET_NAME": s3Bucket.BucketName(),
		},
	})

	awsevents.NewRule(stack, jsii.String("calendar-v2-bucket-lambda-schedule"), &awsevents.RuleProps{
		Enabled: aws.Bool(true),
		Targets: &[]awsevents.IRuleTarget{
			awseventstargets.NewLambdaFunction(lambda, nil),
		},
		EventPattern: &awsevents.EventPattern{},
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Hour:   aws.String("6"),
			Minute: aws.String("0"),
		}),
	})

	s3ObjectPatternPermission := "*"

	s3Bucket.GrantReadWrite(lambda.Role(), s3ObjectPatternPermission)

	originAccessIdentity := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("calendar-v2-cf-oai"), &awscloudfront.OriginAccessIdentityProps{})
	s3Bucket.GrantRead(originAccessIdentity, s3ObjectPatternPermission)

	originConfigs := []*awscloudfront.SourceConfiguration{
		{
			Behaviors: &[]*awscloudfront.Behavior{
				{
					IsDefaultBehavior: aws.Bool(true),
					PathPattern:       &s3ObjectPatternPermission,
					AllowedMethods:    awscloudfront.CloudFrontAllowedMethods_GET_HEAD,
				},
			},
			S3OriginSource: &awscloudfront.S3OriginConfig{
				S3BucketSource:       s3Bucket,
				OriginAccessIdentity: originAccessIdentity,
			},
		},
	}

	awscloudfront.NewCloudFrontWebDistribution(stack, jsii.String("calendar-v2-cf"), &awscloudfront.CloudFrontWebDistributionProps{
		OriginConfigs: &originConfigs,
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewInfraStack(app, "calendar-v2", &InfraStackProps{
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

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

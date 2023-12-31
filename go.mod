module github.com/tanaka-takurou/serverless-simple-cms-go

go 1.21

replace github.com/tanaka-takurou/serverless-simple-cms-go/management/api/controller => ./management/api/controller

require (
	github.com/aws/aws-lambda-go latest
	github.com/aws/aws-sdk-go-v2 latest
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream latest
	github.com/aws/aws-sdk-go-v2/config latest
	github.com/aws/aws-sdk-go-v2/credentials latest
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue latest
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression latest
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds latest
	github.com/aws/aws-sdk-go-v2/feature/s3/manager latest
	github.com/aws/aws-sdk-go-v2/internal/configsources latest
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 latest
	github.com/aws/aws-sdk-go-v2/internal/ini latest
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider latest
	github.com/aws/aws-sdk-go-v2/service/dynamodb latest
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams latest
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding latest
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery latest
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url latest
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared latest
	github.com/aws/aws-sdk-go-v2/service/s3 latest
	github.com/aws/aws-sdk-go-v2/service/sso latest
	github.com/aws/aws-sdk-go-v2/service/sts latest
	github.com/aws/smithy-go latest
	github.com/jmespath/go-jmespath latest
	github.com/tanaka-takurou/serverless-simple-cms-go/management/api/controller latest
)

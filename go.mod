module github.com/tanaka-takurou/serverless-simple-cms-go

go 1.16

replace github.com/tanaka-takurou/serverless-simple-cms-go/management/api/controller => ./management/api/controller

require (
	github.com/aws/aws-lambda-go v1.22.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.2.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.1.1 // indirect
	github.com/tanaka-takurou/serverless-simple-cms-go/management/api/controller v0.0.0-20210205040859-f8821394f0ef // indirect
)

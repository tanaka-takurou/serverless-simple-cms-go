module github.com/tanaka-takurou/serverless-simple-cms-go

go 1.15

replace github.com/tanaka-takurou/serverless-simple-cms-go/management/api/controller => ./management/api/controller

require (
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go-v2 v0.23.0
	github.com/tanaka-takurou/serverless-simple-cms-go/management/api/controller v0.0.0-00010101000000-000000000000
)

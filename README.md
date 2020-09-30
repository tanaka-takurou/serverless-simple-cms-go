# serverless-simple-content-management-system kit
Simple kit for serverless simple content management system using AWS Lambda.


## Dependence

##### Golang
- aws-lambda-go
- aws-sdk-go-v2

##### Javascript
- jQuery

##### CSS
- semantic-ui


## Requirements
- AWS (Lambda, API Gateway, DynamoDB, S3, Cognito)

##### To fix the system
- golang environment
- aws-sam-cli


## Deploy
```bash
make clean build
AWS_PROFILE={profile} AWS_DEFAULT_REGION={region} make bucket={bucket} stack={stack name} deploy
```

Otherwise, You can deploy by Serverless Application Repository.


## Usage

### Front Page
Access the URL displayed as "FrontURI"

### Management Page
Access the URL displayed as "ManagementURI"

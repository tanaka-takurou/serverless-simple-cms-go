package controller

import (
	"os"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var cognitoClient *cognitoidentityprovider.Client

func Login(ctx context.Context, name string, pass string)(string, error) {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.New(cfg)
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: cognitoidentityprovider.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": name,
			"PASSWORD": pass,
		},
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
	}

	req := cognitoClient.InitiateAuthRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	return aws.StringValue(res.InitiateAuthOutput.AuthenticationResult.AccessToken), nil
}

func GetUser(ctx context.Context, token string)(string, error) {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.New(cfg)
	}

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	}

	req := cognitoClient.GetUserRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	return aws.StringValue(res.GetUserOutput.Username), nil
}

func ChangePass(ctx context.Context, token string, pass string, newPass string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.New(cfg)
	}

	input := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(token),
		PreviousPassword: aws.String(pass),
		ProposedPassword: aws.String(newPass),
	}

	req := cognitoClient.ChangePasswordRequest(input)
	_, err := req.Send(ctx)
	return err
}

func Logout(ctx context.Context, token string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.New(cfg)
	}

	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(token),
	}

	req := cognitoClient.GlobalSignOutRequest(input)
	_, err := req.Send(ctx)
	return err
}

func Signup(ctx context.Context, name string, pass string, mail string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.New(cfg)
	}

	ua := &cognitoidentityprovider.AttributeType {
		Name: aws.String("email"),
		Value: aws.String(mail),
	}
	input := &cognitoidentityprovider.SignUpInput{
		Username: aws.String(name),
		Password: aws.String(pass),
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		UserAttributes: []cognitoidentityprovider.AttributeType{
			*ua,
		},
	}

	req := cognitoClient.SignUpRequest(input)
	_, err := req.Send(ctx)
	return err
}

func ConfirmSignup(ctx context.Context, name string, confirmationCode string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.New(cfg)
	}

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		Username: aws.String(name),
		ConfirmationCode: aws.String(confirmationCode),
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
	}

	req := cognitoClient.ConfirmSignUpRequest(input)
	_, err := req.Send(ctx)
	return err
}

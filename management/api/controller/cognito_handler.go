package controller

import (
	"os"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

var cognitoClient *cognitoidentityprovider.Client

func Login(ctx context.Context, name string, pass string)(string, error) {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": name,
			"PASSWORD": pass,
		},
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
	}

	res, err := cognitoClient.InitiateAuth(ctx, input)
	if err != nil {
		return "", err
	}
	return aws.ToString(res.AuthenticationResult.AccessToken), nil
}

func GetUser(ctx context.Context, token string)(string, error) {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	}

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	}

	res, err := cognitoClient.GetUser(ctx, input)
	if err != nil {
		return "", err
	}
	return aws.ToString(res.Username), nil
}

func ChangePass(ctx context.Context, token string, pass string, newPass string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	}

	input := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(token),
		PreviousPassword: aws.String(pass),
		ProposedPassword: aws.String(newPass),
	}

	_, err := cognitoClient.ChangePassword(ctx, input)
	return err
}

func Logout(ctx context.Context, token string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	}

	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(token),
	}

	_, err := cognitoClient.GlobalSignOut(ctx, input)
	return err
}

func Signup(ctx context.Context, name string, pass string, mail string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	}

	ua := &types.AttributeType {
		Name: aws.String("email"),
		Value: aws.String(mail),
	}
	input := &cognitoidentityprovider.SignUpInput{
		Username: aws.String(name),
		Password: aws.String(pass),
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		UserAttributes: []types.AttributeType{
			*ua,
		},
	}

	_, err := cognitoClient.SignUp(ctx, input)
	return err
}

func ConfirmSignup(ctx context.Context, name string, confirmationCode string) error {
	if cognitoClient == nil {
		cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	}

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		Username: aws.String(name),
		ConfirmationCode: aws.String(confirmationCode),
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
	}

	_, err := cognitoClient.ConfirmSignUp(ctx, input)
	return err
}

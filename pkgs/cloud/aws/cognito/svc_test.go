package cognito_test

import (
	"fmt"
	"testing"

	"gitea/pcp-inariam/inariam/pkgs/cloud/aws"
	awstest "gitea/pcp-inariam/inariam/pkgs/cloud/aws"
)

var (
	EMAIL, PASSWORD, APPKEY, CLIENT_ID = "yassinebk23@gmail.com", "12345678ABCD*$ùa", "6c41us69rh65bnehqss8h889b9", ""
)

func TestCognito(t *testing.T) {

	credentials := awstest.GetTestCredentials()
	fmt.Println(credentials)

	sess, err := aws.OpenSession(&credentials)
	if err != nil {
		panic("Error opening session")
	}

	sess.CreateCognitoSvc(CLIENT_ID, APPKEY)

	// t.Run("SimpleSignIn", func(t *testing.T) {

	// 	result, err := sess.CognitoSvc.SimpleSignIn(EMAIL, PASSWORD)

	// 	fmt.Println("Printing here", result, "err", err)

	// 	if err != nil {
	// 		t.Errorf("Error: %s", err)
	// 	} else {
	// 		t.Logf("Result: %s", result)
	// 	}

	// 	fmt.Printf("Result: %s", result)
	// })

	// t.Run("Signup", func(t *testing.T) {

	// 	result, err := sess.CognitoSvc.SignUp(EMAIL, PASSWORD)

	// 	if err != nil {
	// 		t.Errorf("Error: %s", err)
	// 	} else {
	// 		t.Logf("Result: %s", result)
	// 	}

	// 	fmt.Printf("Result: %s", result)
	// })

	// t.Run("ConfirmSignup", func(t *testing.T) {

	// 	result, err := sess.CognitoSvc.ConfirmSignUp(EMAIL, "305741")

	// 	if err != nil {
	// 		t.Errorf("Error: %s", err)
	// 	} else {
	// 		t.Logf("Result: %s", result)
	// 	}

	// 	t.Logf("Result: %s", result)
	// })

	t.Run("Initiate Auth", func(t *testing.T) {
		res, err := sess.CognitoSvc.StartSignInProcess(EMAIL, PASSWORD)
		if err != nil {
			t.Fail()
		}

		res2, err := sess.CognitoSvc.GenerateMFAActivationCode(res.SessionKey, EMAIL)
		fmt.Println("Error generating MFA", err)
		if err != nil {
			t.Fail()
		}

		// fmt.Println("Completing MFA SETUP")

		sess.CognitoSvc.ConfirmMFAActivation("SAMSING", "643351", res2.Session)
		// sess.CognitoSvc.CompleteMFASetup(EMAIL, "MA2RZFJBDHB7QV6KQRPJSHUAUK6W7ZONIBCBB4AWEITDZIZNGAQA", res.Session)

	})

	// t.Run("Add MFA", func(t *testing.T) {
	// 	result, err := sess.CognitoSvc.SimpleSignIn(EMAIL, PASSWORD)
	// 	if err != nil {
	// 		t.Errorf("Error %s", err)
	// 	}

	// 	qrCode, err := sess.CognitoSvc.GenerateMFAActivationCode(result)
	// 	if err != nil {
	// 		t.Errorf("Error %s", err)
	// 	}
	// 	fmt.Println(qrCode)
	// })

	// t.Run("Resend Confirmation code", func(t *testing.T) {
	// 	err := sess.CognitoSvc.ResendConfirmSignUp(EMAIL)
	// 	if err != nil {
	// 		t.Errorf("Error %s", err)
	// 	}

	// })
}

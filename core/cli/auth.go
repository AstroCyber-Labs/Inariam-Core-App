package cli

import (
	"fmt"
	"gitea/pcp-inariam/inariam/core/config"
	inaAws "gitea/pcp-inariam/inariam/pkgs/cloud/aws"
	"gitea/pcp-inariam/inariam/pkgs/log"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "auth [command]",
	Short: "auth",
	Long:  "authentication-command",
}

func newSignUpCmd() *cobra.Command {
	var email, password string
	var createUserCmd = &cobra.Command{
		Use:   "create -e [email] -p [password]",
		Short: "Create a new user account",
		Run: func(cmd *cobra.Command, args []string) {

			cfg, err := config.New()
			if err != nil {
				log.Logger.Panicf("error loading configuration %w", err)
			}

			awsSess, err := retrieveAwsSession(cfg)
			if err != nil {
				log.Logger.Panicln(err)
			}

			user, err := awsSess.CognitoSvc.SignUp(email, password)

			if err != nil {
				log.Logger.Panicln(err)
			}

			color.Green(fmt.Sprintf("[+] User created successfully %s", user))
		},
	}

	createUserCmd.Flags().StringVarP(&email, "email", "e", "", "The email of the user (required)")

	createUserCmd.Flags().StringVarP(&password, "password", "p", "", "The password of the user (required)")

	return createUserCmd
}

// retrieveAwsSession retrieve AwsSession with the config passed to the handler
func retrieveAwsSession(currentConfig *config.Config) (*inaAws.Session, error) {
	creds := inaAws.Credentials{
		AccessKeyID:     currentConfig.AWS.AccessKeyID,
		SecretAccessKey: currentConfig.AWS.SecretAccessKey,
		Region:          currentConfig.AWS.Region,
	}
	awsSession, err := inaAws.OpenSession(&creds)

	if err != nil {
		return nil, err
	}

	awsSession.CreateCognitoSvc(currentConfig.CognitoConfig.ClientId, currentConfig.CognitoConfig.AppSecret)

	return awsSession, nil
}

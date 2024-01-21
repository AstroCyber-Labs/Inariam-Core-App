package aws

import (
	"os"

	"github.com/joho/godotenv"
)

func GetTestCredentials() Credentials {

	if err := godotenv.Load("/home/askee/.inariam/.env"); err != nil {
		panic(err)
	}
	return Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SESSION_TOKEN"),
		Region:          os.Getenv("AWS_DEFAULT_REGION"),
	}
}

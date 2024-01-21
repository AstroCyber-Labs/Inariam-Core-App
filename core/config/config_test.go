package config

import (
	"fmt"
	"testing"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

func TestConfigFlow(t *testing.T) {
	config, err := CreateDefaultConfig()
	if err != nil {
		return
	}

	if err := config.CreateConfigFileIfNotExist(); err != nil {
		fmt.Println("Error creating configuration file:", err)
		return
	}

	config.AWS = &AWSConfig{
		AccessKeyID:     "new-aws-access-key",
		SecretAccessKey: "new-aws-secret-key",
	}
	config.Azure =
		&AzureConfig{
			ClientID:       "new-azure-client-id",
			ClientSecret:   "new-azure-client-secret",
			SubscriptionID: "new-azure-subscription-id",
		}
	config.GCP = &GCPConfig{
		ProjectID: "new-gcp-project-idaab",
	}

	err = config.UpdateConfig()
	if err != nil {
		log.Logger.Errorln("Error updating configuration:", err)
	} else {
		log.Logger.Infoln("Configuration updated successfully.")
	}
	againAConfig, err := CreateDefaultConfig()
	if err != nil {
		log.Logger.Panicln(err)
	}
	err = againAConfig.ParseConfigFile()
	if err != nil {
		log.Logger.Errorln("Error parsing configuration file:", err)
		return
	}

	// Check if GCP Project ID is not nil before printing
	if config.GCP != nil {
		log.Logger.Errorln("GCP Project ID:", config.GCP.ProjectID)
	} else {
		log.Logger.Errorln("GCP Project ID is nil")
	}

	// Check if AWS Access Key ID is not nil before printing
	if config.AWS != nil {
		log.Logger.Errorln("AWS Access Key ID:", config.AWS.AccessKeyID)
	} else {
		log.Logger.Infoln("AWS Access Key ID is nil")
	}

	// Check if Azure Client ID is not nil before printing
	if config.Azure != nil {
		log.Logger.Errorln("Azure Client ID:", config.Azure.ClientID)
	} else {
		log.Logger.Infoln("Azure Client ID is nil")
	}
}

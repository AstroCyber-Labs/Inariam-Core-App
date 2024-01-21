package config

// GCPConfig represents the configuration settings specific to Google Cloud Platform (GCP).
type GCPConfig struct {
	ProjectID       string `mapstructure:"project_id" yaml:"project_id"`
	CredentialsPath string `mapstructure:"credentials_path" yaml:"credentials_path"`
	// Add other GCP-specific fields here
}

// AWSConfig represents the configuration settings specific to Amazon Web Services (AWS).
type AWSConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id" yaml:"access_key_id" json:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key" yaml:"secret_access_key" json:"secret_access_key"`
	Region          string `mapstructure:"region" yaml:"region" json:"region"`
	SessionToken    string `mapstructure:"session_token" yaml:"session_token" json:"session_token"`
	// Add other AWS-specific fields here
}

// AzureConfig represents the configuration settings specific to Microsoft Azure.
type AzureConfig struct {
	ClientID       string `mapstructure:"client_id" yaml:"client_id" json:"client_id"`
	ClientSecret   string `mapstructure:"client_secret" yaml:"client_secret" json:"client_secret"`
	SubscriptionID string `mapstructure:"subscription_id" yaml:"subscription_id" json:"subscription_id"`
	// Add other Azure-specific fields here
}

type ApiConfig struct {
	Port uint   `mapstructure:"port" yaml:"port" json:"port" validate:"required,lte=65535" `
	Host string `mapstructure:"host" yaml:"host" json:"host" validate:"required,hostname"`
}

type DbConfig struct {
	Username     string `mapstructure:"username" yaml:"username" json:"username" validate:"required"`
	Password     string `mapstructure:"password" yaml:"password" json:"password" validate:"required"`
	DatabaseName string `mapstructure:"database_name" yaml:"database_name" json:"database_name" validate:"required"`
	Host         string `mapstructure:"host" yaml:"host" json:"host" validate:"required,hostname"`
	Port         uint   `mapstructure:"port" yaml:"port" json:"port" validate:"required,lte=65535"`
}

type CognitoConfig struct {
	AppSecret string `mapstructure:"app_secret" yaml:"app_secret" json:"app_secret"`
	ClientId  string `mapstructure:"client_id" yaml:"client_id" json:"client_id"`
}

type IDBConfig interface {
	GetDBConfig() *Config
}

type IApiConfig interface {
	GetApiConfig() *ApiConfig
}

// Config is the main configuration structure for the application. It includes configurations for GCP, AWS, Azure, as well as configurations for the database and API.
type Config struct {
	dirname       string
	filename      string
	GCP           *GCPConfig     `mapstructure:"gcp" yaml:"gcp" json:"GCP"`
	AWS           *AWSConfig     `mapstructure:"aws" yaml:"aws" json:"AWS"`
	Azure         *AzureConfig   `mapstructure:"azure" yaml:"azure" json:"Azure"`
	DBConfig      DbConfig       `mapstructure:"db_config" yaml:"db_config" json:"DBConfig"`
	APIConfig     ApiConfig      `mapstructure:"api_config" yaml:"api_config" json:"APIConfig"`
	CognitoConfig *CognitoConfig `mapstructure:"cognito_config" yaml:"cognito_config" json:"cognito_config"`
}

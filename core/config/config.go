// Package config provides a configuration model and validation for the application configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"gitea/pcp-inariam/inariam/pkgs/log"
)

const (
	ErrInvalidConfig      = "error configuration is not correct"
	ErrOpeningConfigFile  = "error opening configuration file"
	ErrDecodingConfig     = "error decoding configuration file"
	ErrCreatingConfigFile = "error creating configuration file"
	ErrEncodingConfig     = "error encoding configuration file"
	ErrCreatingConfigDir  = "error creating configuration directory"
	ErrClosingConfigFile  = "error closing configuration file"
)

var (
	DefaultConfigDir      = filepath.Join(os.Getenv("HOME"), ".inariam")
	DefaultConfigFilename = "config.yaml"
)
var DefaultConfig = Config{
	filename: DefaultConfigFilename,
	dirname:  DefaultConfigDir,
	DBConfig: DbConfig{
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "inariam",
		Host:         "localhost",
		Port:         5432,
	},
	APIConfig: ApiConfig{
		Host: "localhost",
		Port: 8082,
	},
}

func InitConfig() {
	viper.AddConfigPath(DefaultConfigDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // File format
	viper.AutomaticEnv()        // Read from environment variables as well
}

// New Read and returns the configuration with viper
func New() (*Config, error) {

	// Read the configuration file using viper
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if errors.As(err, &configFileNotFoundError) {
			// If config not found, use default values
			return &DefaultConfig, nil
		} else if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("Config.New: %w", err)
	}

	var config Config
	// Unmarshal the config file into the Config struct
	viper.SetConfigType("yaml")
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("Config.New: %w", err)
	}

	return &config, nil
}

// CreateDefaultConfig constructs a default configuration.
// It determines the appropriate user configuration directory and sets up default directory and filename for the application's config.
func CreateDefaultConfig() (Config, error) {
	userConfigDir, err := os.UserConfigDir()

	if err != nil {
		log.Logger.Panicln("Error finding configuration directory", err)
	}

	appConfigDir := filepath.Join(userConfigDir, "inariam")

	return Config{
		filename: "config.yaml",
		dirname:  appConfigDir,
	}, nil
}

// CreateConfigFileIfNotExist checks if the configuration file exists at the specified path.
// If it doesn't exist, it creates a new one with default values.
func (config *Config) CreateConfigFileIfNotExist() error {
	_, err := os.Stat(config.getConfigPath())
	if os.IsNotExist(err) {
		if err := os.MkdirAll(config.dirname, os.ModePerm); err != nil {
			return fmt.Errorf("Config.CreateConfigFileIfNotExist %s %w", ErrCreatingConfigDir, err)
		}

		configFile, err := os.Create(
			config.getConfigPath(),
		)
		if err != nil {
			return fmt.Errorf("Config.CreateConfigFile: %s %w", ErrClosingConfigFile, err)
		}

		defer func(configFile *os.File) {
			err = configFile.Close()
			if err != nil {
				err = fmt.Errorf("Config.CreateConfigFile: %s %w", ErrClosingConfigFile, err)
			}
		}(configFile)

		err = yaml.NewEncoder(configFile).Encode(config)
		if err != nil {
			return fmt.Errorf("Config.CreateConfigFile: %s %w", ErrCreatingConfigFile, err)
		}

		log.Logger.Infoln("Configuration file created:", config.getConfigPath())

		err = nil

		return err
	}
	return nil
}

// getConfigPath constructs the full path for the configuration file.
func (config *Config) getConfigPath() string {
	return filepath.Join(config.dirname, config.filename)
}

// ParseConfigFile opens the configuration file from the computed path, decodes its contents, and fills in missing configurations.
// After parsing, it validates the configuration to ensure all necessary fields are present and correctly set up.
func (config *Config) ParseConfigFile() error {
	configFile, err := os.Open(config.getConfigPath())
	if err != nil {
		return fmt.Errorf("Config.ParseConfig: %s %w", ErrOpeningConfigFile, err)
	}

	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			err = fmt.Errorf("Config.ParseConfig: %s %w", ErrClosingConfigFile, err)
		}
	}(configFile)

	// Parse YAML data
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return fmt.Errorf("Config.ParseConfig: %s %w", ErrDecodingConfig, err)
	}

	// Check if AWS, GCP, and Azure configurations are missing and fill with nil
	if config.AWS == nil {
		config.AWS = &AWSConfig{}
	}

	if config.GCP == nil {
		config.GCP = &GCPConfig{}
	}

	if config.Azure == nil {
		config.Azure = &AzureConfig{}
	}

	err = config.Validate()

	if err != nil {
		return fmt.Errorf("Config.ParseConfig: %s %w", ErrInvalidConfig, err)
	}

	return err
}

// UpdateConfig encodes the current configuration state into YAML, then encodes that YAML into Base64.
// It then writes the Base64-encoded data back to the configuration file.
func (config *Config) UpdateConfig() error {
	// Encode the updated configuration as YAML

	configFile, err := os.OpenFile(config.getConfigPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("Config.UpdateConfig: %s %w", ErrOpeningConfigFile, err)
	}

	err = yaml.NewEncoder(configFile).Encode(config)
	if err != nil {
		return fmt.Errorf("Config.UpdateConfig: %s %w", ErrEncodingConfig, err)
	}

	return nil
}

// GetDBConfig returns the database configuration from the main configuration.
func (config *Config) GetDBConfig() *DbConfig {
	return &config.DBConfig
}

// GetAPIConfig returns the API configuration from the main configuration.
func (config *Config) GetAPIConfig() *ApiConfig {
	return &config.APIConfig
}

func Check(filename string) error {

	cfg := &Config{
		dirname:  "",
		filename: filename,
	}

	err := cfg.ParseConfigFile()

	if err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil

}

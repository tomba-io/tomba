package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Key          string `mapstructure:"key"`
	Secret       string `mapstructure:"secret"`
	AccessToken  string `mapstructure:"access_token"`
	RefreshToken string `mapstructure:"refresh_token"`
	TokenExpiry  string `mapstructure:"token_expiry"`
	AuthMethod   string `mapstructure:"auth_method"` // "oauth" or "apikey"
}

// InitConfig generates a Viper config file with the provided object key and secret.
func InitConfig() {

	// Set the default configuration values
	viper.SetDefault("key", "")
	viper.SetDefault("secret", "")
	viper.SetDefault("access_token", "")
	viper.SetDefault("refresh_token", "")
	viper.SetDefault("token_expiry", "")
	viper.SetDefault("auth_method", "")
	// Check TOMBA_CONFIG env var for custom config path
	configPath := os.Getenv("TOMBA_CONFIG")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Search config in home directory with name ".tomba" (without extension).
		viper.AddConfigPath(Home())
		viper.SetConfigName(".tomba")
		viper.SetConfigType("json")
	}

	if err := viper.ReadInConfig(); err != nil {
		// Write the file
		writePath := configPath
		if writePath == "" {
			writePath = filepath.Join(Home(), ".tomba.json")
		}
		err := viper.WriteConfigAs(writePath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// UpdateConfig update config file with all credential fields.
func UpdateConfig(c Config) error {
	viper.Set("key", c.Key)
	viper.Set("secret", c.Secret)
	viper.Set("access_token", c.AccessToken)
	viper.Set("refresh_token", c.RefreshToken)
	viper.Set("token_expiry", c.TokenExpiry)
	viper.Set("auth_method", c.AuthMethod)

	viper.AddConfigPath(Home())
	viper.SetConfigName(".tomba")
	viper.SetConfigType("json")

	return viper.WriteConfigAs(Home() + "/.tomba.json")
}

// ReadConfigFile reads the configuration from the file and returns the Config struct.
func ReadConfigFile() (*Config, error) {

	// Set the config file name and type
	viper.SetConfigName(".tomba")
	viper.SetConfigType("json")

	// Set the config file's path (you can modify this according to your needs)
	viper.AddConfigPath(Home())

	// Read the configuration from the file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal the config data into the Config struct
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

// Home Find home directory.
func Home() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

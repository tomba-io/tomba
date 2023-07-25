package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Key    string `mapstructure:"key"`
	Secret string `mapstructure:"secret"`
}

// InitConfig generates a Viper config file with the provided object key and secret.
func InitConfig() {

	// Set the default configuration values
	viper.SetDefault("key", "")
	viper.SetDefault("secret", "")
	// Search config in home directory with name ".email" (without extension).
	viper.AddConfigPath(Home())
	viper.SetConfigName(".email")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		// Write the file
		err := viper.WriteConfigAs(Home() + "/.email.json")
		if err != nil {
			fmt.Println(err)
		}
	}
}

// UpdateConfig update config file new object key and secret.
func UpdateConfig(c Config) error {

	// Set the default configuration values
	viper.Set("key", c.Key)
	viper.Set("secret", c.Secret)
	// Search config in home directory with name ".email" (without extension).
	viper.AddConfigPath(Home())
	viper.SetConfigName(".email")
	viper.SetConfigType("json")

	// Write the file
	err := viper.WriteConfigAs(Home() + "/.email.json")
	if err != nil {
		return err
	}
	return nil
}

// readConfigFile reads the configuration from the file and returns the Config struct.
func ReadConfigFile() (*Config, error) {

	// Set the config file name and type
	viper.SetConfigName(".email")
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

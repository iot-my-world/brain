package config

import (
	"github.com/iot-my-world/brain/environment"
	"github.com/iot-my-world/brain/log"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	MongoNodes                []string
	MongoUser                 string
	MongoPassword             string
	MailRedirectBaseUrl       string
	EmailHost                 string
	EmailAddress              string
	EmailPassword             string
	RootPasswordFileLocation  string
	PathToEmailTemplateFolder string
	KeyFilePath               string
	Environment               environment.Type
}

func New(pathToConfigFile string) Config {
	viper.SetConfigFile(pathToConfigFile)

	viper.SetDefault("mongoNodes", []string{"localhost:27017"})
	viper.SetDefault("mongoUser", "")
	viper.SetDefault("mongoPassword", "")
	viper.SetDefault("mailRedirectBaseUrl", "http://localhost:3000")
	viper.SetDefault("emailHost", "")
	viper.SetDefault("emailAddress", "")
	viper.SetDefault("emailPassword", "")
	viper.SetDefault("rootPasswordFileLocation", "")
	viper.SetDefault("pathToEmailTemplateFolder", "assets/email/template")
	viper.SetDefault("keyFilePath", "")
	viper.SetDefault("environment", environment.Development)

	// check if the config file exists
	if _, err := os.Stat(pathToConfigFile); err != nil {
		// if the error is that the file does not exist
		if os.IsNotExist(err) {
			// create the file
			if err := viper.WriteConfigAs(pathToConfigFile); err != nil {
				log.Fatal("error writing config file", err)
			}
		} else {
			// the error is something else
			log.Fatal("error determining if config file exists", err)
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("error reading in config file", err)
	}

	return Config{
		MongoNodes:                viper.GetStringSlice("mongoNodes"),
		MongoUser:                 viper.GetString("mongoUser"),
		MongoPassword:             viper.GetString("mongoPassword"),
		MailRedirectBaseUrl:       viper.GetString("mailRedirectBaseUrl"),
		EmailHost:                 viper.GetString("emailHost"),
		EmailAddress:              viper.GetString("emailAddress"),
		EmailPassword:             viper.GetString("emailPassword"),
		RootPasswordFileLocation:  viper.GetString("rootPasswordFileLocation"),
		PathToEmailTemplateFolder: viper.GetString("pathToEmailTemplateFolder"),
		KeyFilePath:               viper.GetString("keyFilePath"),
		Environment:               environment.Type(viper.GetString("environment")),
	}
}

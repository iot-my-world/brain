package config

import (
	"github.com/iot-my-world/brain/log"
	"github.com/spf13/viper"
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
}

func New(configFilePath string) Config {
	viper.SetConfigFile(configFilePath)

	viper.SetDefault("mongoNodes", []string{"localhost:27015"})
	viper.SetDefault("mongoUser", "")
	viper.SetDefault("mongoPassword", "")
	viper.SetDefault("mailRedirectBaseUrl", "http://localhost:3000")
	viper.SetDefault("emailHost", "")
	viper.SetDefault("emailAddress", "")
	viper.SetDefault("emailPassword", "")
	viper.SetDefault("rootPasswordFileLocation", "")
	viper.SetDefault("pathToEmailTemplateFolder", "communication/email/template")
	viper.SetDefault("keyFilePath", "")

	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Info("using default brain config")
		default:
			log.Fatal("error reading in config file")
		}
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
	}
}

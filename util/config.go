package util

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type HTTPServerConfig struct {
	Address string `mapstructure:"address"`
}

type MySQLConfig struct {
	DBName   string `mapstructure:"db_name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int32  `mapstructure:"port"`
	Host     string `mapstructure:"host"`
}

// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	FilePath    string `mapstructure:"file_path"`

	HTTPServer HTTPServerConfig `mapstructure:"http_server"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AutomaticEnv()
	configFile := viper.GetString("PORTS_CONFIGFILE")
	if configFile == "" {
		logrus.Fatal("No configuration file found. Pleas set the PORTS_CONFIGFILE environment variable.")
	}
	configFileParts := strings.Split(configFile, ".")

	logrus.WithField("configFile", configFile).Debug("loading configuration file")

	viper.AddConfigPath(path)
	viper.SetConfigName(configFileParts[0])
	viper.SetConfigType(configFileParts[1])

	err = viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot read config")
		return
	}

	err = viper.Unmarshal(&config)
	logrus.WithField("environment", config.Environment).Debug("loaded configuration")
	return
}

func LoadTestConfig(path string, testDBName string) (config Config, err error) {
	config, err = LoadConfig(path)
	if err != nil {
		return
	}

	// Override the MySQL.DBName with the test database name
	config.MySQL.DBName = testDBName
	return
}

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logJSON bool
var Verbose int = 3

func getLogger(config *viper.Viper) logrus.FieldLogger {
	log := logrus.New()

	switch Verbose {
	case 0:
		log.Level = logrus.ErrorLevel
	case 1:
		log.Level = logrus.WarnLevel
	case 2:
		log.Level = logrus.InfoLevel
	case 3:
		log.Level = logrus.DebugLevel
	default:
		log.Level = logrus.DebugLevel
	}

	if logJSON {
		log.Formatter = new(logrus.JSONFormatter)
	}

	appName := config.GetString("metadata.appName")
	version := config.GetString("metadata.version")

	fieldLogger := log.WithFields(logrus.Fields{
		"source":  appName,
		"version": version,
	})
	return fieldLogger
}

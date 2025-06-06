package initialization

import (
	"com.say.more.server/config"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func Logger(cfg *config.Config) *log.Logger {
	logger = log.New()
	//if cfg.App.Env == "prod" {
	//	logger.SetFormatter(&log.JSONFormatter{})
	//} else {
	//	logger.SetFormatter(&log.TextFormatter{
	//		TimestampFormat: "2006-01-02 15:04:05",
	//		FullTimestamp:   true,
	//	})
	//}
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetReportCaller(true)
	switch cfg.App.LogLevel {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
	return logger
}

func GetLogger() *log.Logger {
	return logger
}

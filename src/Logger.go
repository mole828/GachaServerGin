package src

import "github.com/sirupsen/logrus"

func logger() *logrus.Logger {
	re := logrus.New()
	//re.SetReportCaller(true)
	//re.SetFormatter(&logrus.TextFormatter{
	//
	//	//ForceColors: true,
	//	//EnvironmentOverrideColors: true,
	//	TimestampFormat: "2006-01-02 15:04:05",
	//})
	return re
}

var Logger = logger()

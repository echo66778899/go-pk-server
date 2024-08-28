package mylog

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	// Set logrus configuration options here, if needed
	// For example, you can set the log level:
	// log.SetLevel(logrus.DebugLevel)
}

package mylog

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Set logrus configuration options here, if needed
	// For example, you can set the log level:
	log.SetLevel(logrus.DebugLevel)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	black  = "\033[30m"
	white  = "\033[37m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

func TestPrintSuits() {
	// Define colors for suits
	suits := map[string]string{
		"♥": red,
		"♦": yellow,
		"♣": green,
		"♠": white,
	}

	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	for suit, color := range suits {
		for _, rank := range ranks {
			card := []string{
				"┌─────┐",
				fmt.Sprintf("│%-2s   │", rank),
				fmt.Sprintf("│  %s%s%s  │", color, suit, reset),
				"└─────┘",
			}
			for _, line := range card {
				fmt.Println(line)
			}
			fmt.Println() // Print an empty line between cards
		}
	}
}

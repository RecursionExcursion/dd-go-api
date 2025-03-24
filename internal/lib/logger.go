package lib

import (
	"log"
	"strconv"
)

// type Logger struct {
// }

// func NewLogger() *Logger {
// 	return &Logger{}
// }

/* Logger.Log(msg, lvl) will log the msgs if the env logLevel is less than the Logs level
* If log level is < 0, it will always be logged
 */

func Log(msg string, lvl int) {
	if lvl >= getLogLevel() || lvl > 0 {
		log.Println(msg)
	}

}

var getLogLevel = initEnvLogLevel()

func initEnvLogLevel() func() int {
	var logLevelSecret = EnvGet("LOG_LEVEL")
	logLvl, err := strconv.ParseInt(logLevelSecret, 10, 0)
	if err != nil {
		panic("could not convert log level secret to int")
	}
	return func() int {
		return int(logLvl)
	}

}

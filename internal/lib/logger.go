package lib

import (
	"log"
	"strconv"
)

/* Logger.Log(msg, lvl) will log the msgs if the env logLevel is less than the Logs level
* If log level is < 0, it will always be logged
 */

func Log(msg string, lvl int) {
	if lvl >= getLogLevel() || lvl < 0 {
		log.Println(msg)
	}

}

func LogFn(fn func(), lvl int) {
	if lvl >= getLogLevel() || lvl < 0 {
		fn()
	}
}

var getLogLevel = initEnvLogLevel()

func initEnvLogLevel() func() int {
	logLvl, err := strconv.Atoi(EnvGet("LOG_LEVEL"))
	if err != nil {
		panic("could not convert log level secret to int")
	}
	return func() int {
		return logLvl
	}
}

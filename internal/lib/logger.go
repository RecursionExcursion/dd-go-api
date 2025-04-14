package lib

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

/* Logger.Log(msg, lvl) will log the msgs if the env logLevel is less than the Logs level
* If log level is < 0, it will always be logged
 */

func Log(msg string, lvl ...int) {
	if len(lvl) == 0 || lvl[0] >= getLogLevel() || lvl[0] < 0 {
		log.Println(msg)
	}
}

func LogError(msg string, err error, stacks ...string) {
	stackStr := ""

	if stacks != nil {

		if len(stacks) == 1 {
			stackStr = stacks[0]
		} else {
			for i := len(stacks) - 1; i >= 0; i-- {
				s := stacks[i]

				if i == len(stacks)-1 {
					stackStr += s
					continue
				}

				if i == 0 {
					stackStr += fmt.Sprintf("(%v%v", s, strings.Repeat(")", i+1))
					continue
				}

				stackStr += fmt.Sprintf("(%v", s)
			}
		}

		stackStr = fmt.Sprintf(" (%v)", stackStr)
	}

	log.Printf("ERROR[%v]%v: %v", msg, stackStr, err.Error())
}

func LogFn(fn func(), lvl int) {
	if lvl >= getLogLevel() || lvl < 0 {
		fn()
	}
}

var getLogLevel = initEnvLogLevel()

func initEnvLogLevel() func() int {
	logLvl, err := strconv.Atoi(EnvGetOrPanic("LOG_LEVEL"))
	if err != nil {
		panic("could not convert log level secret to int")
	}
	return func() int {
		return logLvl
	}
}

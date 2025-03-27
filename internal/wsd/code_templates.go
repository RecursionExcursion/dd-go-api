package wsd

import (
	"strings"
)

type codeTemplate struct {
	imports    []string
	code       string
	delimmiter string
}

func (ct *codeTemplate) inject(in string) {
	ct.code = strings.Replace(ct.code, ct.delimmiter, in, 1)
}

var mainFuncTemplate = codeTemplate{
	code: `func main(){
	<args>
	}`,
	delimmiter: "<args>",
}

var execFnCallTemplate = codeTemplate{
	code:       `execCommand(<args>)`,
	delimmiter: `<args>`,
}

var execFnTemplate = codeTemplate{
	imports: []string{
		"fmt",
		"os/exec",
		"runtime",
		"strings",
		"log",
	},
	code: `func execCommand(path string) error {

	cmdPrefix := "cmd:"
	urlPrefix := "url:"

	var cmd *exec.Cmd
	if strings.HasPrefix(path, cmdPrefix) {
		parts := strings.SplitN(path, ":", 2)
		if len(parts) != 2 {
			log.Fatal("invalid code command")
		}
		cmdParts := strings.SplitN(parts[1], " ", -1)
		cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
	} else if strings.HasPrefix(path, urlPrefix) {
		parts := strings.SplitN(path, ":", 2)
		if len(parts) != 2 {
			log.Fatal("invalid code command")
		}
		url := parts[1]

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		default:
			return fmt.Errorf("unsupported platform")
		}
	} else {
		log.Fatal("Missing command prefix")
	}

	return cmd.Start()
}`,
}

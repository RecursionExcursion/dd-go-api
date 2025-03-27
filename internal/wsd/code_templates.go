package wsd

import "strings"

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

var exeFnCallTemplate = codeTemplate{
	code:       `execCommand(<args>)`,
	delimmiter: `<args>`,
}

var execFnTemplate = codeTemplate{
	imports: []string{
		"fmt",
		"os/exec",
		"runtime",
	},
	code: `func execCommand(path string) error {
	
			var cmd *exec.Cmd
		
			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
			case "darwin":
				cmd = exec.Command("open", path)
			case "linux":
				cmd = exec.Command("xdg-open", path)
			default:
				return fmt.Errorf("unsupported platform")
			}
		
			return cmd.Start()
		}`,
}

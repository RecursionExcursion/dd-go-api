package wsd

import (
	"fmt"
	"os"
	"strings"
)

type script struct {
	code string
}

func (s *script) addLine(l string) {
	s.code += "\n" + l
}

/* GenerateScript - creates a go script that opens files/urls
 * path- is wher the script is going to be created
 * args- urls/paths to open
 */

func GenerateScript(f *os.File, args ...string) error {
	//Create base script and imports
	fileContent := script{
		code: genPackageStatement(),
	}
	fileContent.addLine(createImportStatement(execFnTemplate.imports...))

	//create main fn logic
	mainFnScript := script{
		code: "",
	}

	//add exefn calls with args
	for _, a := range args {
		execCall := execFnCallTemplate
		execCall.inject(fmt.Sprintf("\"%v\"", a))
		mainFnScript.addLine(execCall.code)
	}

	//inject code into mainfn wrapper
	mainFn := mainFuncTemplate
	mainFn.inject(mainFnScript.code)

	//add main fn and logic to file
	fileContent.addLine(mainFn.code)
	//Add exefn at bottom of file
	fileContent.addLine(execFnTemplate.code)

	_, err := f.Write([]byte(fileContent.code))
	if err != nil {
		return err
	}
	return nil
}

func genPackageStatement() string {
	return "package main"
}

func createImportStatement(args ...string) string {
	switch len(args) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("import \"%v\"", args[0])
	default:
		str := "import (<i>)"
		importStr := ""
		for _, a := range args {
			importStr += fmt.Sprintf("\"%v\"\n", a)
		}
		return strings.Replace(str, "<i>", importStr, 1)
	}
}

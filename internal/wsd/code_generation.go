package app

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

func GenerateScript(path string, args ...string) error {
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
		execCall := exeFnCallTemplate
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

	//write file
	err := os.WriteFile(path, []byte(fileContent.code), os.ModePerm)
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

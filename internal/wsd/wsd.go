package wsd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type CreateExeParams struct {
	Name     string
	Arch     string
	Commands []string
}

// TODO update logs to use lib.Log()
func CreateGoExe(params CreateExeParams) ([]byte, string, error) {

	log.Println(`Creating temp dir`)

	tempDir, f, err := createTempDirAndFile()
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, "", err
	}
	defer os.RemoveAll(tempDir)

	log.Println(`gen script`)
	err = GenerateScript(f, params.Commands...)
	if err != nil {
		return nil, "", err
	}

	log.Println(`building binary`)
	binPath, exeName, err := execCmdOnTempProject(tempDir, params)
	if err != nil {
		return nil, "", err
	}

	lib.Log(fmt.Sprintf("Build successful. Binary at:%v", binPath), 5)

	log.Println(`Reading bin`)
	inMemBin, err := os.ReadFile(binPath)
	if err != nil {
		return nil, "", err
	}

	return inMemBin, exeName, nil
}

func createTempDirAndFile() (string, *os.File, error) {
	tempDir, err := CreateTempDir("", "go-app")
	if err != nil {
		return "", nil, err
	}

	f, err := CreateTempFile(tempDir, "foo")
	if err != nil {
		return tempDir, nil, err
	}
	return tempDir, f, nil
}

func execCmdOnTempProject(tempDir string, params CreateExeParams) (string, string, error) {
	// BUILD go.mod, must be done before the exe bin is created or it will fail
	modInit := createExecCmd(tempDir, "go", "mod", "init", "tmp.com/tmp")
	if out, err := modInit.CombinedOutput(); err != nil {
		fmt.Println("go mod init failed:", string(out))
		return "", "", err
	}

	arc := SupportedArchitecture[params.Arch]

	exeName := ""
	if params.Name != "" {
		exeName = params.Name
	} else {
		exeName = "app"
	}

	if params.Arch == win64 {
		exeName += ".exe"
	}

	binPath := filepath.Join(tempDir, exeName)
	binCmd := createExecCmd(tempDir, "go", "build", "-o", binPath)
	binCmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%v", arc.os),
		fmt.Sprintf("GOARCH=%v", arc.arch),
	)

	output, err := binCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Build failed:", string(output))
		return "", "", err
	}
	return binPath, exeName, nil
}

func createExecCmd(dir string, name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...)
	command.Dir = dir
	return command
}

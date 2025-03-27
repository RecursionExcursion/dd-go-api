package wsd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type CreateExeParams struct {
	Arch     string
	Commands []string
}

// TODO update logs to use lib.Log()
func CreateGoExe(params CreateExeParams) ([]byte, string, error) {

	tempDir, err := CreateTempDir("", "go-app")
	if err != nil {
		return nil, "", err
	}
	defer os.RemoveAll(tempDir)

	f, err := CreateTempFile(tempDir, "foo")
	if err != nil {
		return nil, "", err
	}

	err = GenerateScript(f, params.Commands...)
	if err != nil {
		return nil, "", err
	}

	// BUILD go.mod, must be done before the exe bin is created or it will fail
	modInit := createExecCmd(tempDir, "go", "mod", "init", "tmp.com/tmp")
	if out, err := modInit.CombinedOutput(); err != nil {
		fmt.Println("go mod init failed:", string(out))
		return nil, "", err
	}

	arc := compilationPairs[params.Arch]

	exeName := ""
	if params.Arch == "win" {
		exeName = "app.exe"
	} else {
		exeName = "app"
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
		return nil, "", err
	}

	lib.Log(fmt.Sprintf("Build successful. Binary at:%v", binPath), 5)

	inMemBin, err := os.ReadFile(binPath)
	if err != nil {
		return nil, "", err
	}

	return inMemBin, exeName, nil
}

func createExecCmd(dir string, name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...)
	command.Dir = dir
	return command
}

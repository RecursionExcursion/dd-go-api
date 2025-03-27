package wsd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type CreateExeParams struct {
	Arch     string
	Commands []string
}

func CreateGoExe(params CreateExeParams) ([]byte, string, error) {

	tempDir, err := CreateTempDir("", "go-app")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	f, err := CreateTemp(tempDir, "foo")
	if err != nil {
		panic(err)
	}

	err = GenerateScript(f, params.Commands...)
	if err != nil {
		return nil, "", err
	}

	modInit := createExecCmd(tempDir, "go", "mod", "init", "tmp.com/tmp")
	if out, err := modInit.CombinedOutput(); err != nil {
		fmt.Println("go mod init failed:", string(out))
		panic(err)
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
		panic(err)
	}

	fmt.Println("Build successful. Binary at:", binPath)

	ReadDir(tempDir)

	inMemBin, err := os.ReadFile(binPath)
	if err != nil {
		panic(err)
	}

	return inMemBin, exeName, nil
}

func createExecCmd(dir string, name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...)
	command.Dir = dir
	return command
}

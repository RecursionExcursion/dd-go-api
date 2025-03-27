package wsd

import (
	"fmt"
	"log"
	"os"
)

func ReadDir(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Contents of %s:\n", path)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func CreateTempFile(path string, name string) (*os.File, error) {
	return os.CreateTemp(path, fmt.Sprintf("%v-*.go", name))
}

func CreateTempDir(path string, name string) (string, error) {
	return os.MkdirTemp(path, fmt.Sprintf("%v-*", name))
}

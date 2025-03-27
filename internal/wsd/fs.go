package wsd

import (
	"fmt"
	"log"
	"os"
	// "path/filepath"
)

var folderName = "temp"

// var fileName = "temp.go"
// var path = filepath.Join(folderName, fileName)

// func CreateTempFolder() (string, error) {

// 	_, err := os.Getwd()
// 	if err != nil {
// 		return "", err
// 	}

// 	err = os.M(path, os.ModePerm)
// 	if err != nil {
// 		return "", err
// 	}

// 	return path, nil
// }

func cleanUpTempFolder() error {
	return os.RemoveAll(folderName)
}

func deleteFile(n string) error {
	return os.Remove(n)
}

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

func ReadRoot() {
	dir := "/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Contents of %s:\n", dir)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func CreateTemp(path string, name string) (*os.File, error) {
	return os.CreateTemp(path, fmt.Sprintf("%v-*.go", name))
}

func CreateTempDir(path string, name string) (string, error) {
	return os.MkdirTemp(path, fmt.Sprintf("%v-*", name))
}

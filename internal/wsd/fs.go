package app

import (
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

package lib

import "encoding/base64"

func BytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base64ToBytes(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

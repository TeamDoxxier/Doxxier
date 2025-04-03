package internal

import (
	"encoding/base64"
	"path"
	"path/filepath"
	"runtime"
)

func StringToBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func BytesToBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func Base64ToBytes(input string) ([]byte, error) {
	output, err := base64.StdEncoding.DecodeString(input)
	return output, err
}

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(b))
	return filepath.Dir(root)
}

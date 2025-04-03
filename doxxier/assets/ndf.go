package assets

import (
	_ "embed"
)

//go:embed ndf.json
var ndfData string
var ndf map[string]interface{}

func ParseNdf() (string, error) {
	return ndfData, nil
}

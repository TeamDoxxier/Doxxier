package compression

import (
	"bytes"

	"github.com/ulikunitz/xz"
)

func Compress(data []byte) []byte {
	var buffer bytes.Buffer
	w, _ := xz.NewWriter(&buffer)
	w.Write(data)
	w.Close()
	return buffer.Bytes()
}

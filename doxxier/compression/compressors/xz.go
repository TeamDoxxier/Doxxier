package compression

import (
	"bytes"

	"github.com/ulikunitz/xz"
)

type XzCompressor struct{}

// compress implements compression.Compressor.
func (x *XzCompressor) Compress(data []byte) []byte {
	var buffer bytes.Buffer
	w, _ := xz.NewWriter(&buffer)
	w.Write(data)
	w.Close()
	return buffer.Bytes()
}

// uncompress implements compression.Compressor.
func (x *XzCompressor) Uncompress(data []byte) []byte {
	var buffer bytes.Buffer
	r, _ := xz.NewReader(bytes.NewReader(data))
	buffer.ReadFrom(r)
	return buffer.Bytes()
}

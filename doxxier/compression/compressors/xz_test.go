package compression

import "testing"

func TestCompress(t *testing.T) {
	data := []byte("Hello, World!")
	compressor := &XzCompressor{}
	compressed := compressor.Compress(data)
	if len(compressed) == 0 {
		t.Error("Compressed data is empty")
	}
}

func TestUncompress(t *testing.T) {
	data := []byte("Hello, World!")
	var compressor = &XzCompressor{}
	compressed := compressor.Compress(data)
	uncompressed := compressor.Uncompress(compressed)
	if string(uncompressed) != string(data) {
		t.Error("Uncompressed data does not match original data")
	}
}

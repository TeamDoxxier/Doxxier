package compression

import (
	"testing"
)

func TestCompress(t *testing.T) {
	data := []byte("Hello, World!")
	compressed := Compress(AlgXz, data)
	if len(compressed) == 0 {
		t.Error("Compressed data is empty")
	}
}

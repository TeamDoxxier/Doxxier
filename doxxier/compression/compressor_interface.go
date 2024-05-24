package compression

type Compressor interface {
	Compress(data []byte) []byte
	Uncompress(data []byte) []byte
}

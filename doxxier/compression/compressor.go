package compression

type Compressor interface {
	Compress(data []byte) []byte
	Uncompress(data []byte) []byte
}

type CompressionAlgorithm string

const (
	AlgXz     CompressionAlgorithm = "xz"
	AlgGzip   CompressionAlgorithm = "gzip"
	AlgBrotli CompressionAlgorithm = "brotli"
)

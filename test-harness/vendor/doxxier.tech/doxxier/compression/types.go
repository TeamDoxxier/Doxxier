package compression

type CompressionAlgorithm string

const (
	AlgXz     CompressionAlgorithm = "xz"
	AlgGzip   CompressionAlgorithm = "gzip"
	AlgBrotli CompressionAlgorithm = "brotli"
)

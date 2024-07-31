package compression

type CompressionAlgorithm int

const (
	AlgXz CompressionAlgorithm = iota
	AlgGzip
	AlgBrotli
)

func StringToCompressionAlgorithm(s string) CompressionAlgorithm {
	mapping := map[string]CompressionAlgorithm{
		"AlgXz":     AlgXz,
		"AlgGzip":   AlgGzip,
		"AlgBrotli": AlgBrotli,
	}
	return mapping[s]
}

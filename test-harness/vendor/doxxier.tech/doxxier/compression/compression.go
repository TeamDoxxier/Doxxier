package compression

import compression "doxxier.tech/doxxier/compression/compressors"

func Compress(data []byte, alg CompressionAlgorithm) []byte {

	var compressor Compressor
	switch alg {
	case AlgXz:
		compressor = &compression.XzCompressor{}
	default:
		return data
	}
	return compressor.Compress(data)
}

func Uncompress(data []byte, alg CompressionAlgorithm) []byte {
	var compressor Compressor
	switch alg {
	case AlgXz:
		compressor = &compression.XzCompressor{}
	default:
		return data
	}
	return compressor.Uncompress(data)
}

package compression

import compression "doxxier.tech/doxxier/compression/compressors"

func Compress(alg CompressionAlgorithm, data []byte) []byte {

	var compressor Compressor
	switch alg {
	case AlgXz:
		compressor = &compression.XzCompressor{}
	default:
		return data
	}
	return compressor.Compress(data)
}

func Uncompress(alg CompressionAlgorithm, data []byte) []byte {
	var compressor Compressor
	switch alg {
	case AlgXz:
		compressor = &compression.XzCompressor{}
	default:
		return data
	}
	return compressor.Uncompress(data)
}

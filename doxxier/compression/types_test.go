package compression

import (
	"testing"
)

func TestStringToCompressionAlgorithm(t *testing.T) {
	t.Run("Test StringToCompressionAlgorithm", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			input    string
			expected CompressionAlgorithm
		}{
			{
				input:    "AlgXz",
				expected: AlgXz,
			},
			{
				input:    "AlgGzip",
				expected: AlgGzip,
			},
			{
				input:    "AlgBrotli",
				expected: AlgBrotli,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.input, func(t *testing.T) {
				t.Parallel()
				actual := StringToCompressionAlgorithm(tc.input)
				if actual != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, actual)
				}
			})
		}
	})
}

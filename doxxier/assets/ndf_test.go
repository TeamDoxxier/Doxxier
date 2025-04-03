package assets

import (
	"os"
	"testing"
)

func TestParseNdf(t *testing.T) {
	// Test case: Valid NDF data
	NdfData, err := os.ReadFile("ndf.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := string(NdfData)
	result, err := ParseNdf()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != expected {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

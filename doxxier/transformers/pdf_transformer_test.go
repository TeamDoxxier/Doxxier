package transformers

import (
	"os"
	"testing"

	"doxxier.tech/doxxier/models"
	"github.com/stretchr/testify/assert"
)

func TestPdfTransformer_Transform(t *testing.T) {
	dir := "../test_assets/pdf/"
	pdfTransformer := PdfTransformer{}
	ctx := &models.DoxxierContext{}
	f, err := os.ReadFile(dir + "sample.pdf")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	ctx.Content = f
	err = pdfTransformer.Transform(ctx)
	assert.NoError(t, err)
}

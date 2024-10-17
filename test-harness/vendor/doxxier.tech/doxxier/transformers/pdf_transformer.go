package transformers

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"doxxier.tech/doxxier/pkg/models"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PdfTransformer struct{}

func (t *PdfTransformer) Transform(ctx *models.DoxxierPart) error {
	reader := bytes.NewReader(ctx.Content)
	modelCtx, err := api.ReadContext(reader, nil)

	if err != nil {
		return err
	}
	pages := []string{"1"}
	api.ExtractImages(reader, pages, handleImage, nil)

	md, err := pdfcpu.ExtractMetadata(modelCtx)
	if err != nil {
		return err
	}
	fmt.Println(md)
	return nil
}

func handleImage(img model.Image, b bool, e int) error {
	data, err := io.ReadAll(img.Reader)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("../test_assets/images/image%d.jpg", e), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

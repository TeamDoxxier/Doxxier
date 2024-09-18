package transformers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

<<<<<<< Updated upstream
	"doxxier.tech/doxxier/models"
	"github.com/adrium/goheif"
=======
	"doxxier.tech/doxxier/pkg/models"
>>>>>>> Stashed changes
	"github.com/gen2brain/avif"
	"github.com/nfnt/resize"
)

type ImageTransformer struct{}

func (t *ImageTransformer) Transform(ctx *models.DoxxierPart) error {
	writer := new(bytes.Buffer)
	defer writer.Reset()
	reader := bytes.NewReader(ctx.Content)
	reader.Seek(0, 0)
	img, err := decode(ctx.Content)
	if err != nil {
		return err
	}
	if img.Bounds().Dx() > img.Bounds().Dy() {
		img = resize.Resize(1920, 0, img, resize.Lanczos3)
	} else {
		img = resize.Resize(1080, 0, img, resize.Lanczos3)
	}
	avif.Encode(writer, img, avif.Options{Quality: 20, Speed: 3})
	ctx.Content = writer.Bytes()
	return nil
}

func decode(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	reader.Seek(0, 0)
	mimeType := http.DetectContentType(data)
	switch mimeType {
	case "image/png":
		return png.Decode(reader)
	case "image/jpeg":
		return jpeg.Decode(reader)
	case "application/octet-stream":
		format := string(data[4:12])
		if format == "ftypheic" || format == "ftypheix" || format == "ftyphevc" || format == "ftyphevx" || format == "ftypheis" {
			return goheif.Decode(reader)
		}
	}
	return nil, fmt.Errorf("unsupported image type: %s", mimeType)
}

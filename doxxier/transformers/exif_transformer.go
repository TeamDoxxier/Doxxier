package transformers

import (
	"bytes"
	"fmt"

	"doxxier.tech/doxxier/models"
	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
)

// ExifTransformer is a transformer that extracts EXIF data from an image.
type ExifTransformer struct{}

func (t *ExifTransformer) Transform(ctx *models.DoxxierContext) error {
	reader := bytes.NewReader(ctx.Content)

	e, err := imagemeta.Decode(reader)

	if err == imagemeta.ErrNoExif {
		return nil
	} else if err != nil {
		return fmt.Errorf("Error decoding image: %v", err)
	}

	if e.GPS != (exif2.GPSInfo{}) {
		ctx.Metadata.Gps = models.GpsInfo{
			Latitude:  e.GPS.Latitude(), // Call the function to get the float64 value
			Longitude: e.GPS.Longitude(),
			Date:      e.GPS.Date(),
			Altitude:  e.GPS.Altitude(),
		}
	}
	ctx.Metadata.CreationDateTime = e.CreateDate()
	ctx.Metadata.ModifiedDateTime = e.ModifyDate()
	ctx.Metadata.OriginalDateTime = e.DateTimeOriginal()
	return nil
}

package transformers

import (
	"os"
	"testing"

	"doxxier.tech/doxxier/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestExifTransformer_Transform(t *testing.T) {
	for _, i := range test_images {
		t.Run(i.filename, func(t *testing.T) {
			exifTransformer := ExifTransformer{}
			ctx := &models.DoxxierPart{}
			f, err := os.ReadFile(dir + i.filename)
			if err != nil {
				t.Errorf("Error reading file: %v", err)
			}
			ctx.Content = f
			err = exifTransformer.Transform(ctx)
			if err != nil {
				t.Errorf("Error transforming image: %v", err)
			}
		})
	}
}

func TestExifTransformer_Transform_Has_Gps(t *testing.T) {
	for _, i := range test_images {
		if i.hasExif == true {
			t.Run(i.filename, func(t *testing.T) {
				exifTransformer := ExifTransformer{}
				ctx := &models.DoxxierPart{}
				f, err := os.ReadFile(dir + i.filename)
				if err != nil {
					t.Errorf("Error reading file: %v", err)
				}
				ctx.Content = f
				err = exifTransformer.Transform(ctx)
				assert.True(t, ctx.Metadata.Gps.Latitude != 0)
				assert.True(t, ctx.Metadata.Gps.Longitude != 0)
			})
		}
	}
}

func TestExifTransformer_Transform_Correct_Dates(t *testing.T) {
	exifTransformer := ExifTransformer{}
	ctx := &models.DoxxierPart{}
	f, err := os.ReadFile(dir + "heic-with-metadata.heic")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	ctx.Content = f
	err = exifTransformer.Transform(ctx)
	assert.True(t, ctx.Metadata.OriginalDateTime.Year() == 2018)
	assert.True(t, ctx.Metadata.ModifiedDateTime.Year() == 2018)
	assert.True(t, ctx.Metadata.CreationDateTime.Year() == 2018)
}

func TestExifTransformer_Transform_ErrNoExif(t *testing.T) {
	for _, i := range test_images {
		if i.hasExif == false {
			t.Run(i.filename, func(t *testing.T) {
				exifTransformer := ExifTransformer{}
				ctx := &models.DoxxierPart{}
				f, err := os.ReadFile(dir + i.filename)
				if err != nil {
					t.Errorf("Error reading file: %v", err)
				}
				ctx.Content = f
				err = exifTransformer.Transform(ctx)
				assert.NoError(t, err)
			})
		}
	}
}

var (
	dir         = "../test_assets/images/"
	test_images = []struct {
		filename string
		hasExif  bool
	}{
		{"jpg-with-metadata.jpeg", true},
		{"jpg-without-metadata.jpeg", false},
		{"heic-with-metadata.heic", true},
		{"heic-without-metadata.heic", false},
	}
)

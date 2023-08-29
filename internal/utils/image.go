package utils

import (
	"fmt"
	"image"
	"mime/multipart"

	"github.com/disintegration/imaging"
	"github.com/spf13/viper"
)

type CropData struct {
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

func CropImage(img *multipart.FileHeader, cropData CropData) (image.Image, error) {
	file, err := img.Open()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error saving image")
	}

	defer file.Close()

	i, err := imaging.Decode(file, imaging.AutoOrientation(true))
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error saving image")
	}

	cimg := imaging.Crop(i, image.Rect(
		int(cropData.X),
		int(cropData.Y), int(cropData.Width)+int(cropData.X), int(cropData.Height)+int(cropData.Y)),
	)

	return cimg, nil
}

func CreateImagePath(img string) string {
	var up string

	if img == "" {
		return ""
	}

	if viper.GetString("S3_BUCKET") != "" {
		up = fmt.Sprintf("%s/%s", viper.GetString("CDN_URL"), img)
	} else {
		up = fmt.Sprintf("/uploads/%s", img)
	}
	return up
}

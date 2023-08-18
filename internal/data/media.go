package data

import (
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/oliamb/cutter"
)

type CropData struct {
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

func (d *Data) SaveMediaLocally(img *multipart.FileHeader, cropData CropData) (string, error) {
	var (
		ext       = filepath.Ext(img.Filename)
		randID    = uuid.New().String()
		imagePath = fmt.Sprintf("public/assets/images/%s%s", randID, ext)
		i         image.Image
	)

	file, err := img.Open()
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Error saving image")
	}

	defer file.Close()

	i, err = imaging.Decode(file, imaging.AutoOrientation(true))
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Error saving image")
	}

	cimg, err := cutter.Crop(i, cutter.Config{
		Width:  int(cropData.Width),
		Height: int(cropData.Height),
		Anchor: image.Point{int(cropData.X), int(cropData.Y)},
		Mode:   cutter.TopLeft,
	})

	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Error sa ving image")
	}

	newFile, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	err = jpeg.Encode(newFile, cimg, nil)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

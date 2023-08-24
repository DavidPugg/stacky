package data

import (
	"fmt"
	"image"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type CropData struct {
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

func (d *Data) SaveMediaLocally(img *multipart.FileHeader, cropData CropData) (string, error) {
	var (
		id        = fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(img.Filename))
		imagePath = fmt.Sprintf("uploads/%s", id)
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

	cimg := imaging.Crop(i, image.Rect(
		int(cropData.X),
		int(cropData.Y), int(cropData.Width)+int(cropData.X), int(cropData.Height)+int(cropData.Y)),
	)

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	if err = imaging.Save(cimg, imagePath); err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Error saving image")
	}

	return id, nil
}

func (d *Data) DeleteMediaLocally(id string) error {
	err := os.Remove(fmt.Sprintf("uploads/%s", id))
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Error deleting image")
	}

	return nil
}

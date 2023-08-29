package data

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func (d *Data) SaveMediaLocally(img image.Image, ext string) (string, error) {
	var (
		id        = fmt.Sprintf("%s%s", uuid.New().String(), ext)
		imagePath = fmt.Sprintf("uploads/%s", id)
	)

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	if err := imaging.Save(img, imagePath); err != nil {
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

// func (d *Data) SaveToS3() error {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String("us-east-1")},
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

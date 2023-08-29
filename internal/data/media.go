package data

import (
	"bytes"
	"fmt"
	"image"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (d *Data) SaveMedia(img image.Image, ext string) (string, error) {
	if viper.GetString("S3_BUCKET") != "" {
		return saveToS3(img, ext)
	} else {
		return saveMediaLocally(img, ext)
	}
}

func (d *Data) DeleteMedia(imageID string) error {
	if viper.GetString("S3_BUCKET") != "" {
		return deleteFromS3(imageID)
	} else {
		return deleteMediaLocally(imageID)
	}
}

func saveMediaLocally(img image.Image, ext string) (string, error) {
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

func deleteMediaLocally(id string) error {
	err := os.Remove(fmt.Sprintf("uploads/%s", id))
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Error deleting image")
	}

	return nil
}

func saveToS3(img image.Image, ext string) (string, error) {
	id := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	sess, err := getSessionAWS()
	if err != nil {
		return "", err
	}

	file, err := os.CreateTemp("uploads", fmt.Sprintf("%s%s", "*", ext))
	if err != nil {
		return "", err
	}

	defer os.Remove(file.Name())

	var format imaging.Format
	switch ext {
	case ".png":
		format = imaging.PNG
	case ".jpg":
		format = imaging.JPEG
	case ".jpeg":
		format = imaging.JPEG
	case ".gif":
		format = imaging.GIF
	}

	err = imaging.Encode(file, img, format)
	if err != nil {
		return "", err
	}

	fileContent, err := os.ReadFile(file.Name())
	if err != nil {
		return "", err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	fs := fileInfo.Size()

	svc := s3.New(sess)

	input := &s3.PutObjectInput{
		Body:          bytes.NewReader(fileContent),
		Bucket:        aws.String(viper.GetString("S3_BUCKET")),
		Key:           aws.String(id),
		ContentLength: aws.Int64(fs),
		ContentType:   aws.String(fmt.Sprintf("image/%s", ext[1:])),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		return "", err
	}

	return id, nil
}

func deleteFromS3(imageID string) error {
	sess, err := getSessionAWS()
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetString("S3_BUCKET")),
		Key:    aws.String(imageID),
	}

	_, err = svc.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}

func getSessionAWS() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("S3_REGION")),
		Credentials: credentials.NewStaticCredentials(
			viper.GetString("S3_ACCESS_KEY"),
			viper.GetString("S3_SECRET_KEY"),
			"",
		),
	})
}

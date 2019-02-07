package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sumwhere/models"
)

func ProfileSaver(header *multipart.FileHeader, user *models.User, imageName string) string {

	if header == nil {
		return ""
	}

	file, err := header.Open()
	if err != nil {
		return ""
	}

	defer file.Close()
	path := fmt.Sprintf("/images/%d/profile/", user.Id)
	CreateDirIfNotExist(path)

	dst, err := os.Create(path + fmt.Sprintf("%s.jpg", imageName))
	if err != nil {
		return ""
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return ""
	}
	return fmt.Sprintf("/%d/profile/%s.jpg", user.Id, imageName)
}

func ProfileSaver2(files []*multipart.FileHeader, path string) error {
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		CreateDirIfNotExist(path)

		dst, err := os.Create(path + "profile.jpg")
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}
	return nil
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

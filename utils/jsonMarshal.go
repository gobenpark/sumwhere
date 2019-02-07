package utils

import (
	"io/ioutil"
	"mime/multipart"
)

func JSONConverter(m *multipart.FileHeader) ([]byte, error) {
	file, err := m.Open()

	defer file.Close()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

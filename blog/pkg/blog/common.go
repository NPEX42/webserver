package blog

import "os"

func ReadString(filepath string) (string, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func OpenOrCreate(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
}

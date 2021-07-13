package utils

import "os"

func MakeDirsIfNotExists(path string) error {
	_, err := os.Lstat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

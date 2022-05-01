package worker

import (
	"os"
	"time"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsExist(err)
}

func touch(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create("temp.txt")
		if err != nil {
			return err
		}
		defer file.Close()
		return nil
	}
	currentTime := time.Now().Local()
	err = os.Chtimes(path, currentTime, currentTime)
	if err != nil {
		return err
	}
	return nil
}

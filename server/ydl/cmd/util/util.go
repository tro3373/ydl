package util

import (
	"os"
	"time"

	"github.com/fatih/color"
)

func LogInfo(format string, a ...interface{}) {
	logWrap(color.Green, format, a...)
}

func LogWarn(format string, a ...interface{}) {
	logWrap(color.Yellow, format, a...)
}

func logWrap(fn func(format string, a ...interface{}), format string, a ...interface{}) {
	l := len(a)
	if l == 0 {
		fn(format)
		return
	}
	for i := 0; i < l; i++ {
		format += " %s"
	}
	fn(format, a...)
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// func existsPrefix(name string) (bool, error) {
// 	matches, err := filepath.Glob(name + ".*")
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(matches) > 0, nil
// }
//

func Touch(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
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

func ReadDir(dir string, fn func(dir, name string) error) error {
	if len(dir) == 0 {
		return nil
	}
	if !Exists(dir) {
		return nil
	}
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = fn(dir, name)
		if err != nil {
			return err
		}
	}
	return nil
}

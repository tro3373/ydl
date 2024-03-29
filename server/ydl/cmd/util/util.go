package util

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func LogInfo(format string, a ...interface{}) {
	logWrap(color.Green, format, a...)
}

func LogWarn(format string, a ...interface{}) {
	logWrap(color.Yellow, format, a...)
}

func LogError(format string, a ...interface{}) {
	logWrap(color.Red, format, a...)
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
		file, err := os.Create(filepath.Clean(path))
		if err != nil {
			return err
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Println("Failed to close", path)
			}
		}()
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
	f, err := os.Open(filepath.Clean(dir))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Failed to close", dir)
		}
	}()

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

func ReadFileIfExist(filePath string) (string, error) {
	if len(filePath) == 0 || !Exists(filePath) {
		return "", nil
	}
	raw, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func WriteFile(filePath, data string) error {
	return os.WriteFile(filePath, []byte(data), 0664) //#nosec G306
}

func IsEmptyDir(dirPath string) (bool, error) {
	f, err := os.Open(filepath.Clean(dirPath))
	if err != nil {
		return false, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Failed to close", dirPath)
		}
	}()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

func Contains(list interface{}, target interface{}) bool {
	listValue := reflect.ValueOf(list)
	if listValue.Kind() != reflect.Slice {
		return false
	}
	targetType := reflect.TypeOf(target)
	targetValue := reflect.ValueOf(target)
	for i := 0; i < listValue.Len(); i++ {
		item := listValue.Index(i).Interface()
		itemType := reflect.TypeOf(item)
		if !targetType.ConvertibleTo(itemType) {
			continue
		}
		t := targetValue.Convert(itemType).Interface()
		if ok := reflect.DeepEqual(item, t); ok {
			return true
		}
	}
	return false
}

func CheckSha256sum(targetFile, sha256SumFile, key string) error {
	f, err := os.Open(filepath.Clean(targetFile))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Failed to close", targetFile)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	target := fmt.Sprintf("%x", h.Sum(nil))
	sha256Sum, err := readSha256Sums(sha256SumFile, key)
	if err != nil {
		return err
	}
	fmt.Println("==>    target:", target)
	fmt.Println("==> sha256Sum:", sha256Sum)
	if target == sha256Sum {
		fmt.Println("===> sha256Sum ok")
		return nil
	}
	return errors.Errorf("Invalid sha256sum binary. target:%s, sha256sum:%s", targetFile, sha256SumFile)
}

func readSha256Sums(filePath, key string) (string, error) {
	fp, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	defer func() {
		if err := fp.Close(); err != nil {
			fmt.Println("Failed to close", filePath)
		}
	}()

	r := regexp.MustCompile(".*" + key)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if !r.MatchString(line) {
			continue
		}
		idx := strings.Index(line, " ")
		return line[0:idx], nil
	}
	return "", errors.Errorf("No such sha256Sum exist %s", key)
}

func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return -1, err
	}
	return info.Size(), nil
}

func CreateDirsIfNotExist(dirPaths []string) error {
	for _, dir := range dirPaths {
		if Exists(dir) {
			continue
		}
		//#nosec G301
		if err := os.MkdirAll(dir, 0750); err != nil {
			return err
		}
	}
	return nil
}

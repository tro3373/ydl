package sample

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func download(u string) error {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	cmd := exec.Command("youtube-dl", "--write-thumbnail", "-o", "%(title)s.%(ext)s", u)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	f, err := os.Open(dir)
	if err != nil {
		println("dir", err)
		return err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		println("readdir", err)
		return err
	}

	thumbnail := ""
	filename := ""
	for _, fn := range names {
		switch filepath.Ext(fn) {
		case ".jpg", ".png":
			if thumbnail == "" {
				thumbnail = filepath.Join(dir, fn)
			}
		default:
			if filename == "" {
				filename = filepath.Join(dir, fn)
			}
		}
	}

	ext := filepath.Ext(filename)
	name := filepath.Base(filename[:len(filename)-len(ext)])
	output := filepath.Join(dir, name) + ".mp3"

	args := []string{
		`ffmpeg`,
		`-i`,
		filename,
	}

	if thumbnail != "" {
		args = append(args,
			`-i`,
			thumbnail,
			`-map`,
			`0:a`,
			`-map`,
			`1:v`,
		)
	}

	args = append(args,
		`-metadata`, fmt.Sprintf(`title=%s`, name),
		`-metadata`, fmt.Sprintf(`artist=%s`, name),
		`-metadata`, fmt.Sprintf(`album=%s`, name),
		output)

	cmd = exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return os.Rename(output, filepath.Base(output))
}

func main() {
	err := download(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

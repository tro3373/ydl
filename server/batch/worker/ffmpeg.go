package worker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func executeFfmpeg(task *Task) error {

	fmt.Println("=============================================================")
	fmt.Println("=> Start executeFfmpeg", task.String())
	fmt.Println("=============================================================")

	// ctx := task.Ctx
	req := task.Req

	if task.HasAudio() {
		err := os.Remove(task.PathAudio)
		if err != nil {
			return err
		}
	}
	var args []string
	args = append(args, "-i", task.PathMovie)
	if len(task.PathThumbnail) > 0 {
		args = append(args, "-i", task.PathThumbnail, "-map", "0:a", "-map", "1:v")
	}
	args = appendMetaDataIfPresent(args, "title", req.Tag.Title)
	args = appendMetaDataIfPresent(args, "artist", req.Tag.Artist)
	album := req.Tag.Album
	if len(album) == 0 {
		album = req.Tag.Title
	}
	args = appendMetaDataIfPresent(args, "album", album)
	genre := req.Tag.Genre
	f := "Favorite artist of "
	fw := f + "West"
	fj := f + "Japan"
	if len(genre) == 0 {
		genre = fw
	}
	switch genre {
	case "ja", "jap", "Japan", "Ja":
		genre = fj
	case "en", "En", "West":
		genre = fw
	}
	args = appendMetaDataIfPresent(args, "genre", genre)

	args = append(args, task.PathAudio)

	fmt.Println(append([]string{"==> Executing: ffmpeg"}, args...))
	cmd := exec.Command("ffmpeg", args...)
	cmd.Dir = task.DoingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "Failed to execute ffmpeg %s", err)
	}
	return nil
}

func appendMetaDataIfPresent(args []string, key, val string) []string {
	if len(val) == 0 {
		return args
	}
	args = append(args, "-metadata")
	return append(args, fmt.Sprintf("%s='%s'", key, val))
}

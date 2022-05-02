package worker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func executeFfmpeg(task Task) error {

	// ctx := task.Ctx
	req := task.Req

	if task.HasAudio() {
		err := os.Remove(task.FileNameAudio)
		if err != nil {
			return err
		}
	}
	var args []string
	args = append(args, "-i", task.FileNameMovie)
	if len(task.Thumbnail) > 0 {
		args = append(args, "-i", task.Thumbnail, "-map", "0:a", "-map", "1:v")
	}
	args = appendIfPresent(args, "title", req.Tag.Title)
	args = appendIfPresent(args, "artist", req.Tag.Artist)
	args = appendIfPresent(args, "album", req.Tag.Album)
	args = appendIfPresent(args, "genre", req.Tag.Genre)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Dir = task.DstDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "Failed to execute ffmpeg %s", err)
	}
	return nil
}

func appendIfPresent(args []string, key, val string) []string {
	if len(val) == 0 {
		return args
	}
	return append(args, fmt.Sprintf("%s=%s", key, val))
}

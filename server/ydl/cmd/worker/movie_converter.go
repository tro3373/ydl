package worker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/tro3373/ydl/cmd/util"
)

func StartConvert(task *Task) error {

	util.LogInfo("Converting mp3 via ffmpeg..")

	if len(task.TaskPath.Audio) == 0 {
		task.setPathAudioFromPathMovie()
	}
	tag := task.Tag
	taskPath := task.TaskPath

	if task.HasAudio() {
		err := os.Remove(taskPath.Audio)
		if err != nil {
			return err
		}
	}
	var args []string
	args = append(args, "-i", taskPath.Movie)
	if len(taskPath.Thumbnail) > 0 {
		args = append(args, "-i", taskPath.Thumbnail, "-map", "0:a", "-map", "1:v")
	}
	args = appendMetaDataIfPresent(args, "title", tag.Title)
	args = appendMetaDataIfPresent(args, "artist", tag.Artist)
	album := tag.Album
	if len(album) == 0 {
		album = tag.Title
	}
	args = appendMetaDataIfPresent(args, "album", album)
	genre := getGenre(tag.Genre)
	args = appendMetaDataIfPresent(args, "genre", genre)

	args = append(args, taskPath.Audio)

	util.LogInfo("==> Executing: ffmpeg", args)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Dir = taskPath.Doing
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
	return append(args, fmt.Sprintf("%s=%s", key, val))
}

func getGenre(genre string) string {
	// genre := tag.Genre
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
	return genre
}

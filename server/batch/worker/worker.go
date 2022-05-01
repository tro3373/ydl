package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/tro3373/ydl/batch/request"
)

func Start(ctx Ctx, event fsnotify.Event) {
	if isTaskRunning(ctx) {
		fmt.Println("=> Already tasks running")
		return
	}
	touchTaskRunning(ctx)
	msgs := []interface{}{"Success to execute tasks!"}
	defer func() {
		rmTakRunning(ctx)
		fmt.Println(msgs...)
	}()
	err := startTasks(ctx)
	if err != nil {
		msgs = []interface{}{"Failed to execute tasks..", err}
	}
}

func startTasks(ctx Ctx) error {
	err := updateYoutubeDlIfNeeded(ctx)
	if err != nil {
		return err
	}
	// queued := filepath.Join(workDir, "queue")
	jsons, err := findJsons(ctx.QueueDir)
	if err != nil {
		return err
	}
	for _, json := range jsons {
		err := startDownload(ctx, json)
		if err != nil {
			return err
		}
	}
	return nil
}

func findJsons(dir string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, "*", "*.json"))
}

func startDownload(ctx Ctx, jsonPath string) error {
	raw, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read json %s", jsonPath)
	}
	var req request.Exec
	json.Unmarshal(raw, &req)
	return executeYoutubeDl(ctx, req)
}

// func exists(filename string) bool {
// 	_, err := os.Stat(filename)
// 	return os.IsExist(err)
// }
// func existsPrefix(name string) (bool, error) {
// 	matches, err := filepath.Glob(name + ".*")
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(matches) > 0, nil
// }
//

func executeYoutubeDl(ctx Ctx, req request.Exec) error {
	var args []string
	key := req.Key()
	args = append(args, key)
	args = append(args, "-o")
	format := "%(id)s_%(title)s.%(ext)s"
	doned := ctx.DoneDir
	libd := ctx.LibDir
	args = append(args, filepath.Join(doned, format))

	// for audio output
	args = append(args, "-x")
	args = append(args, "--audio-format")
	args = append(args, "mp3")

	err := exec.Command(filepath.Join(libd, ctx.YoutubeDl), args...).Run()
	if err != nil {
		log.Fatalf("Failed to executeYoutubeDl %s: %v", key, err)
		return err
	}
	return nil
}

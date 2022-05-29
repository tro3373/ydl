package worker

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tro3373/ydl/cmd/util"
)

func Start(ctx Ctx, event fsnotify.Event) {
	if isTaskRunning(ctx) {
		util.LogWarn("=> Already tasks running")
		return
	}
	touchTaskRunning(ctx)
	msgs := []interface{}{"Success to execute tasks!"}
	defer func() {
		rmTaskRunning(ctx)
		fmt.Println(msgs...)
	}()
	err := startTasks(ctx)
	if err != nil {
		msgs = []interface{}{"Failed to execute tasks..", err}
	}
}

func startTasks(ctx Ctx) error {
	jsons, err := findJsons(ctx.QueueDir)
	if err != nil {
		return err
	}
	for _, json := range jsons {
		err := handleJson(ctx, json)
		if err != nil {
			return err
		}
	}
	return nil
}

func findJsons(dir string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, "*.json"))
}

func handleJson(ctx Ctx, jsonPath string) error {
	task, err := NewTask(ctx, jsonPath)
	if err != nil {
		return err
	}
	util.LogInfo("=> Starting New Task!", task.String())

	if !task.HasMovie() {
		err = startDownloadMovie(task)
		if err != nil {
			return err
		}
	}

	err = startConvert(task)
	if err != nil {
		return err
	}

	return task.Done()
}

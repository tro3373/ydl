package worker

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
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

	if !task.HasMovie() {
		err = executeYoutubeDl(task)
		if err != nil {
			return err
		}
	}

	err = executeFfmpeg(task)
	if err != nil {
		return err
	}

	return nil
}

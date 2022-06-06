package worker

import (
	"os"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tro3373/ydl/cmd/util"
)

func Start(ctx Ctx, event fsnotify.Event) {
	if isTaskRunning(ctx) {
		util.LogWarn("=> Already tasks running")
		return
	}
	err := touchTaskRunning(ctx)
	if err != nil {
		return
	}

	msg := "=> Success to execute all tasks!"
	defer func() {
		rmTaskRunning(ctx)
		util.LogInfo(msg)
	}()
	err = startTasks(ctx)
	if err != nil {
		msg = "=> Some task was failed.."
	}
}

func getTasksRunningFile(ctx Ctx) string {
	return path.Join(ctx.WorkDir, ".tasks_running")
}

func isTaskRunning(ctx Ctx) bool {
	return util.Exists(getTasksRunningFile(ctx))
}

func touchTaskRunning(ctx Ctx) error {
	file := getTasksRunningFile(ctx)
	if err := util.Touch(file); err != nil {
		util.LogError("Failed to touch", file, err)
		return err
	}
	return nil
}

func rmTaskRunning(ctx Ctx) error {
	file := getTasksRunningFile(ctx)
	if err := os.Remove(file); err != nil {
		util.LogError("Failed to remove", file, err)
		return err
	}
	return nil
}

func cleanTaskRunning(ctx Ctx) error {
	if !isTaskRunning(ctx) {
		return nil
	}
	return rmTaskRunning(ctx)
}

func startTasks(ctx Ctx) error {
	jsons, err := findJsons(ctx.WorkDirs.Queue)
	if err != nil {
		util.LogError("Failed to findJson", err)
		return err
	}
	for _, json := range jsons {
		handleJsonErr := handleJson(ctx, json)
		if handleJsonErr != nil {
			util.LogError("Failed to handle json ", json, handleJsonErr)
			err = handleJsonErr
		}
	}
	return err
}

func findJsons(dir string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, "*.json"))
}

func handleJson(ctx Ctx, jsonPath string) error {
	task, err := NewTask(ctx, jsonPath)
	if err != nil {
		return err
	}
	util.LogInfo("=> Start New Task!", task.String())

	if !task.HasMovie() {
		err = StartDownloadMovie(task)
		if err != nil {
			return err
		}
	}

	err = StartConvert(task)
	if err != nil {
		return err
	}

	err = task.Done()
	if err != nil {
		util.LogError("=> Task failed..", task.String())
		return err
	}
	util.LogInfo("=> Task Done!", task.String())
	return err
}

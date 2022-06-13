package worker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tro3373/ydl/cmd/util"
	"github.com/tro3373/ydl/cmd/worker/ctx"
	"github.com/tro3373/ydl/cmd/worker/lib"
	"github.com/tro3373/ydl/cmd/worker/movie"
	"github.com/tro3373/ydl/cmd/worker/task"
)

func Start(ctx ctx.Ctx) {
	var handleError = func(message string, err error) {
		if err == nil {
			return
		}
		fmt.Println(message, err)
		os.Exit(1)
	}

	handleError("Failed to clean running file", cleanTaskRunningIfNeeded(ctx))
	handleError("Failed to update lib", lib.UpdateLibIFNeeded(ctx))
	watcher, err := fsnotify.NewWatcher()
	handleError("Failed to new watcher", err)

	defer watcher.Close()
	done := make(chan bool)
	go dog(watcher, ctx)

	fmt.Println("==> Watching", ctx.WorkDirs.Queue, "..")
	if err := watcher.Add(ctx.WorkDirs.Queue); err != nil {
		fmt.Println("Failed to watcher.Add", err)
	}
	<-done
}

func dog(watcher *fsnotify.Watcher, ctx ctx.Ctx) {
	for {
		select {
		case event := <-watcher.Events:
			// Receive event! fsnotify.Event{Name:"path/to/the/file", Op:0x1}
			fmt.Printf("Receive event! Name:%s Op:%s\n", event.Name, event.Op.String())
			if event.Op&fsnotify.Create == fsnotify.Create {
				StartTasks(ctx, event)
			}
		case err := <-watcher.Errors:
			fmt.Println("Receive error!", err)
		}
	}
}

func StartTasks(ctx ctx.Ctx, event fsnotify.Event) {
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
	err = handleTasks(ctx)
	if err != nil {
		msg = "=> Some task was failed.."
	}
}

func handleTasks(ctx ctx.Ctx) error {
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

func handleJson(ctx ctx.Ctx, jsonPath string) error {
	task, err := task.NewTask(ctx, jsonPath)
	if err != nil {
		return err
	}
	util.LogInfo("=> Start New Task!", task.String())

	if !task.HasMovie() {
		err = movie.StartDownload(task)
		if err != nil {
			return err
		}
	}

	err = movie.StartConvert(task)
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

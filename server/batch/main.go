package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tro3373/ydl/batch/worker"
)

func main() {
	var workDir string
	if len(os.Args) > 1 {
		workDir = os.Args[1]
	}
	ctx := initializeDir(workDir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	fmt.Println("==> Watching", ctx.QueueDir, "..")

	defer watcher.Close()
	done := make(chan bool)
	go dog(watcher, ctx)
	if err := watcher.Add(ctx.QueueDir); err != nil {
		fmt.Println("Failed to watcher.Add", err)
	}
	<-done
}

func initializeDir(workRootDir string) worker.Ctx {
	if len(workRootDir) == 0 {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		workRootDir = filepath.Join(filepath.Dir(ex), "work")
	}
	lib := createDirIfNotExist(workRootDir, "lib")
	queue := createDirIfNotExist(workRootDir, "queue")
	doing := createDirIfNotExist(workRootDir, "doing")
	done := createDirIfNotExist(workRootDir, "done")

	return worker.NewCtx(workRootDir, lib, queue, doing, done)
}

func createDirIfNotExist(dstRootDir, targetDir string) string {
	dir := filepath.Join(dstRootDir, targetDir)
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0775)
	}
	return dir
}

func dog(watcher *fsnotify.Watcher, ctx worker.Ctx) {
	for {
		select {
		case event := <-watcher.Events:
			// Receive event! fsnotify.Event{Name:"path/to/the/file", Op:0x1}
			fmt.Printf("Receive event! Name:%s Op:%s\n", event.Name, event.Op.String())
			if event.Op&fsnotify.Create == fsnotify.Create {
				worker.Start(ctx, event)
			}
		case err := <-watcher.Errors:
			fmt.Println("Receive error!", err)
		}
	}
}

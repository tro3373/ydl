package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tro3373/ydl/cmd/worker"
)

func StartBatch(ctx worker.Ctx) {
	err := worker.UpdateLibIFNeeded(ctx)
	if err != nil {
		fmt.Println("Update lib Error", err)
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("New Watcher Error", err)
		os.Exit(1)
	}

	defer watcher.Close()
	done := make(chan bool)
	go dog(watcher, ctx)

	fmt.Println("==> Watching", ctx.WorkDirs.Queue, "..")
	if err := watcher.Add(ctx.WorkDirs.Queue); err != nil {
		fmt.Println("Failed to watcher.Add", err)
	}
	<-done
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

package worker

import (
	"fmt"
	"os"
	"path"

	"github.com/tro3373/ydl/cmd/util"
)

func getTasksRunningFile(ctx Ctx) string {
	return path.Join(ctx.WorkDir, ".tasks_running")
}

func isTaskRunning(ctx Ctx) bool {
	return util.Exists(getTasksRunningFile(ctx))
}

func touchTaskRunning(ctx Ctx) {
	file := getTasksRunningFile(ctx)
	if err := util.Touch(file); err != nil {
		fmt.Println("Failed to touch", file, err)
	}
}

func rmTaskRunning(ctx Ctx) {
	file := getTasksRunningFile(ctx)
	if err := os.Remove(file); err != nil {
		fmt.Println("Failed to remove", file, err)
	}
}

func cleanTaskRunning(ctx Ctx) {
	if !isTaskRunning(ctx) {
		return
	}
	rmTaskRunning(ctx)
}

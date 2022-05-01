package worker

import (
	"fmt"
	"os"
	"path"
)

func getTasksRunningFile(ctx Ctx) string {
	return path.Join(ctx.WorkDir, ".tasks_running")
}

func isTaskRunning(ctx Ctx) bool {
	return exists(getTasksRunningFile(ctx))
}

func touchTaskRunning(ctx Ctx) {
	file := getTasksRunningFile(ctx)
	if err := touch(file); err != nil {
		fmt.Println("Failed to touch", file, err)
	}
}

func rmTakRunning(ctx Ctx) {
	file := getTasksRunningFile(ctx)
	if err := os.Remove(file); err != nil {
		fmt.Println("Failed to remove", file, err)
	}
}

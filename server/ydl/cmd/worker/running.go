package worker

import (
	"os"
	"path"

	"github.com/tro3373/ydl/cmd/util"
	"github.com/tro3373/ydl/cmd/worker/ctx"
)

func getTasksRunningFile(ctx ctx.Ctx) string {
	return path.Join(ctx.WorkDir, ".tasks_running")
}

func isTaskRunning(ctx ctx.Ctx) bool {
	return util.Exists(getTasksRunningFile(ctx))
}

func touchTaskRunning(ctx ctx.Ctx) error {
	file := getTasksRunningFile(ctx)
	if err := util.Touch(file); err != nil {
		util.LogError("Failed to touch", file, err)
		return err
	}
	return nil
}

func rmTaskRunning(ctx ctx.Ctx) error {
	file := getTasksRunningFile(ctx)
	if err := os.Remove(file); err != nil {
		util.LogError("Failed to remove", file, err)
		return err
	}
	return nil
}

func cleanTaskRunningIfNeeded(ctx ctx.Ctx) error {
	if !isTaskRunning(ctx) {
		return nil
	}
	return rmTaskRunning(ctx)
}

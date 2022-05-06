package worker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tro3373/ydl/cmd/util"
)

type Ctx struct {
	WorkDir   string
	LibDir    string
	QueueDir  string
	DoingDir  string
	DoneDir   string
	YoutubeDl string
}

func NewCtx(workDir string) (Ctx, error) {
	if len(workDir) == 0 {
		ex, err := os.Executable()
		if err != nil {
			return Ctx{}, err
		}
		workDir = filepath.Join(filepath.Dir(ex), "work")
	}
	ctx := Ctx{
		WorkDir:   createDirIfNotExist(workDir, ""),
		LibDir:    createDirIfNotExist(workDir, "lib"),
		QueueDir:  createDirIfNotExist(workDir, "queue"),
		DoingDir:  createDirIfNotExist(workDir, "doing"),
		DoneDir:   createDirIfNotExist(workDir, "done"),
		YoutubeDl: "youtube-dl",
	}
	err := ctx.Clean()
	return ctx, err
}

func createDirIfNotExist(targetDirPath, subDir string) string {
	dir := targetDirPath
	if len(subDir) > 0 {
		dir = filepath.Join(targetDirPath, subDir)
	}
	if !util.Exists(dir) {
		os.MkdirAll(dir, 0775)
	}
	return dir
}

func (ctx Ctx) Clean() error {
	fmt.Println("==> Cleaning", ctx.DoingDir)
	// TODO duplicate with NewCtx doingDir
	err := os.RemoveAll(ctx.DoingDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(ctx.DoingDir, 0775)
	if err != nil {
		return err
	}
	return nil
}

func (ctx Ctx) GetDoneDir(key string) string {
	return filepath.Join(ctx.DoneDir, key)
}

func (ctx Ctx) GetDoingDir(key string) string {
	return filepath.Join(ctx.DoingDir, key)
}

package worker

import (
	"os"
	"path/filepath"

	"github.com/tro3373/ydl/cmd/util"
)

type Ctx struct {
	WorkDir         string `json:"workDir"`
	LibDir          string `json:"libDir"`
	QueueDir        string `json:"queueDir"`
	DoingDir        string `json:"doingDir"`
	DoneDir         string `json:"doneDir"`
	DownloadLibRepo string `json:"DownloadLibRepo"`
	DownloadLibName string `json:"DownloadLibName"`
}

func NewCtx(args []string) (Ctx, error) {
	workDir, err := chooseWorkDir(args)
	if err != nil {
		return Ctx{}, err
	}
	ctx := Ctx{
		WorkDir:         createDirIfNotExist(workDir, ""),
		LibDir:          createDirIfNotExist(workDir, "lib"),
		QueueDir:        createDirIfNotExist(workDir, "queue"),
		DoingDir:        createDirIfNotExist(workDir, "doing"),
		DoneDir:         createDirIfNotExist(workDir, "done"),
		DownloadLibRepo: "ytdl-org/youtube-dl",
		DownloadLibName: "youtube-dl",
	}
	err = ctx.Clean()
	return ctx, err
}

func chooseWorkDir(args []string) (string, error) {
	var dir string
	if len(args) > 0 {
		dir = args[0]
	}
	if len(dir) == 0 {
		currentDir, err := os.Getwd()
		if err != nil {
			return currentDir, err
		}
		dir = filepath.Join(currentDir, "work")
	}
	return filepath.Abs(dir)
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
	cleanTaskRunning(ctx)

	// fmt.Println("==> Cleaning", ctx.DoingDir)
	// // TODO duplicate with NewCtx doingDir
	// err := os.RemoveAll(ctx.DoingDir)
	// if err != nil {
	// 	return err
	// }
	// err = os.MkdirAll(ctx.DoingDir, 0775)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (ctx Ctx) GetDoneDir(key string) string {
	return filepath.Join(ctx.DoneDir, key)
}

func (ctx Ctx) GetDoingDir(key string) string {
	return filepath.Join(ctx.DoingDir, key)
}

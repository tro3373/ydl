package worker

import (
	"os"
	"path/filepath"

	"github.com/tro3373/ydl/cmd/util"
)

type Ctx struct {
	WorkDir     string      `json:"work"`
	WorkDirs    WorkDirs    `json:"workDirs"`
	DownloadLib DownloadLib `json:"downloadLib"`
}

type WorkDirs struct {
	Lib   string `json:"lib"`
	Queue string `json:"queue"`
	Doing string `json:"doing"`
	Done  string `json:"done"`
}

type DownloadLib struct {
	Repo string `json:"repo"`
	Name string `json:"name"`
	Sums string `json:"sums"`
}

func NewCtx(args []string) (Ctx, error) {
	workDir, err := chooseWorkDir(args)
	if err != nil {
		return Ctx{}, err
	}
	workDirs := WorkDirs{
		Lib:   createDirIfNotExist(workDir, "lib"),
		Queue: createDirIfNotExist(workDir, "queue"),
		Doing: createDirIfNotExist(workDir, "doing"),
		Done:  createDirIfNotExist(workDir, "done"),
	}
	downloadLib := DownloadLib{
		// Repo: "ytdl-org/youtube-dl",
		// Name: "youtube-dl",
		Repo: "yt-dlp/yt-dlp",
		Name: "yt-dlp",
		Sums: "SHA2-256SUMS",
	}
	ctx := Ctx{
		WorkDir:     createDirIfNotExist(workDir, ""),
		WorkDirs:    workDirs,
		DownloadLib: downloadLib,
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
	return cleanTaskRunning(ctx)
}

func (ctx Ctx) GetDoneDir(key string) string {
	return filepath.Join(ctx.WorkDirs.Done, key)
}

func (ctx Ctx) GetDoingDir(key string) string {
	return filepath.Join(ctx.WorkDirs.Doing, key)
}

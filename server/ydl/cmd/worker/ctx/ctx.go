package ctx

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
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

	lib := filepath.Join(workDir, "lib")
	queue := filepath.Join(workDir, "queue")
	doing := filepath.Join(workDir, "doing")
	done := filepath.Join(workDir, "done")

	dirs := []string{lib, queue, doing, done}
	if err := util.CreateDirsIfNotExist(dirs); err != nil {
		return Ctx{}, errors.Wrapf(err, "Failed to create dirs %v", dirs)
	}

	workDirs := WorkDirs{
		Lib:   lib,
		Queue: queue,
		Doing: doing,
		Done:  done,
	}
	downloadLib := DownloadLib{
		// Repo: "ytdl-org/youtube-dl",
		// Name: "youtube-dl",
		Repo: "yt-dlp/yt-dlp",
		Name: "yt-dlp",
		Sums: "SHA2-256SUMS",
	}
	ctx := Ctx{
		WorkDir:     workDir,
		WorkDirs:    workDirs,
		DownloadLib: downloadLib,
	}
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

func (ctx Ctx) GetDoneDir(key string) string {
	return filepath.Join(ctx.WorkDirs.Done, key)
}

func (ctx Ctx) GetDoingDir(key string) string {
	return filepath.Join(ctx.WorkDirs.Doing, key)
}

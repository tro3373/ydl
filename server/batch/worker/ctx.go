package worker

import (
	"os"
	"path/filepath"
)

type Ctx struct {
	WorkDir   string
	LibDir    string
	QueueDir  string
	DoingDir  string
	DoneDir   string
	YoutubeDl string
}

func NewCtx(work, lib, queue, doing, done string) Ctx {
	return Ctx{
		WorkDir:   work,
		LibDir:    lib,
		QueueDir:  queue,
		DoingDir:  doing,
		DoneDir:   done,
		YoutubeDl: "youtube-dl",
	}
}

func (ctx Ctx) Clean() error {
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

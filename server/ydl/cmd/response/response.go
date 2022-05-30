package response

import (
	"github.com/tro3373/ydl/cmd/request"
	"github.com/tro3373/ydl/cmd/worker"
)

type Res struct {
	Req       *request.Req `json:"req"`
	Thumbnail string       `json:"thumbnail"`
	Movie     string       `json:"movie"`
	Audio     string       `json:"audio"`
}

func NewRes(task worker.Task) Res {
	doneDir := task.PathDoneDir
	thumbnail := relPath(doneDir, task.PathThumbnail)
	movie := relPath(doneDir, task.PathMovie)
	audio := relPath(doneDir, task.PathAudio)
	res := Res{
		Req:       task.Req,
		Thumbnail: thumbnail,
		Movie:     movie,
		Audio:     audio,
	}
	return res
}

func relPath(dirPath, filePath string) string {
	return filePath[len(dirPath)+1:]
}

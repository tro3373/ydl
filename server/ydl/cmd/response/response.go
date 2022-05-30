package response

import (
	"fmt"

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
	workDir := task.Ctx.WorkDir
	thumbnail := toResourcePath(workDir, task.PathThumbnail)
	movie := toResourcePath(workDir, task.PathMovie)
	audio := toResourcePath(workDir, task.PathAudio)
	res := Res{
		Req:       task.Req,
		Thumbnail: thumbnail,
		Movie:     movie,
		Audio:     audio,
	}
	return res
}

func toResourcePath(dirPath, filePath string) string {
	return fmt.Sprintf("/resource/%s", filePath[len(dirPath)+1:])
}

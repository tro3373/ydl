package response

import (
	"fmt"

	"github.com/tro3373/ydl/cmd/api/request"
	"github.com/tro3373/ydl/cmd/util"
	"github.com/tro3373/ydl/cmd/worker/task"
)

type Res struct {
	Url       string      `json:"url" binding:"required"`
	Doing     bool        `json:"doing"`
	Tag       request.Tag `json:"tag"`
	CreatedAt string      `json:"createdAt"`
	Thumbnail string      `json:"thumbnail"`
	Movie     string      `json:"movie"`
	Audio     string      `json:"audio"`
	MovieSize int64       `json:"movieSize"`
	AudioSize int64       `json:"audioSize"`
}

func NewRes(task task.Task, doing bool) Res {
	workDir := task.Ctx.WorkDir
	thumbnail := toResourcePath(workDir, task.TaskPath.Thumbnail)
	movie := toResourcePath(workDir, task.TaskPath.Movie)
	audio := toResourcePath(workDir, task.TaskPath.Audio)
	movieSize, _ := util.GetFileSize(task.TaskPath.Movie)
	audioSize, _ := util.GetFileSize(task.TaskPath.Audio)
	res := Res{
		Url:       task.Url,
		Doing:     doing,
		Tag:       task.Tag,
		CreatedAt: task.CreatedAt,
		Thumbnail: thumbnail,
		Movie:     movie,
		Audio:     audio,
		MovieSize: movieSize,
		AudioSize: audioSize,
	}
	return res
}

func getFileSize() {
}

func toResourcePath(dirPath, filePath string) string {
	if len(filePath) == 0 {
		return ""
	}
	return fmt.Sprintf("/resource/%s", filePath[len(dirPath)+1:])
}

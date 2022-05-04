package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tro3373/ydl/batch/request"
)

type Task struct {
	Ctx      Ctx
	Req      *request.Exec
	DoingDir string
	DoneDir  string

	PathReq       string
	PathThumbnail string
	PathMovie     string
	PathAudio     string
}

func NewTask(ctx Ctx, jsonPath string) (Task, error) {
	task := Task{
		Ctx: ctx,
	}
	req, err := task.readJson(jsonPath)
	if err != nil {
		return task, err
	}
	task.Req = req
	key := req.Key()

	task.DoingDir = ctx.GetDoingDir(key)
	if !exists(task.DoingDir) {
		os.MkdirAll(task.DoingDir, 0775)
	}
	task.DoneDir = ctx.GetDoneDir(key)
	err = task.findTargetFile(task.DoneDir)
	if err != nil {
		return task, err
	}
	task.PathReq = filepath.Join(task.DoingDir, "req.json")
	err = os.Rename(jsonPath, task.PathReq)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (task Task) String() string {
	return fmt.Sprintf("%#+v", task)
}

func (task Task) Key() string {
	return task.Req.Key()
}

func (task Task) readJson(jsonPath string) (*request.Exec, error) {
	raw, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read json %s", jsonPath)
	}
	var req request.Exec
	json.Unmarshal(raw, &req)
	return &req, nil
}

func (task Task) findTargetFile(targetDir string) error {
	if !exists(targetDir) {
		return nil
	}
	f, err := os.Open(targetDir)
	if err != nil {
		return err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		fullPath := filepath.Join(targetDir, name)
		switch filepath.Ext(name) {
		case ".json":
			continue
		case ".jpg", ".png", ".webp":
			task.PathThumbnail = fullPath
		case ".mp3":
			task.PathAudio = fullPath
		default:
			if len(task.PathMovie) > 0 {
				continue
			}
			task.PathMovie = fullPath
			task.setPathAudioFromPathMovieIfNeeded()
		}
	}
	return nil
}

func (task Task) setPathAudioFromPathMovieIfNeeded() {
	if len(task.PathAudio) > 0 {
		return
	}
	movie := task.PathMovie
	dir := filepath.Dir(movie)
	ext := filepath.Ext(movie)
	name := filepath.Base(movie[:len(movie)-len(ext)])
	task.PathAudio = filepath.Join(dir, name) + ".mp3"
}

func (task Task) HasMovie() bool {
	return len(task.PathMovie) > 0
}
func (task Task) HasAudio() bool {
	return len(task.PathAudio) > 0
}

func (task Task) Done() error {
	err := task.findTargetFile(task.DoingDir)
	if err != nil {
		return err
	}
	task.RenameDoing2Done(task.PathReq)
	task.RenameDoing2Done(task.PathThumbnail)
	task.RenameDoing2Done(task.PathMovie)
	task.RenameDoing2Done(task.PathAudio)
	return task.Clean()
}

func (task Task) RenameDoing2Done(src string) error {
	if len(src) == 0 {
		return nil
	}
	r := regexp.MustCompile("doing")
	if !r.MatchString(src) {
		return nil
	}
	dst := strings.Replace(src, "doing", "done", -1)
	return os.Rename(src, dst)
}

func (task Task) Clean() error {
	return os.RemoveAll(task.DoingDir)
}

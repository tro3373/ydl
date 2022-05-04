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

func NewTask(ctx Ctx, jsonPath string) (*Task, error) {
	task := Task{
		Ctx: ctx,
	}
	req, err := task.readJson(jsonPath)
	if err != nil {
		return &task, err
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
		return &task, err
	}
	task.PathReq = filepath.Join(task.DoingDir, "req.json")
	err = os.Rename(jsonPath, task.PathReq)
	if err != nil {
		return &task, err
	}
	return &task, nil
}

func (task *Task) String() string {
	return fmt.Sprintf("%#+v", task)
}

func (task *Task) Key() string {
	return task.Req.Key()
}

func (task *Task) readJson(jsonPath string) (*request.Exec, error) {
	raw, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read json %s", jsonPath)
	}
	var req request.Exec
	json.Unmarshal(raw, &req)
	return &req, nil
}

func (task *Task) findTargetFile(targetDir string) error {
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
		// fmt.Println("==> findTargetFile: handling.. ", fullPath)
		switch filepath.Ext(name) {
		case ".json":
			continue
		case ".jpg", ".png", ".webp":
			// fmt.Println("==> findTargetFile: set jpg. ")
			task.PathThumbnail = fullPath
		case ".mp3":
			task.PathAudio = fullPath
		default:
			// fmt.Println("==> findTargetFile: set movie. ", task.PathMovie)
			if len(task.PathMovie) > 0 {
				continue
			}
			task.PathMovie = fullPath
			task.setPathAudioFromPathMovieIfNeeded()
		}
	}
	return nil
}

func (task *Task) setPathAudioFromPathMovieIfNeeded() {
	if len(task.PathAudio) > 0 {
		return
	}
	movie := task.PathMovie
	dir := filepath.Dir(movie)
	ext := filepath.Ext(movie)
	name := filepath.Base(movie[:len(movie)-len(ext)])
	task.PathAudio = filepath.Join(dir, name) + ".mp3"
}

func (task *Task) HasMovie() bool {
	return len(task.PathMovie) > 0 && exists(task.PathMovie)
}
func (task *Task) HasAudio() bool {
	return len(task.PathAudio) > 0 && exists(task.PathAudio)
}

func (task *Task) Done() error {
	if !exists(task.DoneDir) {
		return os.Rename(task.DoingDir, task.DoneDir)
	}
	err := task.findTargetFile(task.DoingDir)
	if err != nil {
		return err
	}
	err = task.RenameDoing2Done(task.PathReq)
	if err != nil {
		return err
	}
	err = task.RenameDoing2Done(task.PathThumbnail)
	if err != nil {
		return err
	}
	err = task.RenameDoing2Done(task.PathMovie)
	if err != nil {
		return err
	}
	return task.RenameDoing2Done(task.PathAudio)
}

func (task *Task) RenameDoing2Done(src string) error {
	if len(src) == 0 {
		return nil
	}
	r := regexp.MustCompile("doing")
	if !r.MatchString(src) {
		return nil
	}
	dst := strings.Replace(src, "doing", "done", -1)
	fmt.Println("===> Renaming from:", src, "to:", dst)
	return os.Rename(src, dst)
}

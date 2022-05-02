package worker

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tro3373/ydl/batch/request"
)

type Task struct {
	Ctx           Ctx
	QueueJsonPath string
	DoneJsonPath  string
	Req           *request.Exec
	DstDir        string

	Thumbnail     string
	FileNameMovie string
	FileNameAudio string
}

func NewTask(ctx Ctx, queueJsonPath string) (Task, error) {
	task := Task{
		Ctx:           ctx,
		QueueJsonPath: queueJsonPath,
	}
	req, err := task.readJson()
	if err != nil {
		return task, err
	}
	task.Req = req
	key := req.Key()
	task.DstDir = ctx.DestDir(key)
	if !exists(task.DstDir) {
		os.MkdirAll(task.DstDir, 0775)
	}
	err = task.findFile()
	if err != nil {
		return task, err
	}
	task.DoneJsonPath = filepath.Join(task.DstDir, "req.json")
	return task, nil
}

func (task Task) Key() string {
	return task.Req.Key()
}

func (task Task) readJson() (*request.Exec, error) {
	raw, err := ioutil.ReadFile(task.QueueJsonPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read json %s", task.QueueJsonPath)
	}
	var req request.Exec
	json.Unmarshal(raw, &req)
	return &req, nil
}

func (task Task) findFile() error {
	f, err := os.Open(task.DstDir)
	if err != nil {
		return err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		switch filepath.Ext(name) {
		case ".json":
			continue
		case ".jpg", ".png":
			task.Thumbnail = name
		case ".mp3":
			task.FileNameAudio = name
		default:
			task.FileNameMovie = name
			task.setAudioFileNameFromMovieIfNeeded()
		}
	}
	return nil
}

func (task Task) setAudioFileNameFromMovieIfNeeded() {
	if len(task.FileNameAudio) > 0 {
		return
	}
	movie := task.FileNameMovie
	dir := filepath.Dir(movie)
	ext := filepath.Ext(movie)
	name := filepath.Base(movie[:len(movie)-len(ext)])
	task.FileNameAudio = filepath.Join(dir, name) + ".mp3"
}

func (task Task) HasMovie() bool {
	return len(task.FileNameMovie) > 0
}
func (task Task) HasAudio() bool {
	return len(task.FileNameAudio) > 0
}

func (task Task) Done() error {
	dstFile := task.DoneJsonPath

	data, _ := json.MarshalIndent(task.Req, "", " ")
	err := ioutil.WriteFile(dstFile, data, 0644)
	if err != nil {
		return err
	}
	return os.Remove(task.QueueJsonPath)
}

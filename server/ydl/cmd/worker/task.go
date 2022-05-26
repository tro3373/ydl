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
	"github.com/tro3373/ydl/cmd/request"
	"github.com/tro3373/ydl/cmd/util"
)

type Task struct {
	Ctx           Ctx           `json:"ctx"`
	Req           *request.Exec `json:"req"`
	PathDoingDir  string        `json:"pathDoingDir"`
	PathDoneDir   string        `json:"pathDoneDir"`
	PathReqJson   string        `json:"pathReqJson"`
	PathInfoJson  string        `json:"pathInfoJson"`
	PathThumbnail string        `json:"pathThumbnail"`
	PathMovie     string        `json:"pathMovie"`
	PathAudio     string        `json:"pathAudio"`
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

	task.PathDoneDir = ctx.GetDoneDir(key)
	findTargetDir := task.PathDoneDir

	task.PathDoingDir = ctx.GetDoingDir(key)
	if util.Exists(task.PathDoingDir) {
		findTargetDir = task.PathDoneDir
	} else {
		os.MkdirAll(task.PathDoingDir, 0775)
	}
	err = task.findTargetFile(findTargetDir)
	if err != nil {
		return &task, err
	}
	task.PathReqJson = filepath.Join(task.PathDoingDir, "req.json")
	err = os.Rename(jsonPath, task.PathReqJson)
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
	err := util.ReadDir(targetDir, task.readDirHandler)
	task.genTitleFromInfoIfEnable()
	return err
}

func (task *Task) readDirHandler(dir, name string) error {
	fullPath := filepath.Join(dir, name)

	fmt.Println("==> readDirHandler: handling.. ", fullPath)
	switch filepath.Ext(name) {
	case ".json":
		if name == "req.json" {
			task.PathReqJson = fullPath
		} else if strings.HasSuffix(name, "info.json") {
			task.PathInfoJson = fullPath
		}
		return nil
	case ".jpg", ".png", ".webp":
		fmt.Println("==> readDirHandler: set jpg. ")
		task.PathThumbnail = fullPath
	case ".mp3":
		task.PathAudio = fullPath
	default:
		fmt.Println("==> readDirHandler: set movie. ", task.PathMovie)
		if len(task.PathMovie) > 0 {
			return nil
		}
		task.PathMovie = fullPath
		task.setPathAudioFromPathMovieIfNeeded()
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

func (task *Task) genTitleFromInfoIfEnable() {
	if len(task.PathInfoJson) == 0 {
		return
	}
	dir := filepath.Dir(task.PathInfoJson)
	matches, _ := filepath.Glob(filepath.Join(dir, "*.title"))
	if len(matches) > 0 {
		return
	}
	raw, err := ioutil.ReadFile(task.PathInfoJson)
	if err != nil {
		fmt.Println("Failed to read info json.", task.PathInfoJson, err)
		return
	}
	var info interface{}
	json.Unmarshal(raw, &info)
	m := info.(map[string]interface{})
	title := m["title"].(string)
	util.Touch(filepath.Join(dir, title+".title"))
}

func (task *Task) HasMovie() bool {
	return len(task.PathMovie) > 0 && util.Exists(task.PathMovie)
}
func (task *Task) HasAudio() bool {
	return len(task.PathAudio) > 0 && util.Exists(task.PathAudio)
}

func (task *Task) Done() error {
	task.save()
	if !util.Exists(task.PathDoneDir) {
		return os.Rename(task.PathDoingDir, task.PathDoneDir)
	}
	return util.ReadDir(task.PathDoingDir, task.RenameDoing2DoneHandler)
}

func (task *Task) save() error {
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}
	dst := filepath.Join(task.PathDoingDir, "task.json")
	return ioutil.WriteFile(dst, data, 0644)
}

func (task *Task) RenameDoing2DoneHandler(dir, name string) error {
	return task.RenameDoing2Done(filepath.Join(dir, name))
}

func (task *Task) RenameDoing2Done(src string) error {
	// logWarn("Renaming", src, "...")
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

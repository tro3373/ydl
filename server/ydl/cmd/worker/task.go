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
	Ctx       Ctx         `json:"ctx"`
	TaskPath  TaskPath    `json:"taskPath"`
	Url       string      `json:"url" binding:"required"`
	Tag       request.Tag `json:"tag"`
	Uuids     []string    `json:"uuids"`
	CreatedAt string      `json:"createdAt"`
}

type TaskPath struct {
	Doing     string `json:"doing"`
	Done      string `json:"done"`
	ReqJson   string `json:"reqJson"`
	InfoJson  string `json:"infoJson"`
	Thumbnail string `json:"thumbnail"`
	Movie     string `json:"movie"`
	Audio     string `json:"audio"`
}

func NewTask(ctx Ctx, jsonPath string) (*Task, error) {
	return newTaskInner(ctx, jsonPath, true)
}

func ReadTask(ctx Ctx, jsonPath string) (*Task, error) {
	return newTaskInner(ctx, jsonPath, false)
}

func newTaskInner(ctx Ctx, jsonPath string, forQueue bool) (*Task, error) {
	task := Task{
		Ctx: ctx,
	}
	req, err := task.readJson(jsonPath)
	if err != nil {
		return &task, err
	}
	task.Uuids = append(task.Uuids, req.Uuid)
	task.CreatedAt = req.CreatedAt
	task.Url = req.Url
	task.Tag = req.Tag

	key := request.Key(req.Url)

	task.TaskPath.Done = ctx.GetDoneDir(key)
	findTargetDir := task.TaskPath.Done

	doingDir := ctx.GetDoingDir(key)
	task.TaskPath.Doing = doingDir
	if util.Exists(doingDir) {
		findTargetDir = task.TaskPath.Done
	} else {
		if forQueue {
			os.MkdirAll(doingDir, 0775)
		}
	}
	err = task.findTargetFile(findTargetDir)
	if err != nil {
		return &task, err
	}
	task.TaskPath.ReqJson = filepath.Join(doingDir, "req.json")
	if forQueue {
		err = os.Rename(jsonPath, task.TaskPath.ReqJson)
		if err != nil {
			return &task, err
		}
	}
	return &task, nil
}

func (task *Task) String() string {
	return fmt.Sprintf("%#+v", task)
}

func (task *Task) Key() string {
	return request.Key(task.Url)
}

func (task *Task) readJson(jsonPath string) (*request.Req, error) {
	raw, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read json %s", jsonPath)
	}
	var req request.Req
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
	switch name {
	case "req.json":
		task.TaskPath.ReqJson = fullPath
	case "src.info.json":
		task.TaskPath.InfoJson = fullPath
	case "src.jpg", "src.png", "src.webp":
		fmt.Println("==> readDirHandler: set jpg. ")
		task.TaskPath.Thumbnail = fullPath
	case "src.mp3":
		task.TaskPath.Audio = fullPath
	default:
		if strings.HasPrefix(name, "src") {
			// fmt.Println("==> readDirHandler: set movie. ", task.TaskPath.Movie)
			// if len(task.TaskPath.Movie) > 0 {
			// 	return nil
			// }
			task.TaskPath.Movie = fullPath
		}
	}

	return nil
}

func (task *Task) setPathAudioFromPathMovie() {
	movie := task.TaskPath.Movie
	dir := filepath.Dir(movie)
	ext := filepath.Ext(movie)
	name := filepath.Base(movie[:len(movie)-len(ext)])
	task.TaskPath.Audio = filepath.Join(dir, name) + ".mp3"
}

func (task *Task) genTitleFromInfoIfEnable() {
	if len(task.TaskPath.InfoJson) == 0 {
		return
	}
	dir := filepath.Dir(task.TaskPath.InfoJson)
	matches, _ := filepath.Glob(filepath.Join(dir, "*.title"))
	if len(matches) > 0 {
		return
	}
	raw, err := ioutil.ReadFile(task.TaskPath.InfoJson)
	if err != nil {
		fmt.Println("Failed to read info json.", task.TaskPath.InfoJson, err)
		return
	}
	var info interface{}
	json.Unmarshal(raw, &info)
	m := info.(map[string]interface{})
	title := m["title"].(string)
	util.Touch(filepath.Join(dir, title+".title"))
}

func (task *Task) HasMovie() bool {
	return len(task.TaskPath.Movie) > 0 && util.Exists(task.TaskPath.Movie)
}
func (task *Task) HasAudio() bool {
	return len(task.TaskPath.Audio) > 0 && util.Exists(task.TaskPath.Audio)
}

func (task *Task) Done() error {
	doneDir := task.TaskPath.Done
	if !util.Exists(doneDir) {
		err := os.MkdirAll(doneDir, 0775)
		if err != nil {
			return err
		}
	}
	doingDir := task.TaskPath.Doing
	err := util.ReadDir(doingDir, task.RenameDoing2DoneHandler)
	if err != nil {
		return err
	}

	err = task.findTargetFile(task.TaskPath.Done)
	if err != nil {
		return err
	}
	err = task.save()
	if err != nil {
		return err
	}

	empty, err := util.IsEmptyDir(doingDir)
	if err != nil {
		return err
	}
	if !empty {
		util.LogWarn("==> Task done but not empty. %s", doingDir)
		return nil
	}
	return os.Remove(doingDir)
}

func (task *Task) save() error {
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}
	dst := filepath.Join(task.TaskPath.Done, "task.json")
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

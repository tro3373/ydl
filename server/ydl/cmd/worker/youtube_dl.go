package worker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tro3373/ydl/cmd/util"
)

func executeYoutubeDl(task *Task) error {

	// fmt.Println("=============================================================")
	// fmt.Println("=> Start executeYoutubeDl", task.String())
	// fmt.Println("=============================================================")
	ctx := task.Ctx
	req := task.Req
	util.LogInfo("=> Downloading via", ctx.YoutubeDl, "..")

	// dir, err := ioutil.TempDir("", "")
	// if err != nil {
	// 	return err
	// }
	// defer os.RemoveAll(dir)

	var args []string
	key := req.Key()

	args = append(args, "--write-thumbnail")
	args = append(args, "--write-info-json")
	// args = append(args, "--write-description")
	args = append(args, "-o")
	// format := "%(id)s_%(title)s.%(ext)s"
	format := "%(title)s.%(ext)s"
	// format := "src.%(ext)s"
	dstd := filepath.Join(ctx.DoingDir, key)
	args = append(args, filepath.Join(dstd, format))

	// // for audio output
	// args = append(args, "-x")
	// args = append(args, "--audio-format")
	// args = append(args, "mp3")

	args = append(args, key)

	fmt.Println(append([]string{"==> Executing: ", filepath.Join(ctx.LibDir, ctx.YoutubeDl)}, args...))
	cmd := exec.Command(filepath.Join(ctx.LibDir, ctx.YoutubeDl), args...)
	cmd.Dir = dstd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "Failed to executeYoutubeDl %s", err)
	}
	err = task.findTargetFile(dstd)
	if err != nil {
		return err
	}
	return nil
}
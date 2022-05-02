package worker

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func executeYoutubeDl(task Task) error {

	ctx := task.Ctx
	req := task.Req

	// dir, err := ioutil.TempDir("", "")
	// if err != nil {
	// 	return err
	// }
	// defer os.RemoveAll(dir)

	var args []string
	key := req.Key()

	args = append(args, "--write-thumbnail")
	args = append(args, "-o")
	format := "%(id)s_%(title)s.%(ext)s"
	dstd := filepath.Join(ctx.DoneDir, key)
	args = append(args, filepath.Join(dstd, format))

	// // for audio output
	// args = append(args, "-x")
	// args = append(args, "--audio-format")
	// args = append(args, "mp3")

	args = append(args, key)

	cmd := exec.Command(filepath.Join(ctx.LibDir, ctx.YoutubeDl), args...)
	cmd.Dir = dstd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "Failed to executeYoutubeDl %s", err)
	}
	err = task.findFile()
	if err != nil {
		return err
	}
	return nil
}

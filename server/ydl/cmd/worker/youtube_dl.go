package worker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tro3373/ydl/cmd/util"
)

func startYoutubeDl(task *Task) error {

	util.LogInfo("=> Downloading via youtube-dl..")
	cmd, dstd := buildCmd(task)
	fmt.Println("==> Executing: ", cmd.String())

	err := runWithRetry(cmd, dstd)
	if err != nil {
		return errors.Wrapf(err, "Failed to executeYoutubeDl %s", err)
	}
	err = task.findTargetFile(dstd)
	if err != nil {
		return err
	}
	return nil
}

func buildCmd(task *Task) (*exec.Cmd, string) {
	ctx := task.Ctx
	req := task.Req
	key := req.Key()

	var args []string
	args = append(args, "--write-thumbnail")
	args = append(args, "--write-info-json")
	// args = append(args, "--write-description")
	args = append(args, "-o")
	// format := "%(id)s_%(title)s.%(ext)s"
	// format := "%(title)s.%(ext)s"
	format := "src.%(ext)s"
	dstd := filepath.Join(ctx.DoingDir, key)
	args = append(args, filepath.Join(dstd, format))

	// // for audio output
	// args = append(args, "-x")
	// args = append(args, "--audio-format")
	// args = append(args, "mp3")

	args = append(args, key)

	absYoutubeDl := filepath.Join(ctx.LibDir, ctx.YoutubeDl)
	cmd := exec.Command(absYoutubeDl, args...)
	cmd.Dir = dstd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, dstd
}

func runWithRetry(cmd *exec.Cmd, dstd string) error {
	fmt.Println("==> Executing: ", cmd.String())
	max := 3
	var err error
	for i := 0; i < max; i++ {
		err = cmd.Run()
		if err != nil {
			fmt.Println("Failed to executeYoutubeDl", err)
			// youtube-dl retry default 10 times(-R option mean)
			// but retry 3times if error occur.
			continue
		}
		// // success but..
		// files, err := filepath.Glob(filepath.Join(dstd, "*.part.*"))
		break
	}
	return err
}

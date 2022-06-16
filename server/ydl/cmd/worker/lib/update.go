package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tro3373/ydl/cmd/worker/ctx"
)

// @see [Go client to github to get latest release and assets for a given repository](https://gist.github.com/metal3d/002e4f0d8545f83c2ace)
func UpdateLibIFNeeded(ctx ctx.Ctx) error {
	fmt.Println("==> Start checking", ctx.DownloadLib.Name, "version..")

	repo := ctx.DownloadLib.Repo
	libd := ctx.WorkDirs.Lib
	libName := ctx.DownloadLib.Name
	sums := ctx.DownloadLib.Sums
	err := DownloadIfNeeded(repo, libd, libName, sums)
	if err != nil {
		return err
	}
	libdAbs, err := filepath.Abs(filepath.Join(libd, libName))
	if err != nil {
		return err
	}
	return os.Chmod(libdAbs, 0775) //#nosec G302
}

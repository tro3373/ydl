package worker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tro3373/ydl/cmd/gh"
)

// @see [Go client to github to get latest release and assets for a given repository](https://gist.github.com/metal3d/002e4f0d8545f83c2ace)
func UpdateLibIFNeeded(ctx Ctx) error {
	fmt.Println("==> Start checking", ctx.DownloadLib.Name, "version..")

	repo := ctx.DownloadLib.Repo
	libd := ctx.WorkDirs.Lib
	libName := ctx.DownloadLib.Name
	sums := ctx.DownloadLib.Sums
	err := gh.DownloadIfNeeded(repo, libd, libName, sums)
	if err != nil {
		return err
	}
	libdAbs, err := filepath.Abs(filepath.Join(libd, libName))
	if err != nil {
		return err
	}
	return os.Chmod(libdAbs, 0775)
}

// func getGhReleaseLatestTagInfo(repo string) (map[string]interface{}, error) {
// 	command := "releases/latest"
// 	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", repo, command)
// 	req, err := prepareRequest(url)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "Failed to prepare request")
// 	}
// 	// Add required headers
// 	req.Header.Add("Accept", "application/vnd.github.v3.text-match+json")
// 	req.Header.Add("Accept", "application/vnd.github.moondragon+json")
//
// 	client := http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if resp.StatusCode < 200 || resp.StatusCode > 299 {
// 		return nil, errors.Errorf("Error: StatusCode: %d, Status: %s", resp.StatusCode, resp.Status)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	bodyText, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "Failed to reading response")
// 	}
// 	result := make(map[string]interface{})
// 	if err := json.Unmarshal(bodyText, &result); err != nil {
// 		return nil, errors.Wrap(err, "Failed to json.Unmarshal")
// 	}
// 	return result, nil
// }
//
// func prepareRequest(url string) (*http.Request, error) {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("User-Agent", "ydl-batch-go-client")
// 	return req, nil
// }
//
// func downloadResource(dstd, repo string, id float64, c chan int) error {
// 	defer func() { c <- 1 }()
// 	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/assets/%.0f", repo, id)
// 	fmt.Printf("Start download: %s\n", url)
// 	req, err := prepareRequest(url)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Add("Accept", "application/octet-stream")
//
// 	client := http.Client{}
// 	resp, _ := client.Do(req)
//
// 	disp := resp.Header.Get("Content-disposition")
// 	dst, err := getDistPathFromDisposition(dstd, disp)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to %s", "parse disposition")
// 	}
// 	err = save(dst, resp.Body)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to save resposence %s", dst)
// 	}
// 	fmt.Printf("Finished download: %s -> %s\n", url, disp)
// 	return nil
// }
//
// func save(dst string, body io.ReadCloser) error {
// 	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to open file %s", dst)
// 	}
// 	defer f.Close()
// 	b := make([]byte, 4096)
// 	var i int
// 	for err == nil {
// 		i, err = body.Read(b)
// 		f.Write(b[:i])
// 	}
// 	return nil
// }
//
// func getDistPathFromDisposition(dstd, disposition string) (string, error) {
// 	re := regexp.MustCompile(`filename=(.+)`)
// 	matches := re.FindAllStringSubmatch(disposition, -1)
// 	if len(matches) == 0 || len(matches[0]) == 0 {
// 		return "", errors.New("Failed to retireave content-disposition")
// 	}
// 	dstf := matches[0][1]
// 	dst := path.Join(dstd, dstf)
// 	return dst, nil
// }

package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

// @see [Go client to github to get latest release and assets for a given repository](https://gist.github.com/metal3d/002e4f0d8545f83c2ace)
func updateYoutubeDlIfNeeded(ctx Ctx) error {
	fmt.Println("==> Start checking", ctx.YoutubeDl, "version..")
	repo := fmt.Sprintf("ytdl-org/%s", ctx.YoutubeDl)
	result, err := getGhReleaseLatestTagInfo(repo)

	tag := result["tag_name"].(string)
	if err != nil || len(tag) == 0 {
		return err
	}
	// fmt.Printf("tag: %s, result:%#+v\n", tag, result)
	fmt.Printf("==> Detect %s tag: %s\n", repo, tag)

	libd := ctx.LibDir
	dstd := filepath.Join(libd, tag)
	if exists(dstd) {
		fmt.Printf("===> Already %s exist.\n", dstd)
		return nil
	}
	os.MkdirAll(dstd, 0775)

	fmt.Println("==> Start download", ctx.YoutubeDl)

	results := make([]interface{}, 0)
	for _, asset := range result["assets"].([]interface{}) {
		a := asset.(map[string]interface{})
		name := a["name"]
		if name != ctx.YoutubeDl && name != fmt.Sprintf("%s.sig", ctx.YoutubeDl) {
			continue
		}
		results = append(results, a["id"])
	}
	c := make(chan int)
	for _, res := range results {
		go downloadResource(dstd, repo, res.(float64), c)
	}
	// TODO sig check
	for i := 0; i < len(results); i++ {
		<-c
	}

	dst := filepath.Join(libd, ctx.YoutubeDl)
	if exists(dst) {
		err := os.Remove(dst)
		if err != nil {
			return errors.Wrapf(err, "Failed to Remove %s", dst)
		}
	}
	libdAbs, err := filepath.Abs(libd)
	if err != nil {
		return errors.Wrapf(err, "Failed to get abs path for %s", libd)
	}
	src := filepath.Join(libdAbs, tag, ctx.YoutubeDl)
	err = os.Chmod(src, 0775)
	if err != nil {
		return errors.Wrapf(err, "Failed to chmod %s", src)
	}
	err = os.Symlink(src, dst)
	if err != nil {
		return errors.Wrapf(err, "Failed to create symlink %s", src)
	}
	return nil
}

func getGhReleaseLatestTagInfo(repo string) (map[string]interface{}, error) {
	command := "releases/latest"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", repo, command)
	req, err := prepareRequest(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to prepare request")
	}
	// Add required headers
	req.Header.Add("Accept", "application/vnd.github.v3.text-match+json")
	req.Header.Add("Accept", "application/vnd.github.moondragon+json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.Errorf("Error: StatusCode: %d, Status: %s", resp.StatusCode, resp.Status)
	}
	if err != nil {
		return nil, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to reading response")
	}
	result := make(map[string]interface{})
	if err := json.Unmarshal(bodyText, &result); err != nil {
		return nil, errors.Wrap(err, "Failed to json.Unmarshal")
	}
	return result, nil
}

func prepareRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "ydl-batch-go-client")
	return req, nil
}

func downloadResource(dstd, repo string, id float64, c chan int) error {
	defer func() { c <- 1 }()
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/assets/%.0f", repo, id)
	fmt.Printf("Start download: %s\n", url)
	req, err := prepareRequest(url)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/octet-stream")

	client := http.Client{}
	resp, _ := client.Do(req)

	disp := resp.Header.Get("Content-disposition")
	dst, err := getDistPathFromDisposition(dstd, disp)
	if err != nil {
		return errors.Wrapf(err, "Failed to %s", "parse disposition")
	}
	err = save(dst, resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Failed to save resposence %s", dst)
	}
	fmt.Printf("Finished download: %s -> %s\n", url, disp)
	return nil
}

func save(dst string, body io.ReadCloser) error {
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return errors.Wrapf(err, "Failed to open file %s", dst)
	}
	defer f.Close()
	b := make([]byte, 4096)
	var i int
	for err == nil {
		i, err = body.Read(b)
		f.Write(b[:i])
	}
	return nil
}

func getDistPathFromDisposition(dstd, disposition string) (string, error) {
	re := regexp.MustCompile(`filename=(.+)`)
	matches := re.FindAllStringSubmatch(disposition, -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		return "", errors.New("Failed to retireave content-disposition")
	}
	dstf := matches[0][1]
	dst := path.Join(dstd, dstf)
	return dst, nil
}

package lib

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
	"github.com/tro3373/ydl/cmd/util"
)

func DownloadIfNeeded(repo, dstd string, target, sha256sum string) error {
	fmt.Printf("==> Start check update %s... target:%s, sha256sum:%s\n", repo, target, sha256sum)

	result, err := GetReleaseLatestTagInfo(repo)
	if err != nil {
		return err
	}

	tag := result["tag_name"].(string)
	if err != nil || len(tag) == 0 {
		return err
	}

	name := filepath.Base(repo)
	versionFile := filepath.Join(dstd, name+".version")
	currentVersion, err := util.ReadFileIfExist(versionFile)
	if err != nil {
		return err
	}
	if len(currentVersion) != 0 && tag == currentVersion {
		fmt.Printf("===> No update %s exist.\n", repo)
		return nil
	}

	fmt.Printf("===> Detect %s new tag: %s\n", repo, tag)

	if !util.Exists(dstd) {
		err = os.MkdirAll(dstd, 0775)
		if err != nil {
			return err
		}
	}

	fmt.Println("===> Downloading...")

	targetAssetNames := []string{target, sha256sum}
	results := make([]interface{}, 0)
	for _, asset := range result["assets"].([]interface{}) {
		a := asset.(map[string]interface{})
		name := a["name"]
		if len(targetAssetNames) != 0 && !util.Contains(targetAssetNames, name) {
			continue
		}
		results = append(results, a["id"])
	}
	c := make(chan int)
	for _, res := range results {
		go downloadResource(dstd, repo, res.(float64), c)
	}
	for i := 0; i < len(results); i++ {
		<-c
	}
	err = util.CheckSha256sum(
		filepath.Join(dstd, target),
		filepath.Join(dstd, sha256sum),
		target,
	)
	if err != nil {
		return err
	}
	return util.WriteFile(versionFile, tag)
}

func GetReleaseLatestTagInfo(repo string) (map[string]interface{}, error) {
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

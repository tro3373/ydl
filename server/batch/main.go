package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tro3373/ydl/batch/worker"
)

// const youtubeDl = "youtube-dl"
// const workd = "./work"

func main() {
	workDir := os.Args[1]

	ctx := initializeDir(workDir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	fmt.Println("==> Watching", ctx.QueueDir, "..")

	defer watcher.Close()
	done := make(chan bool)
	go dog(watcher, ctx)
	if err := watcher.Add(ctx.QueueDir); err != nil {
		fmt.Println("Failed to watcher.Add", err)
	}
	<-done
}

func initializeDir(workRootDir string) worker.Ctx {
	if len(workRootDir) == 0 {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		workRootDir = filepath.Join(filepath.Dir(ex), "work")
	}
	queue := createDirIfNotExist(workRootDir, "queue")
	lib := createDirIfNotExist(workRootDir, "lib")
	done := createDirIfNotExist(workRootDir, "done")

	return worker.NewCtx(workRootDir, queue, lib, done)
}

func createDirIfNotExist(dstRootDir, targetDir string) string {
	dir := filepath.Join(dstRootDir, targetDir)
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0775)
	}
	return dir
}

func dog(watcher *fsnotify.Watcher, ctx worker.Ctx) {
	for {
		select {
		case event := <-watcher.Events:
			// Receive event! fsnotify.Event{Name:"path/to/the/file", Op:0x1}
			fmt.Printf("Receive event! %#v\n", event)
			if event.Op&fsnotify.Create == fsnotify.Create {
				worker.Start(ctx, event)
			}
		case err := <-watcher.Errors:
			fmt.Println("Receive error!", err)
		}
	}
}

// func startTasks(dir string) error {
// 	err := updateYoutubeDlIfNeeded()
// 	if err != nil {
// 		return err
// 	}
// 	queued := filepath.Join(workd, "queue")
// 	jsons, err := findJsons(queued)
// 	if err != nil {
// 		return err
// 	}
// 	for _, json := range jsons {
// 		err := startDownload(json)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
//
// func findJsons(dir string) ([]string, error) {
// 	return filepath.Glob(filepath.Join(dir, "*", "*.json"))
// }
//
// // @see [Go client to github to get latest release and assets for a given repository](https://gist.github.com/metal3d/002e4f0d8545f83c2ace)
// func updateYoutubeDlIfNeeded() error {
// 	fmt.Println("==> Start checking", youtubeDl, "version..")
// 	repo := fmt.Sprintf("ytdl-org/%s", youtubeDl)
// 	result, err := getGhReleaseLatestTagInfo(repo)
//
// 	tag := result["tag_name"].(string)
// 	if err != nil || len(tag) == 0 {
// 		return err
// 	}
// 	// fmt.Printf("tag: %s, result:%#+v\n", tag, result)
// 	fmt.Printf("==> Detect %s tag: %s\n", repo, tag)
//
// 	libd := filepath.Join(workd, "lib")
// 	dstd := filepath.Join(libd, tag)
// 	if exists(dstd) {
// 		fmt.Printf("===> Already %s exist.\n", dstd)
// 		return nil
// 	}
// 	os.MkdirAll(dstd, 0775)
//
// 	fmt.Println("==> Start download", youtubeDl)
//
// 	results := make([]interface{}, 0)
// 	for _, asset := range result["assets"].([]interface{}) {
// 		a := asset.(map[string]interface{})
// 		name := a["name"]
// 		if name != youtubeDl && name != fmt.Sprintf("%s.sig", youtubeDl) {
// 			continue
// 		}
// 		results = append(results, a["id"])
// 	}
// 	c := make(chan int)
// 	for _, res := range results {
// 		go downloadResource(dstd, repo, res.(float64), c)
// 	}
// 	// TODO sig check
// 	for i := 0; i < len(results); i++ {
// 		<-c
// 	}
//
// 	dst := filepath.Join(libd, youtubeDl)
// 	if exists(dst) {
// 		err := os.Remove(dst)
// 		if err != nil {
// 			return errors.Wrapf(err, "Failed to Remove %s", dst)
// 		}
// 	}
// 	libdAbs, err := filepath.Abs(libd)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to get abs path for %s", libd)
// 	}
// 	src := filepath.Join(libdAbs, tag, youtubeDl)
// 	err = os.Chmod(src, 0775)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to chmod %s", src)
// 	}
// 	err = os.Symlink(src, dst)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to create symlink %s", src)
// 	}
// 	return nil
// }
//
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
// func exists(filename string) bool {
// 	_, err := os.Stat(filename)
// 	return os.IsExist(err)
// }
//
// func existsPrefix(name string) (bool, error) {
// 	matches, err := filepath.Glob(name + ".*")
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(matches) > 0, nil
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
//
// func startDownload(jsonPath string) error {
// 	raw, err := ioutil.ReadFile(jsonPath)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to read json %s", jsonPath)
// 	}
// 	var req request.Exec
// 	json.Unmarshal(raw, &req)
// 	return executeYoutubeDl(req)
// }
//
// func executeYoutubeDl(req request.Exec) error {
// 	var args []string
// 	key := req.Key()
// 	args = append(args, key)
// 	args = append(args, "-o")
// 	format := "%(id)s_%(title)s.%(ext)s"
// 	doned := filepath.Join(workd, "done")
// 	libd := filepath.Join(workd, "lib")
// 	args = append(args, filepath.Join(doned, format))
//
// 	// for audio output
// 	args = append(args, "-x")
// 	args = append(args, "--audio-format")
// 	args = append(args, "mp3")
//
// 	err := exec.Command(filepath.Join(libd, youtubeDl), args...).Run()
// 	if err != nil {
// 		log.Fatalf("Failed to executeYoutubeDl %s: %v", key, err)
// 		return err
// 	}
// 	return nil
// }

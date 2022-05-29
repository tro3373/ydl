package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tro3373/ydl/cmd/worker"
)

// // batchCmd represents the batch command
// var batchCmd = &cobra.Command{
// 	Use:   "batch",
// 	Short: "A brief description of your command",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:
//
// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		StartBatch(args)
// 	},
// }
//
// func init() {
// 	rootCmd.AddCommand(batchCmd)
//
// 	// Here you will define your flags and configuration settings.
//
// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// batchCmd.PersistentFlags().String("foo", "", "A help for foo")
//
// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// batchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
//
func StartBatch(ctx worker.Ctx) {
	err := worker.UpdateLibIFNeeded(ctx)
	if err != nil {
		fmt.Println("Update lib Error", err)
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("New Watcher Error", err)
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
			fmt.Printf("Receive event! Name:%s Op:%s\n", event.Name, event.Op.String())
			if event.Op&fsnotify.Create == fsnotify.Create {
				worker.Start(ctx, event)
			}
		case err := <-watcher.Errors:
			fmt.Println("Receive error!", err)
		}
	}
}

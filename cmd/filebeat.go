package cmd

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var FilebeatCmd = &cobra.Command{
	Use:   "filebeat",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//opt, _ := cmd.Flags().GetString("opt")
		path, _ := cmd.Flags().GetString("dir")
		interval, _ := cmd.Flags().GetInt("interval")
		g, _ := cmd.Flags().GetInt("goroutine")
		duration, _ := cmd.Flags().GetDuration("duration")
		rate, _ := cmd.Flags().GetInt("rate")

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		timer := time.NewTimer(duration)
		defer timer.Stop()
		dir := filepath.Dir(path)

		// 检查目录是否存在
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// 目录不存在，创建目录
			err := os.MkdirAll(dir, 0755) // 使用 MkdirAll 递归创建所需的所有父目录
			NoErr(err)
		}
		for {
			select {
			case <-ticker.C:
				for i := 0; i < g; i++ {
					go func() {
						for i := 0; i < rate; i++ {
							filename := fmt.Sprintf("%stestfile_%d.txt", dir, i)
							content := []byte("This is a test message.")
							// Create a new file and write the content to it
							if err := ioutil.WriteFile(filename, content, 0666); err != nil {
								panic(err)
							}

						}
					}()
				}

			}
		}

		time.Sleep(100000)
	},
}

func init() {

	FilebeatCmd.Flags().StringP("opt", "o", "", "")
	FilebeatCmd.Flags().StringP("dir", "d", "./input/", "")
	FilebeatCmd.Flags().IntP("interval", "i", 1, "")
	FilebeatCmd.Flags().IntP("goroutine", "g", 1, "")
	FilebeatCmd.Flags().DurationP("duration", "t", 0, "")
	FilebeatCmd.Flags().IntP("rate", "r", 1, "")

}

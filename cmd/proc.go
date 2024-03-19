package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "sysinfo is a CLI for system information",
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "List processes with optional name filtering",
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("name")
		processes, err := listProcesses(filter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		for _, p := range processes {
			fmt.Println(p)
		}
	},
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Show files opened by the process",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid := args[0]
		files, err := listOpenFiles(pid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		for _, f := range files {
			fmt.Println(f)
		}
	},
}

func listProcesses(filter string) ([]string, error) {
	// On UNIX systems, process information is available in /proc.
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var processes []string
	for _, dir := range dirs {
		if dir.IsDir() && isNumeric(dir.Name()) {
			cmdline, err := ioutil.ReadFile("/proc/" + dir.Name() + "/cmdline")
			if err != nil {
				continue
			}
			// The process name is the first argument in the cmdline file.
			// Splitting by the NULL character as cmdline arguments are separated by NULL in the filesystem.
			procName := strings.Split(string(cmdline), "\x00")[0]

			if filter == "" || strings.Contains(procName, filter) {
				processes = append(processes, fmt.Sprintf("PID: %s - Name: %s", dir.Name(), procName))
			}
		}
	}
	return processes, nil
}

func listOpenFiles(pid string) ([]string, error) {
	fdDir := fmt.Sprintf("/proc/%s/fd", pid)
	entries, err := ioutil.ReadDir(fdDir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		link, err := os.Readlink(fmt.Sprintf("%s/%s", fdDir, entry.Name()))
		if err != nil {
			continue
		}
		files = append(files, link)
	}
	return files, nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func init() {
	processCmd.Flags().StringP("name", "n", "", "Filter processes by name")
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(filesCmd)
}

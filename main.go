package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

const Version = "0.2.0"

var options struct {
	origin       string
	printVersion bool
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ws URL",
		Short: "websocket tool",
		Run:   root,
	}
	rootCmd.Flags().StringVarP(&options.origin, "origin", "o", "", "websocket origin")
	rootCmd.Flags().BoolVarP(&options.printVersion, "version", "v", false, "print version")

	rootCmd.Execute()
}

func root(cmd *cobra.Command, args []string) {
	if options.printVersion {
		fmt.Printf("ws v%s\n", Version)
		os.Exit(0)
	}

	if len(args) != 1 {
		cmd.Help()
		os.Exit(1)
	}

	url := args[0]
	origin := url
	if options.origin != "" {
		origin = options.origin
	}

	var historyFile string
	user, err := user.Current()
	if err == nil {
		historyFile = filepath.Join(user.HomeDir, ".ws_history")
	}

	err = connect(url, origin, &readline.Config{
		Prompt:      "> ",
		HistoryFile: historyFile,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err != io.EOF && err != readline.ErrInterrupt {
			os.Exit(1)
		}
	}
}

package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ws URL",
		Short: "websocket tool",
		Run:   root,
	}
	rootCmd.Flags().StringP("origin", "o", "", "websocket origin")

	rootCmd.Execute()
}

func root(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(1)
	}

	url := args[0]
	origin := url

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
		os.Exit(1)
	}
}

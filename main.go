package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("prompt", "> ")

	rootCmd := &cobra.Command{
		Use:   "ws URL",
		Short: "websocket tool",
		Run:   root,
	}
	rootCmd.Flags().StringP("origin", "o", "", "websocket origin")
	viper.BindPFlag("origin", rootCmd.Flags().Lookup("origin"))

	rootCmd.Execute()
}

func root(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(1)
	}

	url := args[0]
	origin := url
	if viper.IsSet("origin") {
		origin = viper.GetString("origin")
	}

	err := connect(url, origin, &readline.Config{
		Prompt:      viper.GetString("prompt"),
		HistoryFile: "/tmp/readline.tmp",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

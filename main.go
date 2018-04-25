package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

const Version = "0.2.1"

var options struct {
	origin       string
	printVersion bool
	authHeader   string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ws URL",
		Short: "websocket tool",
		Run:   root,
	}
	rootCmd.Flags().StringVarP(&options.origin, "origin", "o", "", "websocket origin")
	rootCmd.Flags().BoolVarP(&options.printVersion, "version", "v", false, "print version")
	rootCmd.Flags().StringVarP(&options.authHeader, "auth_header", "a", "", "auth header")

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

	dest, err := url.Parse(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var origin string
	if options.origin != "" {
		origin = options.origin
	} else {
		originURL := *dest
		if dest.Scheme == "wss" {
			originURL.Scheme = "https"
		} else {
			originURL.Scheme = "http"
		}
		origin = originURL.String()
	}

	var authHeader string
	authHeader = ""
	if options.authHeader != "" {
		authHeader = options.authHeader
	}

	var historyFile string
	user, err := user.Current()
	if err == nil {
		historyFile = filepath.Join(user.HomeDir, ".ws_history")
	}

	err = connect(dest.String(), origin, authHeader, &readline.Config{
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

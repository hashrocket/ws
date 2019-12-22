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
	insecure     bool
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ws URL",
		Short: "websocket tool",
		Run:   root,
	}
	rootCmd.Flags().StringVarP(&options.origin, "origin", "o", "", "websocket origin")
	rootCmd.Flags().BoolVarP(&options.printVersion, "version", "v", false, "print version")
	rootCmd.Flags().BoolVarP(&options.insecure, "insecure", "k", false, "skip ssl certificate check")

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

	// Correct and add missing schemes.
	switch dest.Scheme {
	case "ws", "wss":
	case "http":
		dest.Scheme = "ws"
	case "https":
		dest.Scheme = "wss"
	default:
		// Likely no scheme at all, e.g. "localhost:8000".
		dest, err = url.Parse("ws://" + args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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

	var historyFile string
	user, err := user.Current()
	if err == nil {
		historyFile = filepath.Join(user.HomeDir, ".ws_history")
	}

	err = connect(dest.String(), origin, &readline.Config{
		Prompt:      "> ",
		HistoryFile: historyFile,
	}, options.insecure)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err != io.EOF && err != readline.ErrInterrupt {
			os.Exit(1)
		}
	}
}

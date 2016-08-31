package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"golang.org/x/net/websocket"
)

func connect(url, origin string, rlConf *readline.Config) error {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(rlConf)
	if err != nil {
		return err
	}
	defer rl.Close()

	go rx(ws, rl.Stdout())

	for {
		line, err := rl.Readline()
		switch err {
		case io.EOF, readline.ErrInterrupt:
			return nil
		case nil:
		default:
			return err
		}

		_, err = io.WriteString(ws, line)
		if err != nil {
			return err
		}
	}
}

func rx(ws *websocket.Conn, stdout io.Writer) {
	buf := make([]byte, 4096)

	for {
		n, err := ws.Read(buf)
		if n > 0 {
			fmt.Fprintln(stdout, string(buf[:n]))
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

package main

import (
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"golang.org/x/net/websocket"
)

type session struct {
	ws      *websocket.Conn
	rl      *readline.Instance
	errChan chan error
}

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

	sess := &session{
		ws:      ws,
		rl:      rl,
		errChan: make(chan error),
	}

	go sess.readConsole()
	go sess.readWebsocket()

	return <-sess.errChan
}

func (s *session) readConsole() {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			s.errChan <- err
			return
		}

		_, err = io.WriteString(s.ws, line)
		if err != nil {
			s.errChan <- err
			return
		}
	}
}

func (s *session) readWebsocket() {
	buf := make([]byte, 4096)

	for {
		n, err := s.ws.Read(buf)
		if n > 0 {
			fmt.Fprintln(s.rl.Stdout(), "<", string(buf[:n]))
		}
		if err != nil {
			s.errChan <- err
			return
		}
	}
}

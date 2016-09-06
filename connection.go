package main

import (
	"fmt"
	"net/http"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type session struct {
	ws      *websocket.Conn
	rl      *readline.Instance
	errChan chan error
}

func connect(url, origin string, rlConf *readline.Config) error {
	headers := make(http.Header)
	headers.Add("Origin", origin)

	ws, _, err := websocket.DefaultDialer.Dial(url, headers)
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

		err = s.ws.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			s.errChan <- err
			return
		}
	}
}

func (s *session) readWebsocket() {
	rxSprintf := color.New(color.FgGreen).SprintfFunc()

	for {
		_, buf, err := s.ws.ReadMessage()
		if err != nil {
			s.errChan <- err
			return
		}
		fmt.Fprint(s.rl.Stdout(), rxSprintf("< %s\n", string(buf)))
	}
}

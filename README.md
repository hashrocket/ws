# ws

ws is a simple command line websocket client designed for exploring and debugging websocket servers.

```
$ ws ws://localhost:3000/ws
> {"type": "echo", "payload": "Hello, world"}
< {"type":"echo","payload":"Hello, world"}
> {"type": "broadcast", "payload": "Hello, world"}
< {"type":"broadcast","payload":"Hello, world"}
< {"type":"broadcastResult","payload":"Hello, world","listenerCount":1}
> ^D
EOF
```

ws includes readline-style keyboard shortcuts and history.

## WebSocks

[中文说明](https://github.com/lzjluzijie/websocks/blob/master/README-zh.md)

A secure proxy based on websocket.

This project is still working in progress, more features are still in development. If you are interested in this project, please star this project in order to support me. Thank you.

If you have any problems or suggestions, please do not hesitate to submit issues or contact me @halulu

### Example (Enable TLS)

#### Server

```
go get -v github.com/lzjluzijie/websocks/cmd/websocks-server
websocks-server -l :2333 -p /password
```

#### Local

```
go get -v github.com/lzjluzijie/websocks/cmd/websocks-local
websocks-local -l :1080 -u wss://server.com/password
```

#### Caddyfile
```
https://server.com {
  proxy /password localhost:2333 {
    websocket
  }
}
```

### TO-DO

 - [ ] Configuration
 - [ ] Optimize code
 - [ ] WebSocket mux
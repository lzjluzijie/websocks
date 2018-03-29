## WebSocks

一个基于 WebSocket 的代理工具

请注意，本项目目前还在开发中，仅供测试使用，更多功能仍在完善中。

目前已经可以测试使用了，具体可以看下面的介绍。有任何问题或建议可以直接联系我[@halulu](https://t.me/halulu)

### Example (Enable tls)

#### Server

```
go get -v github.com/lzjluzijie/websocks/cmd/websocks-server
websocks-server -l :2333
```

#### Local

```
go get -v github.com/lzjluzijie/websocks/cmd/websocks-local
websocks-local -l :1080 -u wss://server.com/ws
```

#### Caddyfile
```
https://server.com {
  proxy /ws localhost:2333 {
    websocket
  }
}
```

### TO-DO

 - [ ] Config
 - [ ] 优化代码
 - [ ] ws复用

优点
 - ws+tls 不是伪装而是正经的网站，很隐蔽
 - 可以走CDN，根本不怕被墙

缺点
 - 可能比较慢
 - 配套软件差

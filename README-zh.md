## WebSocks

一个基于 WebSocket 的代理工具

本项目目前还在开发中，更多功能仍在完善中。如果你对这个项目感兴趣，请star它来支持我，蟹蟹。

有任何问题或建议可以直接发issue或者联系我 [@halulu](https://t.me/halulu)

开发记录可以看[我的博客](https://halu.lu/post/websocks-development/)

优点
 - ws+tls 不是伪装而是正经的网站，很隐蔽
 - 可以走CDN，根本不怕被墙

缺点
 - 可能比较慢
 - 配套软件差

### 示例 (开启 TLS)

#### 服务端

```
go get -v github.com/lzjluzijie/websocks/cmd/websocks-server
websocks-server -l :2333 -p /password
```

#### 客户端

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

 - [ ] Config
 - [ ] 优化代码
 - [ ] ws复用

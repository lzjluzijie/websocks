## WebSocks

一个基于 WebSocket 的代理工具

本项目目前还在开发中，更多功能仍在完善中。如果你对这个项目感兴趣，请star它来支持我，蟹蟹。

有任何问题或建议可以直接发issue或者联系我 [@halulu](https://t.me/halulu)

开发记录可以看[我的博客](https://halu.lu/post/websocks-development/)

优点
 - 使用WS+TLS，十分安全且不易被检测，和普通HTTPS网站一样
 - 可以搭配使用cloudflare这类cdn，完全不怕被墙！

缺点就是刚刚开始开发，没有GUI客户端，功能也比较少，如果你能来帮我那就太好了！

## 一键脚本
### 项目首页:[1715173329/websocks-onekey](https://github.com/1715173329/websocks-onekey)
```bash
apt-get install -y curl && curl -O https://raw.githubusercontent.com/lzjluzijie/websocks/master/script/websocks-go.sh && bash websocks-go.sh
```

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

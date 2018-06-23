# WebSocks

[English](https://github.com/lzjluzijie/websocks/blob/master/README-en.md)

一个基于 WebSocket 的代理工具

本项目目前还在开发中，更多功能仍在完善中。如果你对这个项目感兴趣，请star它来支持我，蟹蟹

有任何问题或建议可以直接发issue或者联系我 [@halulu](https://t.me/halulu)，也可以来[TG群](https://t.me/websocks)水一水，开发记录可以看[我的博客](https://halu.lu/post/websocks-development/)

优点
 - 使用WS+TLS，十分安全且不易被检测，和普通HTTPS网站一样
 - 可以搭配使用cloudflare这类cdn，完全不怕被墙！

缺点就是刚刚开始开发，没有GUI客户端，功能也比较少，如果你能来帮我那就太好了！

[官网](https://websocks.org/)|[社区](https://zhuji.lu/tags/websocks)|[一键脚本](https://zhuji.lu/topic/15/websocks-一键脚本-简易安装教程)|[电报群](https://t.me/websocks)

## 示例

### 内置 TLS 混淆域名并反向代理

#### 服务端
```
./websocks cert
./websocks config server -l :2333 -p /password --reverse-proxy https://www.centos.org/ --tls
./websocks server
```

#### 客户端
```
./websocks config client -l :1080 -s wss://server.com:2333/password -n www.centos.com --insecure
./websocks client
```


### Caddy TLS

#### 服务端
```
./websocks config server -l :2333 -p /password
./websocks server
```

#### 客户端
```
./websocks config client -l :1080 -s wss://server.com/password
./websocks client
```

#### Caddyfile
```
https://server.com {
  proxy /password localhost:2333 {
    websocket
  }
}
```
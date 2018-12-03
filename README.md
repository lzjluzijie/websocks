# WebSocks

[English](https://github.com/lzjluzijie/websocks/blob/master/README-en.md)

一个基于 WebSocket 的代理工具

**由于本人学业的原因，websocks暂时停更几个月，各位大佬们对不住了，等我搞定大学一定会填坑的**

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
./websocks server -l :443 -p websocks --reverse-proxy http://mirror.centos.org --tls
```

#### 客户端
```
./websocks client -l :1080 -s wss://websocks.org:443/websocks -sni mirror.centos.com --insecure
```


### Caddy TLS

#### 服务端
```
./websocks server -l :2333 -p /websocks
```

#### 客户端
```
./websocks client -l :1080 -s wss://websocks.org/websocks
```

#### Caddyfile
```
https://websocks.org {
  proxy /websocks localhost:2333 {
    websocket
  }
}
```
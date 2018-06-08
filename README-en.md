# WebSocks

A secure proxy based on websocket.

This project is still working in progress, more features are still in development. If you are interested in this project, please star this project in order to support me. Thank you.

If you have any problems or suggestions, please do not hesitate to submit issues or contact me [@halulu](https://t.me/halulu).

Advantages:

- Using WebSocket and TLS which are very secure and difficult to be detected, same as regular HTTPS websites
- Can be used with cdn such as cloudflare, not afraid of gfw at all!

The disadvantage is that I have just started development, there is no GUI client, and features are not enough. I will appreciate if you can help me!

[Official site](https://websocks.org/)|[Community](https://zhuji.lu/tags/websocks)|[Test node](https://zhuji.lu/topic/39/websocks测试节点)|[One-click script](https://zhuji.lu/topic/15/websocks-一键脚本-简易安装教程)|[Telegram group](https://t.me/websocks)

## Example

### Built-in TLS with fake server name and reversing proxy

#### Server
```
./websocks cert
./websocks config server -l :2333 -p /password --reverse-proxy https://www.centos.org/ --tls
./websocks server
```

#### Local
```
./websocks config client -l :1080 -s wss://server.com:2333/password -n www.centos.com --insecure
./websocks client
```

### Caddy TLS

#### Server
```
./websocks config server -l :2333 -p /password
./websocks server
```

#### Local
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
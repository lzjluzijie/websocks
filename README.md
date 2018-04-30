# WebSocks

[中文说明](https://github.com/lzjluzijie/websocks/blob/master/README-zh.md)

A secure proxy based on websocket.

This project is still working in progress, more features are still in development. If you are interested in this project, please star this project in order to support me. Thank you.

If you have any problems or suggestions, please do not hesitate to submit issues or contact me [@halulu](https://t.me/halulu)

Advantages:

- Using WebSocket and TLS which are very secure and difficult to be detected, same as regular HTTPS websites
- Can be used with cdn such as cloudflare, not afraid of gfw at all!

The disadvantage is that I have just started development, there is no GUI client, and features are not enough. I will appreciate if you can help me!

To-Do:

## Example

### Caddy TLS

#### Server
```
./websocks server -l :2333 -p /password
```

#### Local
```
./websocks client -l :1080 -s wss://server.com/password
```

#### Caddyfile
```
https://server.com {
  proxy /password localhost:2333 {
    websocket
  }
}
```

### Built-in TLS with fake server name

#### Server
```
./websocks cert
./websocks server -l :2333 -p /password --tls
```

#### Local
```
./websocks client -l :1080 -s wss://server.com:2333/password -n github.com --insecure
```

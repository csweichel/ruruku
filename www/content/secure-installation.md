---
title: Securing your installation
weight: 3
menu: true
hideFromIndex: true
---

By default Ruruku runs without any transport encryption/security. This way getting started is easy, but it's not exactly secure.

Especially if you intent to expose your Ruruku installation to the internet, it's a good idea to use TLS for the gRPC API server,
and HTTPS/TLS for the UI server.

To enable TLS for the API server and HTTPS for the UI server, add the following to your config
```
server:
    ui:
        port: 443
        https: enabled
        cert: ui.crt
        key: ui.key
    tls:
        enabled: true
        cert: server.crt
        key: server.key
```
and use the `--tls` flag to pass the certificate to the client (e.g. `ruruku session --tls server.crt list`).

Hint: to generate a pair of self-signed certificates run:
```
openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -subj "/CN=*/" -out server.crt -keyout server.key
```

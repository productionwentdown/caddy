
# caddy

A tiny 4MB Caddy image compressed with [UPX](https://github.com/upx/upx).

# Caveats 

Since this image is `FROM scratch`, it does not have the certificates 
necessary to connect to external HTTPS servers, including Let's Encrypt's 
ACME server. This means that automatic TLS will not work in this Docker 
image. 

TODO: add [ca-certificates.crt](https://curl.haxx.se/ca/cacert.pem) in /etc/ssl/certs/

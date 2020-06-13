FROM caddy:2.1.0-beta.1
2.0.0-alpine as build

RUN apk add --no-cache upx ca-certificates     && upx --ultra-brute /usr/bin/caddy     && caddy version

FROM scratch

COPY --from=build /etc/caddy/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /usr/share/caddy/index.html /usr/share/caddy/index.html
COPY --from=build /usr/bin/caddy /usr/bin/caddy
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV XDG_CONFIG_HOME=/config
ENV XDG_DATA_HOME=/data
ENV HOME=/

VOLUME /config
VOLUME /data

LABEL org.opencontainers.image.version=v2.1.0-beta.1
2.0.0
LABEL org.opencontainers.image.title="Caddy (productionwentdown fork)"
LABEL org.opencontainers.image.description="a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go"
LABEL org.opencontainers.image.url=https://productionwentdown.makerforce.io
LABEL org.opencontainers.image.documentation=https://caddyserver.com/docs
LABEL org.opencontainers.image.vendor="productionwentdown"
LABEL org.opencontainers.image.licenses=Apache-2.0
LABEL org.opencontainers.image.source="https://github.com/productionwentdown/caddy"

EXPOSE 80
EXPOSE 443
EXPOSE 2019

CMD ["/usr/bin/caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]

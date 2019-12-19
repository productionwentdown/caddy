#
# Build stage
#
FROM golang:1.13-alpine as build

# args
ARG version="1.0.4"
# add plugin import paths here separated by commas
ARG plugins=""
ARG telemetry="true"

# build root
WORKDIR /build

# plugins
COPY plugger.go ./

# build & test
RUN apk add --no-cache git upx ca-certificates \
    && echo -e "module caddy\nrequire github.com/caddyserver/caddy v${version}" > go.mod \
    && go run plugger.go -plugins="${plugins}" -telemetry="${telemetry}" \
    && CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build \
    && upx --ultra-brute caddy \
    && ./caddy -version

#
# Final image
#
FROM scratch

# labels
LABEL org.label-schema.vcs-url="https://github.com/productionwentdown/caddy"
LABEL org.label-schema.version=${version}
LABEL org.label-schema.schema-version="1.0"

# copy binary and ca certs
COPY --from=build /build/caddy /bin/caddy
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# copy default caddyfile
COPY Caddyfile /etc/Caddyfile

# set default caddypath
ENV CADDYPATH=/etc/.caddy
VOLUME /etc/.caddy

# serve from /srv
WORKDIR /srv
EXPOSE 2015

ENTRYPOINT ["/bin/caddy", "--conf", "/etc/Caddyfile", "--log", "stdout"]

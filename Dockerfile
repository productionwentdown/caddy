#
# Build stage by @abiosoft https://github.com/abiosoft/caddy-docker
#
FROM golang:1.12-alpine as build

# args
ARG version="0.11.5"
# add plugins here separated by commas
ARG plugins=""

# deps
RUN apk add --no-cache git

# source
RUN git clone https://github.com/mholt/caddy -b "v${version}" $GOPATH/src/github.com/mholt/caddy
WORKDIR $GOPATH/src/github.com/mholt/caddy
RUN git checkout -b "v${version}"

# plugin helper
RUN go get -v github.com/abiosoft/caddyplug/caddyplug

# plugins
RUN for plugin in $(echo $plugins | tr "," " "); do \
    go get -v $(caddyplug package $plugin); \
    printf "package caddyhttp\nimport _ \"$(caddyplug package $plugin)\"" > \
        $GOPATH/src/github.com/mholt/caddy/caddyhttp/$plugin.go ; \
    done

# builder dependency
RUN git clone https://github.com/caddyserver/builds $GOPATH/src/github.com/caddyserver/builds

# build
WORKDIR $GOPATH/src/github.com/mholt/caddy/caddy
RUN git checkout -f
RUN go run build.go
RUN mv caddy /


#
# Compress Caddy with upx
#
FROM debian:stable as compress

ARG upx_version="3.94"

# dependencies
RUN apt-get update && apt install -y --no-install-recommends \
    tar \
    xz-utils \
    curl \
    ca-certificates

# get official upx binary
RUN curl --silent --show-error --fail --location -o - \
    "https://github.com/upx/upx/releases/download/v${upx_version}/upx-${upx_version}-amd64_linux.tar.xz" \
    | tar --no-same-owner -C /usr/bin/ -xJ \
    --strip-components 1 upx-${upx_version}-amd64_linux/upx

# copy and compress
COPY --from=build /caddy /usr/bin/caddy
RUN /usr/bin/upx --ultra-brute /usr/bin/caddy

# test
RUN /usr/bin/caddy -version


#
# Final image
#
FROM scratch

# labels
LABEL org.label-schema.vcs-url="https://github.com/productionwentdown/caddy"
LABEL org.label-schema.version=${version}
LABEL org.label-schema.schema-version="1.0"

# copy binary and ca certs
COPY --from=compress /usr/bin/caddy /bin/caddy
COPY --from=compress /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# copy default caddyfile
COPY Caddyfile /etc/Caddyfile

# set default caddypath
ENV CADDYPATH=/etc/.caddy
VOLUME /etc/.caddy

# serve from /srv
WORKDIR /srv

ENTRYPOINT ["/bin/caddy", "--conf", "/etc/Caddyfile", "--log", "stdout"]

FROM debian:stable as fetch

RUN apt-get update && apt install -y --no-install-recommends \
    tar \
    xz-utils \
    curl \
    ca-certificates

RUN curl --silent --show-error --fail --location -o - \
    "https://caddyserver.com/download/linux/amd64" \
    | tar --no-same-owner -C /usr/bin/ -xz caddy

RUN curl --silent --show-error --fail --location -o - \
    "https://github.com/upx/upx/releases/download/v3.94/upx-3.94-amd64_linux.tar.xz" \
    | tar --no-same-owner -C /usr/bin/ -xJ \
    --strip-components 1 upx-3.94-amd64_linux/upx

RUN ls -l /usr/bin/caddy
RUN /usr/bin/upx --ultra-brute /usr/bin/caddy
RUN ls -l /usr/bin/caddy

RUN /usr/bin/caddy -version


FROM scratch

LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.vcs-url="https://github.com/productionwentdown/caddy"
LABEL org.label-schema.version=$VERSION
LABEL org.label-schema.schema-version="1.0"

COPY --from=fetch /usr/bin/caddy /bin/caddy
COPY --from=fetch /etc/ssl/certs/ca-certificates.crt
COPY Caddyfile /etc/Caddyfile

ENV CADDYPATH=/etc/.caddy
VOLUME /etc/.caddy

WORKDIR /srv
COPY index.html /srv/index.html

ENTRYPOINT ["/bin/caddy"]
CMD ["--conf", "/etc/Caddyfile", "--log", "stdout"]

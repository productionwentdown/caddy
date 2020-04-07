
# WARNING: Upcoming Caddy v2 WILL fundamentally change the way Caddy is configured. The `latest` tag will be changed to reflect that. This LIKELY breaks existing containers using the `latest` tag and the stock Caddyfile

# [caddy](https://hub.docker.com/r/productionwentdown/caddy/) [![](https://images.microbadger.com/badges/version/productionwentdown/caddy.svg)](https://microbadger.com/images/productionwentdown/caddy "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/productionwentdown/caddy.svg)](https://microbadger.com/images/productionwentdown/caddy "Get your own image badge on microbadger.com")

A tiny 9MB Caddy image compressed with [UPX](https://github.com/upx/upx).

# Usage

Serve files in `$PWD`:
```
docker run -it --rm -p 2015:2015 -v $PWD:/srv productionwentdown/caddy
```

Overwrite `Caddyfile`:
```
docker run -it --rm -p 2015:2015 -v $PWD:/srv -v $PWD/Caddyfile:/etc/Caddyfile productionwentdown/caddy
```

Persist `.caddy` to avoid hitting Let's Encrypt's rate limit:
```
docker run -it --rm -p 2015:2015 -v $PWD:/srv -v $PWD/Caddyfile:/etc/Caddyfile -v $HOME/.caddy:/etc/.caddy productionwentdown/caddy
```

# Build with plugins

Using `docker build` arguments:
```
docker build -t caddy --build-arg plugins=github.com/abiosoft/caddy-git,github.com/zikes/gopkg .
```

You can also fork and edit plugger.go for more advanced plugin configuration

# Build without telemetry

```
docker build . -t caddy --build-arg telemetry=false
```

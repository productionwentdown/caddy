
# [caddy](https://hub.docker.com/r/productionwentdown/caddy/) [![](https://images.microbadger.com/badges/version/productionwentdown/caddy.svg)](https://microbadger.com/images/productionwentdown/caddy "Get your own version badge on microbadger.com") [![](https://images.microbadger.com/badges/image/productionwentdown/caddy.svg)](https://microbadger.com/images/productionwentdown/caddy "Get your own image badge on microbadger.com")

A tiny 4MB Caddy image compressed with [UPX](https://github.com/upx/upx).

> Notice: I keep this manually updated. If it goes out of date, ping me via Twitter [@serverwentdown](https://twitter.com/serverwentdown).

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

.PHONY: all
all: update library/caddy

.PHONY: update
update: scratch/Dockerfile
	./update.sh

library/caddy: generate-stackbrew-library.sh
	./generate-stackbrew-library.sh > library/caddy

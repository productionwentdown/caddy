.PHONY: all
all: update library/caddy

.PHONY: update
update: update.sh
	./update.sh

.PHONY: library/caddy
library/caddy: generate-stackbrew-library.sh
	./generate-stackbrew-library.sh > library/caddy

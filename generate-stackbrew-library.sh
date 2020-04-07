#!/bin/bash

set -Eeuo pipefail

gitHubUrl='https://github.com/productionwentdown/caddy'
gitHubUpstreamUrl='https://github.com/caddyserver/caddy-docker'

join() {
	local sep="$1"; shift
	local out; printf -v out "${sep//%/%%}%s" "$@"
	echo "${out#$sep}"
}

_wget() {
	wget -qO- -o /dev/null "$@"
}

cat <<-EOH
# see https://github.com/caddyserver/caddy-docker

Maintainers: Ambrose Chua <ambrose@makerforce.io> (@serverwentdown)
GitRepo: https://github.com/productionwentdown/caddy.git
GitCommit: $(git log --format='format:%H' -1)

EOH

caddyVersion="$(_wget "$gitHubUpstreamUrl/raw/master/stackbrew-config.yaml" | grep -oP '(?<=caddy_version: ).+$')"
version="${caddyVersion#v}"
versionAliases=($version)

arches=(amd64)

if [[ "$version" =~ "-rc" ]]; then
	versionAliases+=(test)
elif [[ "$version" =~ "-beta" ]]; then
	versionAliases+=()
else
	versionAliases+=(latest)
fi

cat <<-EOE
	Tags: $(join ', ' "${versionAliases[@]}")
	Architectures: $(join ', ' "${arches[@]}")
	Directory: scratch
EOE

package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/docker-library/go-dockerlibrary/manifest"
)

var manifestURL = "https://github.com/docker-library/official-images/raw/master/library/caddy"

var maintainers = []string{"Ambrose Chua <ambrose@makerforce.io> (@serverwentdown)"}
var gitRepo = "https://github.com/productionwentdown/caddy.git"

var readmeTemplate = template.Must(template.New("readme").Parse(`
# [caddy](https://hub.docker.com/r/productionwentdown/caddy) ![Docker Pulls](https://img.shields.io/docker/pulls/productionwentdown/caddy?style=flat-square) ![Docker Image Size (tag)](https://img.shields.io/docker/image-size/productionwentdown/caddy/alpine?style=flat-square)

A tiny &lt;10MB Caddy image compressed with [UPX](https://github.com/upx/upx).

- [Docker Hub](https://hub.docker.com/r/productionwentdown/caddy)
- [GitHub](https://github.com/productionwentdown/caddy)

## Simple Tags

New versions are tracked within 4 hours. Currently available versions:

{{range $entry := .Entries}}* {{range $index, $tag := $entry.Tags}}{{if $index}}, {{end}}` + "`{{$tag}}`" + `{{end}}
{{end}}
## Shared Tags

{{range $entry := .Entries}}* {{range $index, $tag := $entry.SharedTags}}{{if $index}}, {{end}}` + "`{{$tag}}`" + `{{end}}
  * ` + "`{{index $entry.Tags 0}}`" + `
{{end}}
# Usage

See the [official image](https://hub.docker.com/_/caddy) for documentation. This image behaves the same way, except that it is slightly slimmer.
`))

var dockerfileTemplate = template.Must(template.New("name").Parse(`FROM caddy:{{.Tag}} as build

RUN apk add --no-cache upx ca-certificates \
    && upx --ultra-brute /usr/bin/caddy \
    && caddy version

FROM scratch

COPY --from=build /etc/caddy/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /usr/share/caddy/index.html /usr/share/caddy/index.html
COPY --from=build /usr/bin/caddy /usr/bin/caddy
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/nsswitch.conf /etc/nsswitch.conf

ENV XDG_CONFIG_HOME=/config
ENV XDG_DATA_HOME=/data
ENV HOME=/

VOLUME /config
VOLUME /data

LABEL org.opencontainers.image.version=v{{.Version}}
LABEL org.opencontainers.image.title="Caddy (productionwentdown build)"
LABEL org.opencontainers.image.description="a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go"
LABEL org.opencontainers.image.url=https://productionwentdown.makerforce.io
LABEL org.opencontainers.image.documentation=https://caddyserver.com/docs
LABEL org.opencontainers.image.vendor="productionwentdown"
LABEL org.opencontainers.image.licenses=Apache-2.0
LABEL org.opencontainers.image.source="https://github.com/productionwentdown/caddy"

EXPOSE 80
EXPOSE 443
EXPOSE 2019

WORKDIR /srv

CMD ["/usr/bin/caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
`))

type dockerfileData struct {
	Tag     string
	Version string
	Arch    string
}

func main() {
	doDockerfiles := flag.Bool("dockerfiles", false, "Update Dockerfiles")
	doManifest := flag.Bool("manifest", false, "Update manifest")
	doReadme := flag.Bool("readme", false, "Update README.md")
	commit := flag.String("commit", "", "Current commit hash")
	flag.Parse()

	_, _, man, err := manifest.Fetch("nonexistent", manifestURL)
	if err != nil {
		panic(err)
	}

	man = subsetManifest(man)

	if *doDockerfiles {
		err = updateDockerfiles(man)
		if err != nil {
			panic(err)
		}
	}
	if *doManifest {
		err = updateManifest(man, *commit)
		if err != nil {
			panic(err)
		}
	}
	if *doReadme {
		err = updateReadme(man)
		if err != nil {
			panic(err)
		}
	}
}

// subsetManifest returns the subset of the manifest that can be built
func subsetManifest(man *manifest.Manifest2822) *manifest.Manifest2822 {
	entries := make([]manifest.Manifest2822Entry, 0, 1)
	for _, entry := range man.Entries {

		isWindows := false
		for _, arch := range entry.Architectures {
			if strings.Contains(arch, "windows") {
				isWindows = true
			}
		}

		isBuilder := false
		for _, tag := range entry.Tags {
			if strings.Contains(tag, "builder") {
				isBuilder = true
			}
		}

		if isWindows || isBuilder {
			continue
		}
		entries = append(entries, entry)
	}
	return &manifest.Manifest2822{
		Global:  man.Global,
		Entries: entries,
	}
}

// updateDockerfiles generates a new set of Dockerfiles based on the manifest
func updateDockerfiles(man *manifest.Manifest2822) error {
	for _, entry := range man.Entries {
		tag := entry.Tags[0]
		version := strings.Split(entry.Tags[0], "-")[0]
		for _, arch := range entry.Architectures {
			data := dockerfileData{Tag: tag, Version: version, Arch: arch}
			err := writeDockerfile(entry, data)
			if err != nil {
				return err
			}
		}
	}
	// TODO: Remove old Dockerfiles
	return nil
}

func writeDockerfile(entry manifest.Manifest2822Entry, data dockerfileData) error {
	dockerfilePath := path.Join(entry.ArchDirectory(data.Arch), entry.ArchFile(data.Arch))
	err := os.MkdirAll(entry.ArchDirectory(data.Arch), 0755)
	if err != nil {
		return err
	}
	dockerfileFile, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dockerfileFile.Close()

	err = dockerfileTemplate.Execute(dockerfileFile, data)
	if err != nil {
		return err
	}

	return nil
}

// updateManifest replaces upstream maintainer and Git info with this repo
func updateManifest(man *manifest.Manifest2822, gitCommit string) error {
	man.Global.Maintainers = maintainers
	for i := range man.Entries {
		entry := &man.Entries[i]
		entry.Maintainers = maintainers
		entry.GitRepo = gitRepo
		entry.GitCommit = gitCommit
	}
	manifestPath := path.Join("library", "caddy")
	err := os.MkdirAll("library", 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(manifestPath, []byte(man.String()), 0644)
}

// updateReadme generates a new README.md
func updateReadme(man *manifest.Manifest2822) error {
	readmeFile, err := os.OpenFile("README.md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	err = readmeTemplate.Execute(readmeFile, man)
	if err != nil {
		return err
	}

	return nil
}

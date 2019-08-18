// +build ignore

package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
)

var plugins string
var telemetry bool

func main() {
	flag.StringVar(&plugins, "plugins", "", "Specify plugins by full paths, seperated by commas")
	flag.BoolVar(&telemetry, "telemetry", false, "Enable telemetry")
	flag.Parse()

	d := &data{
		Plugins:         strings.FieldsFunc(plugins, func(c rune) bool { return c == ',' }),
		EnableTelemetry: telemetry,
	}

	log.Printf("Additional plugins: %v", d.Plugins)
	log.Printf("Enabled telemetry: %v", d.EnableTelemetry)

	f, err := os.Create("caddy.go")
	defer f.Close()
	if err != nil {
		log.Fatal("Unable to open file")
	}

	t := template.Must(template.New("caddy.go").Parse(caddyTemplate))
	t.Execute(f, d)
}

type data struct {
	Plugins         []string
	EnableTelemetry bool
}

var caddyTemplate = `
package main

import (
	"github.com/caddyserver/caddy/caddy/caddymain"

	// plug in plugins here
	{{range $plugin := .Plugins}}
	_ "{{$plugin}}"
	{{end}}
)

func main() {
	// optional: disable telemetry
	caddymain.EnableTelemetry = {{.EnableTelemetry}}
	caddymain.Run()
}
`

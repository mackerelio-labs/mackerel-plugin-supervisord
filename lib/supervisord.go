package supervisord

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/kolo/xmlrpc"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

type process struct {
	Name  string `xmlrpc:"name"`
	State int    `xmlrpc:"state"`
}

type supervisordPlugin struct {
	Uri string
}

func (m supervisordPlugin) FetchMetrics() (map[string]float64, error) {
	transport := http.DefaultTransport
	if strings.HasPrefix(m.Uri, "unix:") {
		path, _ := strings.CutPrefix(m.Uri, "unix:")
		transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		}
		m.Uri = "http://unix/RPC2"
	}

	client, err := xmlrpc.NewClient(m.Uri, transport)
	if err != nil {
		log.Fatal(err)
	}
	var result []process
	err = client.Call("supervisor.getAllProcessInfo", nil, &result)
	if err != nil {
		log.Fatal(err)
	}

	stat := make(map[string]float64)
	for _, i := range result {
		stat[fmt.Sprintf("supervisord.%s.state", i.Name)] = float64(i.State)
	}
	return stat, nil
}

func (m supervisordPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"supervisord.#": {
			Label: "Ssupervisord Process State",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "state", Label: "State", Diff: false},
			},
		},
	}
}

func Do() {
	uri := flag.String("uri", "unix:/var/run/supervisor.sock", "socket `uri`")
	flag.Parse()

	var plugin supervisordPlugin
	plugin.Uri = *uri
	helper := mp.NewMackerelPlugin(plugin)
	helper.Run()
}

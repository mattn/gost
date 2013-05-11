package main

import (
    "fmt"
    "encoding/json"
    "flag"
    "log"
    "net/http"
    "os"
    "os/exec"
)

type config struct {
	Addr string `json:"addr"`
	Root string `json:"root"`
	RPC string `json:"rpc"`
	Apps map[string]struct {
		Name string `json:"name"`
		Path string `json:"path"`
		PullCommand string `json:"pull_command"`
	} `json:"apps"`
}

type payload struct {
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

var configFile = flag.String("c", "config.json", "config file")

func main() {
	flag.Parse()

	f, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	var c config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Root == "" || c.Root[0] != '/' {
		c.Root = "/" + c.Root
	}

    http.HandleFunc(c.Root, func (w http.ResponseWriter, r *http.Request) {
		var p payload
		if json.Unmarshal([]byte(r.FormValue("payload")), &p) == nil {
			app, ok := c.Apps[p.Repository.Name]
			if ok {
				cs := []string{"/bin/bash", "-c", app.PullCommand}
				cmd := exec.Command(cs[0], cs[1:]...)
				err = cmd.Run()
				if err != nil {
					log.Println(err)
				}
			}
		}
		fmt.Fprintf(w, "")
	})
    http.ListenAndServe(c.Addr, nil)
}

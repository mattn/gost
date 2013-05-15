package main

import (
    "fmt"
    "encoding/json"
    "errors"
    "flag"
    "log"
    "net/http"
    "net/rpc"
    "os"
    "os/exec"
    "runtime"
)

type config struct {
	Addr string `json:"addr"`
	Root string `json:"root"`
	RPC string `json:"rpc"`
	Apps map[string]struct {
		Proc string `json:"proc"`
		Path string `json:"path"`
		UpdateCommand string `json:"update_command"`
		BuildCommand string `json:"build_command"`
		ReleaseCommand string `json:"build_command"`
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

func rpcCommand(server, cmd, proc string) error {
	client, err := rpc.Dial("tcp", server)
	if err != nil {
		return err
	}
	defer client.Close()
	var ret string
	switch cmd {
	case "start":
		return client.Call("Goreman.Start", proc, &ret)
	case "stop":
		return client.Call("Goreman.Stop", proc, &ret)
	case "restart":
		return client.Call("Goreman.Restart", proc, &ret)
	}
	return errors.New("Unknown command")
}

func runCommand(name, dir, command string) error {
	if command != "" {
		cs := []string{"/bin/bash", "-c", command}
		if runtime.GOOS == "windows" {
			cs = []string{"cmd", "/c", command}
		}
		cmd := exec.Command(cs[0], cs[1:]...)
		cmd.Dir = dir
		return cmd.Run()
	}
	return nil
}

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
		var succeeded = false
		if json.Unmarshal([]byte(r.FormValue("payload")), &p) == nil {
			name := p.Repository.Name
			app, ok := c.Apps[name]
			if ok {
				updateCommand := app.UpdateCommand
				if updateCommand == "" {
					updateCommand = "git pull origin master"
				}
				commands := []string {
					updateCommand,
					app.BuildCommand,
					app.ReleaseCommand,
				}

				if c.RPC != "" && app.Proc != "" {
					err = rpcCommand(c.RPC, "stop", app.Proc)
					if err != nil {
						log.Printf("%s: %s\n", name, err.Error())
						return
					}
				}
				for _, command := range commands {
					if command != "" {
						log.Printf("%s: %s\n", name, command)
						err = runCommand(name, app.Path, command)
						if err != nil {
							log.Printf("%s: %s\n", name, err.Error())
							http.Error(w, "Failed to update", http.StatusBadRequest)
							break
						} else {
							succeeded = true
						}
					}
				}
				if c.RPC != "" && app.Proc != "" {
					err = rpcCommand(c.RPC, "start", app.Proc)
					if err != nil {
						log.Printf("%s: %s\n", name, err.Error())
					}
				}
			}
		}
		if succeeded {
			fmt.Fprintf(w, "OK")
		}
	})
    http.ListenAndServe(c.Addr, nil)
}

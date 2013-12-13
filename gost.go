package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

type App struct {
	Proc           string `json:"proc"`
	Path           string `json:"path"`
	UpdateCommand  string `json:"update_command"`
	BuildCommand   string `json:"build_command"`
	TestCommand    string `json:"test_command"`
	ReleaseCommand string `json:"release_command"`
}

type config struct {
	Addr string         `json:"addr"`
	Root string         `json:"root"`
	RPC  string         `json:"rpc"`
	Apps map[string]App `json:"apps"`
	Log  string         `json:"log"`
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

func loadConfig() (*config, error) {
	f, err := os.Open(*configFile)
	if err != nil {
		return nil, err
	}

	var c config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}
	f.Close()

	if c.Root == "" || c.Root[0] != '/' {
		c.Root = "/" + c.Root
	}
	return &c, nil
}

func main() {
	flag.Parse()

	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if c.Log != "" {
		logFile, err := os.OpenFile(c.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile)
	}
	log.Println("start gost")

	if c.Root == "" || c.Root[0] != '/' {
		c.Root = "/" + c.Root
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP)
	go func() {
		for _ = range sc {
			log.Println("reloading configuration")
			ct, err := loadConfig()
			if err != nil {
				log.Println(err)
			} else {
				c = ct
			}
		}
	}()

	http.HandleFunc(c.Root, func(w http.ResponseWriter, r *http.Request) {
		var p payload
		var succeeded = false
		value := r.FormValue("payload")
		log.Println(value)
		err := json.Unmarshal([]byte(value), &p)
		if err != nil {
			log.Println(err)
		} else {
			name := p.Repository.Name
			app, ok := c.Apps[name]
			if ok {
				log.Printf("found app: ", name)
				updateCommand := app.UpdateCommand
				if updateCommand == "" {
					updateCommand = "git pull origin master"
				}
				commands := []struct {
					task    string
					command string
				}{
					{"update", updateCommand},
					{"build", app.BuildCommand},
					{"test", app.TestCommand},
					{"release", app.ReleaseCommand},
				}

				if c.RPC != "" && app.Proc != "" {
					log.Printf("stopping app: ", name)
					err = rpcCommand(c.RPC, "stop", app.Proc)
					if err != nil {
						log.Printf("%s: %s\n", name, err.Error())
						//return
					}
				}
				for _, command := range commands {
					if command.command != "" {
						log.Printf("%s: %s\n", name, command.command)
						err = runCommand(name, app.Path, command.command)
						if err != nil {
							log.Printf("%s: %s (%s)\n", name, command.task, err.Error())
							http.Error(w, fmt.Sprintf("Failed to %s", command.task), http.StatusBadRequest)
							break
						} else {
							succeeded = true
						}
					}
				}
				if c.RPC != "" && app.Proc != "" {
					log.Printf("starting app: ", name)
					err = rpcCommand(c.RPC, "start", app.Proc)
					if err != nil {
						log.Printf("%s: %s\n", name, err.Error())
					}
				}
			} else {
				log.Printf("invalid app name")
			}
		}
		if succeeded {
			fmt.Fprintf(w, "OK")
		} else {
			fmt.Fprintf(w, "NG")
		}
	})
	err = http.ListenAndServe(c.Addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

# gost

github push deployer 

## usage

    $ gost -c config.json

## setting

    
    {
        "addr": "127.0.0.1:10092",
        "rpc": "127.0.0.1:5555",
        "apps": {
            "my-app1": {
                "proc": "today_bot",
                "path": "/home/mattn/dev/today_bot/"
            },
            "my-app2": {
                "proc": "today_bot",
                "path": "/home/mattn/dev/today_bot/"
            }
        }
    }

* addr: server address
* rpc: goreman port
* my-app1: this mean github repository name
* proc: goreman proc name to start or stop
* path: path that is working the proc

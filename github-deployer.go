package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		io.Copy(os.Stdout, r.Body)
		fmt.Fprintf(w, "")
	})
    http.ListenAndServe(":10092", nil)
}

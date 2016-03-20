// file: ./goserv.go
// desc: open a simple local go server into input path.
// author: gouvinb
// legal: see LICENSE.md
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gouvinb/goserv/tools"
)

var (
	path = flag.String("path", "./", "which directory to create doc server")
	port = flag.String("port", "4000", "port of server")
)

func main() {
	flag.Parse()

	tools.IsCorrectPort(*port)
	tools.CheckPath(*path)

	go func() {
		url := "http://localhost:"
		if tools.WaitServer(url, *port) && tools.StartBrowser(url, *port) {
			log.Printf("A browser window should open. If not, please visit %s", url+*port)
		} else {
			log.Printf("Please open your web browser and visit %s", url+*port)
		}
	}()

	err := http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*path)))
	for err != nil {
		if err.Error() == "listen tcp :"+*port+": bind: address already in use" {
			tools.IncrementPort(port)
			err = http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*path)))
		} else {
			log.Fatal(err)
		}
	}
}

// file: ./goserv.go
// desc: open a simple local go server into input path.
// author: gouvinb
// legal: see LICENSE.md
package main

import (
	// std
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

var (
	path = flag.String("path", "./", "which directory to create doc server")
	port = flag.String("port", "4000", "port of server")
)

func main() {
	flag.Parse()

	isCorrectPort(*port)
	checkPath(*path)

	go func() {
		url := "http://localhost:"
		if waitServer(url) && startBrowser(url) {
			log.Printf("A browser window should open. If not, please visit %s", url+*port)
		} else {
			log.Printf("Please open your web browser and visit %s", url+*port)
		}
	}()

	err := http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*path)))
	for err != nil {
		if err.Error() == "listen tcp :"+*port+": bind: address already in use" {
			incrementPort(port)
			err = http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*path)))
		} else {
			log.Fatal(err)
		}
	}
}

// isCorrectPort verify if the port is a valide number.
func isCorrectPort(port string) {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("port : ", portNum)
	}
}

// checkPath verify if the directory exist.
func checkPath(path string) {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("path : ", path)
	}
}

// incrementPort change *port if the port is already used.
func incrementPort(port *string) {
	portNum, err := strconv.Atoi(*port)
	if err != nil {
		log.Fatal(err)
	} else {
		portNum++
		log.Print("port : ", portNum)
		*port = strconv.Itoa(portNum)
	}
}

// waitServer waits some time for the http Server to start
// serving url. The return value reports whether it starts.
func waitServer(url string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url + *port)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	log.Print("url : ", url, *port)
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url+*port)...)
	return cmd.Start() == nil
}

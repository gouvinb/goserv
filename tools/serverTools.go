package tools

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// IsCorrectPort verify if the port is a valide number.
func IsCorrectPort(port string) {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	} else {
		ip, err := ExternalIP()
		if err != nil {
			log.Panic(err)
		}
		log.Print("ip : ", ip)
		log.Print("port : ", portNum)
		log.Print("result : ", ip, ":", portNum)
	}
}

// ExternalIP return current ip
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("network is disable")
}

// checkPath verify if the directory exist.
func CheckPath(path string) {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("path : ", path)
	}
}

// IncrementPort change *port if the port is already used.
func IncrementPort(port *string) {
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
func WaitServer(url string, port string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url + port)
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
func StartBrowser(url string, port string) bool {
	// try to start the browser
	log.Print("url : ", url, port)
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url+port)...)
	return cmd.Start() == nil
}

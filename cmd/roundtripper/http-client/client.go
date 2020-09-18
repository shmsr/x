package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	rtick = 1 // request ticker (in s)
	ctick = 5 // cache clear ticker (in s)
)

const (
	flagHost   = "host"
	flagPort   = "port"
	flagScheme = "scheme"
)

const (
	defaultHost   = "127.0.0.1"
	defaultPort   = "8080"
	defaultScheme = "http"
)

const (
	usageHost   = "enter host"
	usagePort   = "enter port"
	usageScheme = "enter scheme"
)

var (
	host   string
	port   string
	scheme string
)

func responseHandle(client *http.Client, req *http.Request) (string, error) {
	resp, err := client.Do(req)
	if err != nil && err.Error() != errInitCache {
		return "", err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func main() {
	flag.StringVar(&host, flagHost, defaultHost, usageHost)
	flag.StringVar(&port, flagPort, defaultPort, usagePort)
	flag.StringVar(&scheme, flagScheme, defaultScheme, usageScheme)
	flag.Parse()

	csize := int64(1024)
	tr := getCacheTransport(csize)
	client := &http.Client{
		Transport: tr,
	}

	url := net.JoinHostPort(host, port)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s", scheme, url), strings.NewReader(""))
	if err != nil {
		log.Fatalln(err)
	}

	sterm := make(chan os.Signal, 1)
	signal.Notify(sterm, syscall.SIGHUP, syscall.SIGTERM)

	var (
		rticker = time.NewTicker(rtick * time.Second)
		cticker = time.NewTicker(ctick * time.Second)
	)

	for {
		select {
		case <-rticker.C:
			resp, err := responseHandle(client, req)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(resp)
		case <-cticker.C:
			tr.Clear()
		case <-sterm:
			rticker.Stop()
			cticker.Stop()
			return
		}
	}
}

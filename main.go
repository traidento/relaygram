package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var nekoXProxyString string

var nekoXProxyBaseDomain string
var nekoXProxyDomains []string

var client *http.Client

func main() {
	listen := flag.String("l", "127.0.0.1:26641", "HttpProxy listen port")
	_nekoXProxyString := flag.String("p", "", "NekoX Proxy URL")
	flag.Parse()

	var ok bool
	if *_nekoXProxyString != "" {
		ok = parseNekoXString(*_nekoXProxyString)
	}

	if !ok {
		log.Println("NekoX Proxy URL is required")
		return
	}

	client = &http.Client{
		Timeout: 25 * 2 * time.Second,
	}

	http.HandleFunc("/", relay)
	server := &http.Server{
		Addr:         *listen,
		WriteTimeout: 25 * 2 * time.Second,
		ReadTimeout:  25 * 2 * time.Second,
	}

	log.Println("Telegram HTTP Proxy started at", *listen)
	server.ListenAndServe()
}

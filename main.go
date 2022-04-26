package main

import (
	"encoding/base64"
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
	_nekoXProxyString := flag.String("p", "", "NekoX Proxy URL (keep empty if you don't know)")
	flag.Parse()

	var ok bool
	if *_nekoXProxyString != "" {
		ok = parseNekoXString(base64.RawURLEncoding.EncodeToString([]byte(*_nekoXProxyString)))
	}

	if !ok {
		log.Println("Failed to parse NekoX proxy.")
		return
	}

	client = &http.Client{}

	http.HandleFunc("/", relay)
	server := &http.Server{
		Addr:         *listen,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Telegram HTTP Proxy started at", *listen)
	server.ListenAndServe()
}

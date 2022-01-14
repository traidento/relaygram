package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var proxyBaseDomain string
var proxyDomains []string

var client *http.Client

func main() {

	listen := *flag.String("l", getEnv("LISTEN", "127.0.0.1:26641"), "HttpProxy listen addr")
	proxy := *flag.String("p", getEnv("PROXY", ""), "Proxy URL (keep empty if you don't know)")
	flag.Parse()

	var ok bool
	if proxy != "" {
		ok = parseRelayProxy(base64.RawURLEncoding.EncodeToString([]byte(proxy)))
	} else {
		log.Println("Getting public proxy...")
		ok = parseRelayProxy(getPublicRelay())
	}

	if !ok {
		log.Println("Failed to parse proxy.")
		return
	}

	client = &http.Client{}

	http.HandleFunc("/", relay)
	server := &http.Server{
		Addr:         listen,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("HTTP Proxy started at", listen)
	server.ListenAndServe()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

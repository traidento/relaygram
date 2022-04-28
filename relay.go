package main

import (
	"io"
	"log"
	"net/http"
)

func relay(w http.ResponseWriter, r *http.Request) {
	reqip := r.URL.Hostname()
	wsurl, found := ip2wsurl(reqip)
	if !found {
		log.Println("New DC address found:", reqip)
		return
	}

	req, _ := http.NewRequest(r.Method, wsurl, r.Body)
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(502)
		log.Println(502, err)
		return
	}

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

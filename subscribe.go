package main

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var nekoXSubscriptionDomain = "_tg.chapro.xyz"
var subscriptionDoh = []string{
	"https://1.1.1.1/dns-query",
	"https://1.0.0.1/dns-query",
	"https://101.101.101.101/dns-query",
	"https://8.8.8.8/resolve",
	"https://8.8.4.4/resolve",
	"https://[2606:4700:4700::1111]/dns-query",
}

var _subscribeGood int32 = 0
var _subscribeBad int32 = 0

func getPublicRelay() string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	in := make(chan string, len(subscriptionDoh))
	out := make(chan string)

	for i := 0; i < 10; i++ {
		go getPublicProxyWorker(ctx, in, out, cancel)
	}

	go func() {
		for _, doh := range subscriptionDoh {
			in <- doh
		}
	}()

	defer cancel()
	return <-out
}

func getPublicProxyWorker(ctx context.Context, in, out chan string, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case doh := <-in:
			ret := getTXTUsingDoH(ctx, doh)

			_, err1 := base64.StdEncoding.DecodeString(ret)
			_, err2 := base64.URLEncoding.DecodeString(ret)

			if ret == "" || (err1 != nil && err2 != nil) {
				if atomic.AddInt32(&_subscribeBad, 1) == int32(len(subscriptionDoh)) {
					cancel()
					out <- ""
					return
				}
				continue
			}
			if atomic.AddInt32(&_subscribeGood, 1) == 1 {
				cancel()
				out <- ret
			}
		}
	}
}

func getTXTUsingDoH(ctx context.Context, doh string) string {
	dohURL := doh + "?name=" + nekoXSubscriptionDomain + "&type=TXT"
	req, _ := http.NewRequestWithContext(ctx, "GET", dohURL, nil)
	req.Header.Set("accept", "application/dns-json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	data, _ := ioutil.ReadAll(resp.Body)
	data2 := strings.ReplaceAll(string(data), "\\\"", "")
	data2 = strings.ReplaceAll(data2, " ", "")
	return between(data2, "#NekoXStart#", "#NekoXEnd#")
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

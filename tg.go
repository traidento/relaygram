package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
)

var mapper = make(map[string]int)

func parseRelayProxy(a string) bool {
	if a == "" {
		return false
	}

	b, _ := base64.StdEncoding.DecodeString(a)
	url, err := url.Parse(string(b))
	if err != nil {
		return false
	}

	pldstr, _ := base64.RawURLEncoding.DecodeString(url.Query().Get("payload"))
	plds := strings.Split(string(pldstr), ",")

	proxyBaseDomain = url.Host
	proxyDomains = append([]string{""}, plds...)

	// fmt.Println(url.String())

	putmapper := func(ip string, dc int) {
		mapper[ip] = dc
	}

	putmapper("149.154.175.5", 1)
	putmapper("95.161.76.100", 2)
	putmapper("149.154.175.100", 3)
	putmapper("149.154.167.91", 4)
	putmapper("149.154.167.92", 4)
	putmapper("149.154.171.", 5)
	putmapper("2001:b28:f23d:f001::", 1)
	putmapper("2001:67c:4e8:f002::", 2)
	putmapper("2001:b28:f23d:f003::", 3)
	putmapper("2001:67c:4e8:f004::", 4)
	putmapper("2001:b28:f23f:f005::", 5)
	putmapper("149.154.161.144", 2)
	putmapper("149.154.167.", 2)
	putmapper("149.154.175.1", 3)
	putmapper("91.108.4.", 4)
	putmapper("149.154.164.", 4)
	putmapper("149.154.165.", 4)
	putmapper("149.154.166.", 4)
	putmapper("91.108.56.", 5)
	putmapper("149.154.175.40", 6)
	putmapper("149.154.167.40", 7)
	putmapper("149.154.175.117", 8)
	putmapper("2001:b28:f23d:f001::e", 6)
	putmapper("2001:67c:4e8:f002::e", 7)
	putmapper("2001:b28:f23d:f003::e", 8)

	return true
}

func dc2wsurl(dc int) string {
	if dc == 0 {
		return ""
	}
	return fmt.Sprintf("https://%s.%s/api", proxyDomains[dc], proxyBaseDomain)
}

func ip2dc(ip string) int {
	for k, v := range mapper {
		if ip == k {
			return v
		}
	}
	for k, v := range mapper {
		if strings.HasPrefix(ip, k) {
			return v
		}
	}
	log.Printf("missing: %v", ip)

	return 0
}

func ip2wsurl(ip string) string {
	return dc2wsurl(ip2dc(ip))
}

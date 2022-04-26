package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

var mapper = make(map[string]int)

func parseNekoXString(a string) bool {
	// fmt.Println(a)

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

	nekoXProxyBaseDomain = url.Host
	nekoXProxyDomains = append([]string{""}, plds...)

	// fmt.Println(url.String())

	putmapper := func(ip string, dc int) {
		mapper[ip] = dc
	}

	putmapper("149.154.175.5", 1)
	putmapper("95.161.76.100", 2)
	putmapper("149.154.175.100", 3)
	putmapper("149.154.167.91", 4)
	putmapper("149.154.167.92", 4)
	putmapper("149.154.171.5", 5)
	putmapper("2001:b28:f23d:f001::a", 1)
	putmapper("2001:67c:4e8:f002::a", 2)
	putmapper("2001:b28:f23d:f003::a", 3)
	putmapper("2001:67c:4e8:f004::a", 4)
	putmapper("2001:b28:f23f:f005::a", 5)
	putmapper("149.154.161.144", 2)
	putmapper("149.154.167.", 2)
	putmapper("149.154.175.1", 3)
	putmapper("91.108.4.", 4)
	putmapper("149.154.164.", 4)
	putmapper("149.154.165.", 4)
	putmapper("149.154.166.", 4)
	putmapper("91.108.56.", 5)
	putmapper("2001:b28:f23d:f001::d", 1)
	putmapper("2001:67c:4e8:f002::d", 2)
	putmapper("2001:b28:f23d:f003::d", 3)
	putmapper("2001:67c:4e8:f004::d", 4)
	putmapper("2001:b28:f23f:f005::d", 5)
	putmapper("149.154.175.40", 6)
	putmapper("149.154.167.40", 7)
	putmapper("149.154.175.117", 8)
	putmapper("2001:b28:f23d:f001::e", 6)
	putmapper("2001:67c:4e8:f002::e", 7)
	putmapper("2001:b28:f23d:f003::e", 8)

	putmapper("2001:b28:f23d:f001::b", 1)
	putmapper("2001:67c:4e8:f002::b", 2)
	putmapper("2001:b28:f23d:f003::b", 3)
	putmapper("2001:67c:4e8:f004::b", 4)
	putmapper("2001:b28:f23f:f005::b", 5)

	return true
}

func ip2dc(ip string) (int, bool) {
	for k, v := range mapper {
		if ip == k {
			return v, true
		}
	}
	for k, v := range mapper {
		if strings.HasPrefix(ip, k) {
			return v, true
		}
	}
	return 0, false
}

func ip2wsurl(ip string) (string, bool) {
	dc, found := ip2dc(ip)
	if !found {
		return "", false
	}
	url := fmt.Sprintf("https://%s.%s/api", nekoXProxyDomains[dc], nekoXProxyBaseDomain)

	return url, true
}

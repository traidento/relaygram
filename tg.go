package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/asergeyev/nradix"
)

var mapper = nradix.NewTree(0)

func parseNekoXString(a string) bool {
	// fmt.Println(a)

	if a == "" {
		return false
	}

	url, err := url.Parse(a)
	if err != nil {
		return false
	}

	pldstr, _ := base64.RawURLEncoding.DecodeString(url.Query().Get("payload"))
	nekoXProxyDomains = strings.Split(string(pldstr), ",")

	nekoXProxyBaseDomain = url.Host

	mapper.AddCIDR("149.154.175.5", 1)
	mapper.AddCIDR("95.161.76.100", 2)
	mapper.AddCIDR("149.154.175.100", 3)
	mapper.AddCIDR("149.154.167.91", 4)
	mapper.AddCIDR("149.154.167.92", 4)
	mapper.AddCIDR("149.154.171.5", 5)
	mapper.AddCIDR("2001:b28:f23d:f001::a", 1)
	mapper.AddCIDR("2001:67c:4e8:f002::a", 2)
	mapper.AddCIDR("2001:b28:f23d:f003::a", 3)
	mapper.AddCIDR("2001:67c:4e8:f004::a", 4)
	mapper.AddCIDR("2001:b28:f23f:f005::a", 5)
	mapper.AddCIDR("149.154.161.144", 2)
	mapper.AddCIDR("149.154.167.0/24", 2)
	mapper.AddCIDR("149.154.175.1", 3)
	mapper.AddCIDR("91.108.4.0/24", 4)
	mapper.AddCIDR("149.154.164.0/24", 4)
	mapper.AddCIDR("149.154.165.0/24", 4)
	mapper.AddCIDR("149.154.166.0/24", 4)
	mapper.AddCIDR("91.108.56.0/24", 5)
	mapper.AddCIDR("2001:b28:f23d:f001::d", 1)
	mapper.AddCIDR("2001:67c:4e8:f002::d", 2)
	mapper.AddCIDR("2001:b28:f23d:f003::d", 3)
	mapper.AddCIDR("2001:67c:4e8:f004::d", 4)
	mapper.AddCIDR("2001:b28:f23f:f005::d", 5)
	mapper.AddCIDR("149.154.175.40", 6)
	mapper.AddCIDR("149.154.167.40", 7)
	mapper.AddCIDR("149.154.175.117", 8)
	mapper.AddCIDR("2001:b28:f23d:f001::e", 6)
	mapper.AddCIDR("2001:67c:4e8:f002::e", 7)
	mapper.AddCIDR("2001:b28:f23d:f003::e", 8)

	mapper.AddCIDR("2001:b28:f23d:f001::b", 1)
	mapper.AddCIDR("2001:67c:4e8:f002::b", 2)
	mapper.AddCIDR("2001:b28:f23d:f003::b", 3)
	mapper.AddCIDR("2001:67c:4e8:f004::b", 4)
	mapper.AddCIDR("2001:b28:f23f:f005::b", 5)

	mapper.AddCIDR("149.154.175.55", 1)

	return true
}

func ip2dc(ip string) (int, bool) {
	dc, err := mapper.FindCIDR(ip)
	if dc == nil || err != nil {
		return 0, false
	} else {
		return dc.(int), true
	}
}

func ip2wsurl(ip string) (string, bool) {
	dc, found := ip2dc(ip)
	if !found {
		return "", false
	}
	url := fmt.Sprintf("https://%s.%s/api", nekoXProxyDomains[dc-1], nekoXProxyBaseDomain)

	return url, true
}

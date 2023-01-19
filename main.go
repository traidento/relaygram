package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/asergeyev/nradix"
	"nhooyr.io/websocket"
)

func main() {
	listen := flag.String("l", "127.0.0.1:26641", "listen address")
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		fmt.Fprintln(os.Stderr, "NekoX Proxy URL is required")
		return
	}

	router, err := NewRouter(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if network != "tcp" {
			return nil, fmt.Errorf("tcp only: %s", network)
		}
		wsurl, found := router.IP2URL(addr)
		if !found {
			return nil, fmt.Errorf("New DC address found: %s", addr)
		}

		dial_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		c, _, err := websocket.Dial(dial_ctx, wsurl,
			&websocket.DialOptions{
				Subprotocols: []string{"binary"},
				// disable compression for faster media download
				CompressionMode: websocket.CompressionDisabled,
			})
		if err != nil {
			return nil, err
		}
		c.SetReadLimit(math.MaxInt64 - 1) // disable the buggy limitReader
		return websocket.NetConn(ctx, c, websocket.MessageBinary), nil
	}

	server := Server{Dialer: dialer}

	l, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Println("Telegram Socks5 Proxy started at", *listen)
	log.Fatal(server.Serve(l))
}

type Router struct {
	baseDomain      string
	subDomainMapper []string
	dcMapper        *nradix.Tree
}

func NewRouter(nekoXURL string) (*Router, error) {
	url, err := url.ParseRequestURI(nekoXURL)
	if err != nil {
		return nil, err
	}

	pldstr, _ := base64.RawURLEncoding.DecodeString(url.Query().Get("payload"))

	router := Router{
		baseDomain:      url.Host,
		subDomainMapper: strings.Split(string(pldstr), ","),
		dcMapper:        nradix.NewTree(0),
	}

	router.initDcMapper()
	return &router, nil
}

func (router *Router) initDcMapper() {
	router.dcMapper.AddCIDR("149.154.175.5", 1)
	router.dcMapper.AddCIDR("95.161.76.100", 2)
	router.dcMapper.AddCIDR("149.154.175.100", 3)
	router.dcMapper.AddCIDR("149.154.167.91", 4)
	router.dcMapper.AddCIDR("149.154.167.92", 4)
	router.dcMapper.AddCIDR("149.154.171.5", 5)
	router.dcMapper.AddCIDR("2001:b28:f23d:f001::a", 1)
	router.dcMapper.AddCIDR("2001:67c:4e8:f002::a", 2)
	router.dcMapper.AddCIDR("2001:b28:f23d:f003::a", 3)
	router.dcMapper.AddCIDR("2001:67c:4e8:f004::a", 4)
	router.dcMapper.AddCIDR("2001:b28:f23f:f005::a", 5)
	router.dcMapper.AddCIDR("149.154.161.144", 2)
	router.dcMapper.AddCIDR("149.154.167.0/24", 2)
	router.dcMapper.AddCIDR("149.154.175.1", 3)
	router.dcMapper.AddCIDR("91.108.4.0/24", 4)
	router.dcMapper.AddCIDR("149.154.164.0/24", 4)
	router.dcMapper.AddCIDR("149.154.165.0/24", 4)
	router.dcMapper.AddCIDR("149.154.166.0/24", 4)
	router.dcMapper.AddCIDR("91.108.56.0/24", 5)
	router.dcMapper.AddCIDR("2001:b28:f23d:f001::d", 1)
	router.dcMapper.AddCIDR("2001:67c:4e8:f002::d", 2)
	router.dcMapper.AddCIDR("2001:b28:f23d:f003::d", 3)
	router.dcMapper.AddCIDR("2001:67c:4e8:f004::d", 4)
	router.dcMapper.AddCIDR("2001:b28:f23f:f005::d", 5)
	router.dcMapper.AddCIDR("149.154.175.40", 6)
	router.dcMapper.AddCIDR("149.154.167.40", 7)
	router.dcMapper.AddCIDR("149.154.175.117", 8)
	router.dcMapper.AddCIDR("2001:b28:f23d:f001::e", 6)
	router.dcMapper.AddCIDR("2001:67c:4e8:f002::e", 7)
	router.dcMapper.AddCIDR("2001:b28:f23d:f003::e", 8)

	router.dcMapper.AddCIDR("2001:b28:f23d:f001::b", 1)
	router.dcMapper.AddCIDR("2001:67c:4e8:f002::b", 2)
	router.dcMapper.AddCIDR("2001:b28:f23d:f003::b", 3)
	router.dcMapper.AddCIDR("2001:67c:4e8:f004::b", 4)
	router.dcMapper.AddCIDR("2001:b28:f23f:f005::b", 5)

	router.dcMapper.AddCIDR("149.154.175.55", 1)
}

func (router *Router) IP2URL(addr string) (string, bool) {
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err.Error())
	}
	dc, err := router.dcMapper.FindCIDR(ip)
	if dc == nil || err != nil {
		return "", false
	}
	url := fmt.Sprintf("wss://%s.%s/api", router.subDomainMapper[dc.(int)-1], router.baseDomain)
	log.Printf("new connection to DC%d (%s)\n", dc, addr)

	return url, true
}

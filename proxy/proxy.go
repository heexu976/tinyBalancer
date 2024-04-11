/*
 * @Date: 2024-04-10 21:10:58
 * @LastEditors: HeXu
 * @LastEditTime: 2024-04-11 10:36:45
 * @FilePath: /tinyBalancer/proxy/proxy.go
 */
package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"log"

	"github.com/heexu976/tinybalancer/balancer"
)

var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
)

var (
	ReverseProxy = "Balance-Reverse-Proxy"
)

type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer
	sync.RWMutex
	alive map[string]bool
}

func NewHTTPProxy(targetHosts []string, algorithm string) (*HTTPProxy, error) {
	hosts := make([]string, 0)
	hostMap := make(map[string]*httputil.ReverseProxy)
	alive := make(map[string]bool)
	for _, targetHost := range targetHosts {
		url, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		originDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originDirector(req)
			req.Header.Set(XProxy, ReverseProxy)
			req.Header.Set(XRealIP, GetIP(req))
		}
		host := GetHost(url)
		alive[host] = true
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}
	lb, err := balancer.Build(algorithm, hosts)
	if err != nil {
		return nil, err
	}
	return &HTTPProxy{
		hostMap: hostMap,
		lb:      lb,
		alive:   alive,
	}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy causes panic %s", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.(error).Error()))
		}
	}()
	host, err := h.lb.Balance(GetIP(req))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("balance error")))
		return
	}
	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.hostMap[host].ServeHTTP(w, req)
}

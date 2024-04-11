/*
 * @Date: 2024-04-10 21:10:50
 * @LastEditors: HeXu
 * @LastEditTime: 2024-04-11 09:51:42
 * @FilePath: /tinyBalancer/proxy/net.go
 */
package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ConnectionTimeout = 3 * time.Second

func GetIP(req *http.Request) string {
	clientIP, _, _ := net.SplitHostPort(req.RemoteAddr)
	if len(req.Header.Get(XForwardedFor)) != 0 {
		xff := req.Header.Get(XForwardedFor)
		s := strings.Index(xff, ", ")
		if s == -1 {
			s = len(req.Header.Get(XForwardedFor))
		}
		clientIP = xff[:s]
	} else if len(req.Header.Get(XRealIP)) != 0 {
		clientIP = req.Header.Get(XRealIP)
	}
	return clientIP
}

func GetHost(url *url.URL) string {
	if _, _, err := net.SplitHostPort(url.Host); err == nil {
		return url.Host
	}
	if url.Scheme == "http" {
		return fmt.Sprintf("%s:%s", url.Host, "80")
	} else if url.Scheme == "https" {
		return fmt.Sprintf("%s:%s", url.Host, "443")
	}
	return url.Host
}

func IsBackendAlive(host string) bool {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return false
	}
	resolveAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	conn, err := net.DialTimeout("tcp", resolveAddr, ConnectionTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

/*
 * @Date: 2024-04-10 21:12:04
 * @LastEditors: HeXu
 * @LastEditTime: 2024-04-11 11:19:54
 * @FilePath: /tinyBalancer/proxy/health_check.go
 */
package proxy

import (
	"log"
	"time"
)

func (h *HTTPProxy) ReadAlive(url string) bool {
	h.RLock()
	defer h.RUnlock()
	return h.alive[url]
}

func (h *HTTPProxy) SetAlive(url string, alive bool) {
	h.Lock()
	defer h.Unlock()
	h.alive[url] = true
}

func (h *HTTPProxy) HealthCheck(interval uint) {
	for host := range h.hostMap {
		go h.healthCheck(host, interval)
	}
}

func (h *HTTPProxy) healthCheck(host string, interval uint) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		if !IsBackendAlive(host) && h.ReadAlive(host) {
			log.Printf("Site unreachable, remove %s from load balancer.", host)
			h.SetAlive(host, false)
			h.lb.Remove(host)
		} else if IsBackendAlive(host) && !h.ReadAlive(host) {
			log.Printf("Site reachable, add %s to load balancer.", host)
			h.SetAlive(host, true)
			h.lb.Add(host)
		}
	}

}

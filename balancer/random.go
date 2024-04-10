/*
 * @Date: 2024-04-10 20:33:52
 * @LastEditors: HeXu
 * @LastEditTime: 2024-04-10 20:42:11
 * @FilePath: /tinyBalancer/balancer/random.go
 */
package balancer

import (
	"math/rand"
	"time"
)

func init() {
	factories[RandomBalancer] = NewRandom
}

type Random struct {
	BaseBalancer
	rnd *rand.Rand
}

func NewRandom(hosts []string) Balancer {
	return &Random{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r Random) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", ErrNoHost
	}
	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}

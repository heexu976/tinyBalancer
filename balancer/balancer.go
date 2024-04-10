/*
 * @Date: 2024-04-10 20:14:35
 * @LastEditors: HeXu
 * @LastEditTime: 2024-04-10 20:25:58
 * @FilePath: /tinyBalancer/balancer/balancer.go
 */
package balancer

import (
	"errors"
)

var (
	ErrNoHost                = errors.New("no host")
	ErrAlgorithmNotSupported = errors.New("algorithm not supported")
)

type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}

// Factory is the factory that generates Balancer,
// and the factory design pattern is used here
type Factory func([]string) Balancer

var factories = make(map[string]Factory)

func Build(algorithm string, hosts []string) (Balancer, error) {
	factory, ok := factories[algorithm]
	if !ok {
		return nil, ErrAlgorithmNotSupported
	}
	return factory(hosts), nil
}

/*
 * @Date: 2024-04-10 20:29:55
 * @LastEditors: HeXu
 * @LastEditTime: 2024-04-10 20:33:37
 * @FilePath: /tinyBalancer/balancer/const.go
 */
package balancer

const (
	IPHashBalancer         = "ip-hash"
	ConsistentHashBalancer = "consistent-hash"
	P2CBalancer            = "p2c"
	RandomBalancer         = "random"
	R2Balancer             = "round-robin"
	LeastLoadBalancer      = "least-load"
	BoundedBalancer        = "bounded"
)

/*
 * @Date: 2024-04-11 14:49:20
 * @LastEditors: HeXu
 * @LastEditTime: 2024-04-11 15:35:53
 * @FilePath: /tinyBalancer/middleware.go
 */
package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func maxAllowedMiddleware(n uint) mux.MiddlewareFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acquire()
			defer release()
			next.ServeHTTP(w, r)
		})
	}
}

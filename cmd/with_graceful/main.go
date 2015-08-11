package main

import (
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/freeformz/grb/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/freeformz/grb/graceful"
)

func nothing(c martini.Context) {
	c.Next()
}
func main() {
	m := martini.Classic()
	gracefulShutdown := &graceful.GracefulShutdown{Timeout: time.Duration(20) * time.Second}
	m.Use(gracefulShutdown.Handler)
	m.Post("/signup", nothing)

	go func() {
		if err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", m); err != nil {
			fmt.Print(err)
		}
	}()
	err := gracefulShutdown.WaitForSignal(syscall.SIGTERM, syscall.SIGINT)
	if err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"context"
	"fmt"
	"go-distributed/registry"
	"log"
	"net/http"
)

func main() {
	http.Handle("/services", &registry.Service{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Println("Registry service started, press any key to stop.")
		var s string
		_, _ = fmt.Scanln(&s)
		_ = srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down registry service.")

}

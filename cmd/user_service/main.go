package main

import (
	"context"
	"fmt"
	"go-distributed/registry"
	"go-distributed/service"
	"go-distributed/user"
	syslog "log"
)

func main() {
	host, port := "localhost", "8080"
	serviceURL := fmt.Sprintf("http://%s:%v", host, port)

	r := registry.Registration{
		ServiceName: registry.UserServiceName,
		ServiceURL:  serviceURL,
	}

	ctx, err := service.Start(
		context.Background(), host, port, r, user.RegisterHandler,
	)

	if err != nil {
		syslog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down user service.")
}

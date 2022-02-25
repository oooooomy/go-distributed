package main

import (
	"context"
	"fmt"
	"go-distributed/log"
	"go-distributed/registry"
	"go-distributed/service"
	syslog "log"
)

func main() {
	log.Run("./system.log")
	host, port := "localhost", "4000"
	serviceURL := fmt.Sprintf("http://%s:%v", host, port)

	r := registry.Registration{
		ServiceName:      registry.LogServiceName,
		ServiceURL:       serviceURL,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceURL + "/services",
	}

	ctx, err := service.Start(
		context.Background(), host, port, r, log.RegisterHandler,
	)

	if err != nil {
		syslog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service.")
}

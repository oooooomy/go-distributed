package main

import (
	"context"
	"fmt"
	"go-distributed/log"
	"go-distributed/registry"
	"go-distributed/service"
	"go-distributed/user"
	syslog "log"
)

func main() {
	host, port := "localhost", "8080"
	serviceURL := fmt.Sprintf("http://%s:%v", host, port)

	r := registry.Registration{
		ServiceName:      registry.UserServiceName,
		ServiceURL:       serviceURL,
		RequiredServices: []registry.ServiceName{registry.LogServiceName},
		ServiceUpdateURL: serviceURL + "/services",
		HeartbeatURL:     serviceURL + "/heart",
	}

	ctx, err := service.Start(
		context.Background(), host, port, r, user.RegisterHandler,
	)

	if err != nil {
		syslog.Fatalln(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogServiceName); err == nil {
		fmt.Println("Logging service found at ", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down user service.")
}

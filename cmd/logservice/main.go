package main

import (
	"context"
	"fmt"
	"go-distributed/log"
	"go-distributed/service"
	syslog "log"
)

func main() {
	log.Run("./system.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(
		context.Background(), "Log Service", host, port, log.RegisterHandler,
	)

	if err != nil {
		syslog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service.")
}

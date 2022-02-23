package service

import (
	"context"
	"fmt"
	"go-distributed/registry"
	"log"
	"net/http"
)

// Start 启动服务
func Start(ctx context.Context, host, port string, registration registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {

	//启动客户端的web服务
	registerHandlersFunc()
	ctx = startService(ctx, registration.ServiceName, host, port)
	//客户端将自己注册到注册中心
	err := registry.RegisterService(registration)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		//http 启动错误 取消 ctx
		cancel()
	}()

	go func() {
		fmt.Printf("%v started, press any key to stop. \n", serviceName)
		var s string
		_, _ = fmt.Scanln(&s)
		_ = srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}

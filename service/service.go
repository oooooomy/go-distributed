package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Start 启动服务
func Start(ctx context.Context, serviceName, host, port string, registerHandlersFunc func()) (context.Context, error) {
	//执行注册函数
	registerHandlersFunc()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
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

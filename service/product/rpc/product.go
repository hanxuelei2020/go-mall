package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"

	"go-mall/service/product/rpc/internal/config"
	"go-mall/service/product/rpc/internal/server"
	"go-mall/service/product/rpc/internal/svc"
	"go-mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile *string

// linux 不会出现此问题, 但是 win goland 中会出现此类问题
func init() {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	etcFile, _ := filepath.Abs(dir + "/etc/product.yaml")
	configFile = flag.String("f", etcFile, "the config file")
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		product.RegisterProductServer(grpcServer, server.NewProductServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

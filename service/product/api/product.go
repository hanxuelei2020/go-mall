package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"

	"go-mall/service/product/api/internal/config"
	"go-mall/service/product/api/internal/handler"
	"go-mall/service/product/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
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

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

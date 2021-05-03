package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

func middleware(name string, h fasthttp.RequestHandler, fh *xray.FastHTTPHandler) fasthttp.RequestHandler {
	f := func(ctx *fasthttp.RequestCtx) {
		h(ctx)
	}

	return fh.Handler(xray.NewFixedSegmentNamer(name), f)
}

func init() {
	if err := xray.Configure(xray.Config{
		DaemonAddr:     "xray:2000",
		ServiceVersion: "0.1",
	}); err != nil {
		panic(err)
	}

	xray.SetLogger(xraylog.NewDefaultLogger(os.Stdout, xraylog.LogLevelDebug))
}

func main() {
	fh := xray.NewFastHTTP(nil)
	r := router.New()
	r.GET("/", middleware("index", index, fh))
	r.GET("/hello/{name}", middleware("hello", hello, fh))

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}

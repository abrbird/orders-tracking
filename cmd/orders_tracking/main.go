package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"

	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/db"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/repository/sql_repository"
	"gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/service/implemented_service"
	wrkr "gitlab.ozon.dev/zBlur/homework-3/orders-tracking/internal/worker"
)

type GoroutinesNum struct{}

func (g *GoroutinesNum) String() string {
	return strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
}

func main() {
	cfg, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	dbConnPool, err := db.New(cfg.Database, ctx)
	if err != nil {
		log.Fatal(err)
	}

	//tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	//defer closer.Close()
	//opentracing.SetGlobalTracer(
	//	tracer,
	//)

	repository := sql_repository.New(dbConnPool)
	service := implemented_service.New()

	fmt.Println("Start orders...")
	fmt.Println("config", cfg)
	fmt.Println("repository", repository)
	fmt.Println("service", service)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker, err := wrkr.New(cfg, repository, service)
	if err != nil {
		log.Fatal(err)
	}
	err = worker.StartConsuming(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		g := &GoroutinesNum{}
		expvar.Publish("GoroutinesNum", g)
		fmt.Println("serving pprof")
		http.ListenAndServe("127.0.0.1:7999", nil)
	}()

	<-ctx.Done()
}

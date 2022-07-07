package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/abrbird/orders-tracking/internal/cache/redis_cache"
	"github.com/abrbird/orders-tracking/internal/metrics/prom_metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"

	"github.com/abrbird/orders-tracking/config"
	"github.com/abrbird/orders-tracking/internal/db"
	"github.com/abrbird/orders-tracking/internal/repository/sql_repository"
	"github.com/abrbird/orders-tracking/internal/service/implemented_service"
	wrkr "github.com/abrbird/orders-tracking/internal/worker"
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

	cache := redis_cache.New(cfg.Cache.Redis)

	//tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	//defer closer.Close()
	//opentracing.SetGlobalTracer(
	//	tracer,
	//)

	repository := sql_repository.New(dbConnPool)
	service := implemented_service.New(cache)
	metrics := prom_metrics.New(cfg)

	log.Printf("Start %s...", cfg.Application.Name)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker, err := wrkr.New(cfg, repository, service, metrics)
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

		err = http.ListenAndServe(
			fmt.Sprintf("%s:%v", cfg.Monitoring.Pprof.Host, cfg.Monitoring.Pprof.Port),
			nil,
		)
		if err != nil {
			log.Print(err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(
			fmt.Sprintf("%s:%v", cfg.Monitoring.Metrics.Host, cfg.Monitoring.Metrics.Port),
			nil,
		)
		if err != nil {
			log.Print(err)
		}
	}()

	<-ctx.Done()
}

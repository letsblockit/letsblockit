package server

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/users/auth"
)

func buildDogstatsMiddleware(dsd statsd.ClientInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == healthPath {
				return next(c)
			}

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			loggedTag := fmt.Sprintf("logged:%t", auth.HasAuth(c))
			duration := time.Since(start)
			_ = dsd.Distribution("letsblockit.request_duration", float64(duration.Nanoseconds()), []string{loggedTag}, 1)
			_ = dsd.Incr("letsblockit.request_count", []string{loggedTag, fmt.Sprintf("status:%d", c.Response().Status)}, 1)
			return nil
		}
	}
}

func collectBusinessStats(log echo.Logger, store db.Store, dsd statsd.ClientInterface) {
	collect := func() {
		stats, err := store.GetStats(context.Background())
		if err != nil {
			log.Error("cannot collect db stats: " + err.Error())
			return
		}
		_ = dsd.Gauge("letsblockit.total_list_count", float64(stats.ListsTotal), nil, 1)
		_ = dsd.Gauge("letsblockit.active_list_count", float64(stats.ListsActive), nil, 1)
		_ = dsd.Gauge("letsblockit.fresh_list_count", float64(stats.ListsFresh), nil, 1)

		instances, err := store.GetInstanceStats(context.Background())
		if err != nil {
			log.Error("cannot collect db stats: " + err.Error())
			return
		}
		for _, i := range instances {
			tags := []string{"filter_name:" + i.TemplateName}
			_ = dsd.Gauge("letsblockit.instance_count", float64(i.Total), tags, 1)
			_ = dsd.Gauge("letsblockit.fresh_instance_count", float64(i.Fresh), tags, 1)
		}
	}

	_ = dsd.Incr("letsblockit.startup", nil, 1)
	collect()
	for range time.Tick(5 * time.Minute) {
		collect()
	}
}

func collectMemStats(dsd statsd.ClientInterface) {
	oldS, newS := &runtime.MemStats{}, &runtime.MemStats{}
	gauge := func(name string, value uint64) {
		_ = dsd.Gauge(name, float64(value), nil, 1)
	}
	rate32 := func(name string, oldValue, newValue uint32) {
		if newValue > oldValue {
			_ = dsd.Count(name, int64(newValue-oldValue), nil, 1)
		} else {
			_ = dsd.Count(name, int64(math.MaxUint32-oldValue+newValue), nil, 1)
		}
	}
	rate64 := func(name string, oldValue, newValue uint64) {
		if newValue > oldValue {
			_ = dsd.Count(name, int64(newValue-oldValue), nil, 1)
		} else {
			_ = dsd.Count(name, int64(math.MaxUint64-oldValue+newValue), nil, 1)
		}
	}
	collect := func() {
		runtime.ReadMemStats(newS)
		gauge("go_expvar.memstats.heap_alloc", newS.HeapAlloc)
		gauge("go_expvar.memstats.heap_idle", newS.HeapIdle)
		gauge("go_expvar.memstats.heap_inuse", newS.HeapInuse)
		gauge("go_expvar.memstats.heap_objects", newS.HeapObjects)
		gauge("go_expvar.memstats.heap_released", newS.HeapReleased)
		gauge("go_expvar.memstats.heap_sys", newS.HeapSys)
		rate64("go_expvar.memstats.frees", oldS.Frees, newS.Frees)
		rate64("go_expvar.memstats.lookups", oldS.Lookups, newS.Lookups)
		rate64("go_expvar.memstats.mallocs", oldS.Mallocs, newS.Mallocs)
		rate64("go_expvar.memstats.total_alloc", oldS.TotalAlloc, newS.TotalAlloc)
		rate32("go_expvar.memstats.num_gc", oldS.NumGC, newS.NumGC)
		rate64("go_expvar.memstats.pause_total_ns", oldS.PauseTotalNs, newS.PauseTotalNs)
		oldS, newS = newS, oldS
	}

	collect()
	for range time.Tick(time.Minute) {
		collect()
	}
}

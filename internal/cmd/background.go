package cmd

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Background struct {
	command                    Command
	interval                   time.Duration
	log                        log.Logger
	executionsCounter          metrics.Counter
	executionDurationHistogram metrics.Histogram
	stopped                    int32
}

func NewBackground(cmd Command, opts ...BackgroundOption) *Background {
	c := &Background{
		command:                    cmd,
		interval:                   time.Millisecond * 200,
		executionsCounter:          discard.NewCounter(),
		executionDurationHistogram: discard.NewHistogram(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type BackgroundOption func(*Background)

// BackgroundWithInterval sets the time duration within command executions
func BackgroundWithInterval(interval time.Duration) BackgroundOption {
	return func(bg *Background) {
		bg.interval = interval
	}
}

func BackgroundWithLogger(l log.Logger) BackgroundOption {
	return func(bg *Background) {
		bg.log = l
	}
}

func BackgroundWithExecutionDurationHistogramCounter(h metrics.Histogram) BackgroundOption {
	return func(bg *Background) {
		bg.executionDurationHistogram = h
	}
}

func (bg *Background) Run() {
	level.Info(bg.log).Log("msg", "running background process")

	t := time.NewTicker(bg.interval)
	defer t.Stop()

	for !bg.isStopped() {
		<-t.C

		bg.executionsCounter.Add(1)
		timer := metrics.NewTimer(bg.executionDurationHistogram)
		timer.Unit(time.Millisecond)

		err := bg.command.Execute(context.Background())
		if err != nil {
			level.Error(bg.log).Log("msg", "command error", "err", err)
		}
		timer.ObserveDuration()
	}

	level.Info(bg.log).Log("msg", "stopping background process")
}

func (c *Background) isStopped() bool {
	return atomic.LoadInt32(&(c.stopped)) != 0
}

func (c *Background) Stop() {
	atomic.StoreInt32(&(c.stopped), int32(1))
}

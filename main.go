package main

import (
	"context"
	"fmt"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler/storage"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// example
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sch := scheduler.New(logger)

	j := scheduler.NewJob(
		"kekw-printer",
		func(ctx context.Context, arg ...any) {
			fmt.Println("kekw")
		},
		time.Second*10,
	)

	l := scheduler.NewJob(
		"lol-printer",
		func(ctx context.Context, arg ...any) {
			fmt.Println("lol")
		},
		time.Minute*20,
	)

	sch.AddJob(j)
	sch.AddJob(l)

	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
	defer cancel()

	conn := "./storage/jobs.sqlite"
	tableName := "jobs"

	jds, err := storage.NewSqliteJobDataStorage(conn, tableName)
	if err != nil {
		panic(err)
	}
	sch.SetJobDataStorage(jds)

	err = sch.Run(ctx)
	if err != nil {
		panic(err)
	}
}

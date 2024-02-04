package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler/storage"
	_ "github.com/mattn/go-sqlite3"
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

	db, err := sql.Open("sqlite3", "./storage/jobs.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tableName := "jobs"
	jds, err := storage.NewSqliteJobDataStorage(db, tableName)
	if err != nil {
		panic(err)
	}
	sch.SetJobDataStorage(jds)

	err = sch.Run(ctx)
	if err != nil {
		panic(err)
	}
}

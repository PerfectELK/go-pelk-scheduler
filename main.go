package main

import (
	"context"
	"fmt"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler/storage"
	"log/slog"
	"os"
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

	sch.AddJob(j)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*6))
	//ctx, cancel := context.WithCancel(context.Background())
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

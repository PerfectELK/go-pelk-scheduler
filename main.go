package main

import (
	"context"
	"fmt"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler"
	"log/slog"
	"os"
	"time"
)

// example
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sch := scheduler.New(logger)

	j := scheduler.NewJob("kekw", func(ctx context.Context, arg ...any) {
		fmt.Println("kekw")
	}, time.Second*10)

	sch.AddJob(j)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*35))
	defer cancel()

	sch.Run(ctx)
}

package main

import (
	"context"
	"fmt"
	"github.com/frisbm/cron"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	schedule, err := cron.Parse("* * * * *")
	if err != nil {
		log.Fatal(err)
	}

	task := NewTask(schedule, func(ctx context.Context) error {
		fmt.Println(fmt.Sprintf("The time is currently: %s", time.Now().UTC().Format(time.RFC3339)))
		return nil
	})

	scheduler := NewScheduler().AddTask(task)
	err = scheduler.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

type Task struct {
	cron *cron.Cron
	fn   func(ctx context.Context) error
}

type Scheduler struct {
	tasks []*Task
}

func NewTask(cron *cron.Cron, fn func(ctx context.Context) error) *Task {
	return &Task{
		cron: cron,
		fn:   fn,
	}
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) AddTask(task *Task) *Scheduler {
	s.tasks = append(s.tasks, task)
	return s
}

func (s *Scheduler) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	for _, task := range s.tasks {
		task := task
		g.Go(func() error {
			return runTaskSchedule(ctx, task)
		})
	}

	g.Go(func() error {
		return handleShutdownSignal(ctx, cancel)
	})

	<-ctx.Done()

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func runTaskSchedule(ctx context.Context, task *Task) error {
	ticker := time.NewTicker(minuteTick())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if task.cron.Now() {
				if err := task.fn(ctx); err != nil {
					return fmt.Errorf("task execution failed: %w", err)
				}
			}
			ticker.Reset(minuteTick())
		}
	}
}

func minuteTick() time.Duration {
	return time.Second * time.Duration(60-time.Now().Second())
}

func handleShutdownSignal(ctx context.Context, cancel context.CancelFunc) error {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signalCh:
		fmt.Printf("Received signal: %v\n", sig)
	case <-ctx.Done():
		return nil
	}

	cancel()
	return nil
}

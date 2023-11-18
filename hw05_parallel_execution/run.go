package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkers           = errors.New("no workers")
)

type Task func() error

func worker(tasks chan Task, errorCounter *int64, maxErrorsCount int64, stop chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-stop:
			return
		case task := <-tasks:
			if atomic.LoadInt64(errorCounter) < maxErrorsCount && task != nil {
				if err := task(); err != nil {
					atomic.AddInt64(errorCounter, 1)
				}
			}
			wg.Done()
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if n == 0 {
		return ErrNoWorkers
	}
	var (
		errorCounter   = int64(0)
		maxErrorsCount = int64(m)
		tasksCh        = make(chan Task, m)
		stop           = make(chan struct{})
		wg             = &sync.WaitGroup{}
	)
	for i := 0; i < n; i++ {
		go worker(tasksCh, &errorCounter, maxErrorsCount, stop, wg)
	}
	for _, task := range tasks {
		if atomic.LoadInt64(&errorCounter) >= maxErrorsCount {
			break
		}
		wg.Add(1)
		tasksCh <- task
	}
	wg.Wait()
	close(stop)
	if errorCounter >= maxErrorsCount {
		return ErrErrorsLimitExceeded
	}
	return nil
}

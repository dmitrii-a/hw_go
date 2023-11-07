package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrNoWorkers           = errors.New("no workers")
)

type Task func() error

type TaskGroup struct {
	tasksCh      chan Task
	errorCounter chan struct{}
	stop         chan struct{}
	wg           *sync.WaitGroup
	mutex        *sync.Mutex
}

func (g *TaskGroup) worker() {
	for {
		select {
		case <-g.stop:
			return
		case task := <-g.tasksCh:
			g.wg.Add(1)
			if err := task(); err != nil {
				g.addError()
			}
			g.wg.Done()
		}
	}
}

func (g *TaskGroup) waitWorkers() {
	g.wg.Wait()
}

func (g *TaskGroup) stopWorkers() {
	close(g.stop)
}

func (g *TaskGroup) addError() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	select {
	case <-g.stop:
	default:
		if len(g.errorCounter) == cap(g.errorCounter) {
			g.stopWorkers()
		} else {
			g.errorCounter <- struct{}{}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if n == 0 {
		return ErrNoWorkers
	}
	group := TaskGroup{
		tasksCh:      make(chan Task),
		errorCounter: make(chan struct{}, m),
		stop:         make(chan struct{}),
		wg:           &sync.WaitGroup{},
		mutex:        &sync.Mutex{},
	}
	defer close(group.errorCounter)
	for i := 0; i < n; i++ {
		go group.worker()
	}
	for _, task := range tasks {
		select {
		case <-group.stop:
			break
		case group.tasksCh <- task:
		}
	}
	group.waitWorkers()
	if len(group.errorCounter) >= m {
		return ErrErrorsLimitExceeded
	}
	group.stopWorkers()
	return nil
}

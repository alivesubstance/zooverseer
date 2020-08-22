package ui

import (
	"context"
)

var createChan = make(chan Task, 1)
var errorChan = make(chan Task, 1)
var completeChan = make(chan Task, 1)

type Task interface {
	process() /* <-chan interface{}*/
	complete()
	fail()
	cancel()
}

type baseTask struct {
	Context    context.Context
	CancelFunc context.CancelFunc
}

func addTask(task Task) {
	createChan <- task
}

func failTask(task Task) {
	createChan <- task
}

func init() {
	go func() {
		for {
			select {
			case task := <-createChan:
				task.process()
			case task := <-completeChan:
				task.complete()
			case task := <-errorChan:
				task.fail()
			}
		}
	}()
}

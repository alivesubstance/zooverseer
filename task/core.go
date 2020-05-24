package task

import (
	"github.com/alivesubstance/zooverseer/core/zk"
)

var createChan = make(chan Handler, 1)
var completeChan = make(chan Handler, 1)

type Handler interface {
	Process()
	Complete()
}

type BaseTask struct {
	Error     error
	OnStart   func()
	OnError   func(error)
	OnSuccess func(interface{})
}

type ImportTask struct {
	path  string
	nodes []*zk.Node
	Handler
	BaseTask
}

type SearchTask struct {
	value string
	Handler
}

func init() {
	go func() {
		for {
			select {
			case task := <-createChan:
				task.Process()
			case task := <-completeChan:
				task.Complete()
			}
		}
	}()
}

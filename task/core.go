package task

import (
	"github.com/alivesubstance/zooverseer/core/zk"
)

// todo failed to use channels and goroutine together with gtk code.
// unable to call gtk code from goroutine. always fails. one of the error:
// fatal error: unexpected signal during runtime execution
// [signal SIGSEGV: segmentation violation code=0x1 addr=0x6 pc=0x6]

//var createChan = make(chan Handler, 1)
//var completeChan = make(chan Handler, 1)

type Handler interface {
	Process() /* <-chan interface{}*/
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

//func init() {
//	go func() {
//		for {
//			select {
//			case task := <-createChan:
//				task.Process()
//			case task := <-completeChan:
//				task.Complete()
//			}
//		}
//	}()
//}

package task

// todo failed to use channels and goroutine together with gtk code.
// unable to call gtk code from goroutine. always fails. one of the error:
// fatal error: unexpected signal during runtime execution
// [signal SIGSEGV: segmentation violation code=0x1 addr=0x6 pc=0x6]

var CreateChan = make(chan Task, 1)
var ErrorChan = make(chan Task, 1)
var CompleteChan = make(chan Task, 1)

type Task interface {
	Process() /* <-chan interface{}*/
	Complete(task Task)
	Fail(task Task)
}

//type ImportTask struct {
//	path  string
//	nodes []*zk.Node
//	Task
//}
//
//type SearchTask struct {
//	value string
//	Task
//}

func init() {
	go func() {
		for {
			select {
			case task := <-CreateChan:
				task.Process()
			case task := <-CompleteChan:
				task.Complete(task)
			case task := <-ErrorChan:
				task.Fail(task)
			}
		}
	}()
}

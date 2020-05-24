package task

import (
	"github.com/alivesubstance/zooverseer/core/zk"
	log "github.com/sirupsen/logrus"
)

type ExportTask struct {
	ZkPath string
	Tree   *zk.Node
	Handler
	*BaseTask
}

func CreateExportTask(
	zkPath string,
	onStart func(),
	onError func(err error),
	onSuccess func(tree interface{}),
) {
	baseTask := &BaseTask{
		OnStart:   onStart,
		OnError:   onError,
		OnSuccess: onSuccess,
	}
	exportTask := &ExportTask{
		ZkPath:   zkPath,
		BaseTask: baseTask,
	}

	createChan <- exportTask
}

func (t *ExportTask) Process() {
	log.Infof("Start exporting %s", t.ZkPath)

	t.OnStart()

	node, err := zk.CachingRepo.Export(t.ZkPath)
	t.Tree = node
	t.Error = err

	completeChan <- t
}

func (t *ExportTask) Complete() {
	if t.Error != nil {
		t.OnError(t.Error)
		return
	}

	log.Infof("Finish exporting %s", t.ZkPath)
	t.OnSuccess(t.Tree)
}

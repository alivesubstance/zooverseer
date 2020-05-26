package ui

import (
	"encoding/json"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/task"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

type ExportTask struct {
	ZkPath string
	Tree   *zk.Node
	task.Handler
	*task.BaseTask
}

type Metadata struct {
	RootPath  string
	CreatedAt time.Time
	// todo think to add connection
	//Connection struct {
	//	Name string
	//	Host string
	//	Port int
	//}
}

type ExportNodeReport struct {
	Metadata Metadata
	Node     *zk.Node
}

func ExportSelectedNode() {
	treeSelection, _ := getNodesTreeView().GetSelection()
	zkPath, _ := getTreeSelectedZkPath(treeSelection)

	onStartTask := func() {
		enableSpinner(true)
	}
	onError := func(err error) {
		CreateErrorDialog(GetMainWindow(), "Export from "+zkPath+" failed: "+err.Error())
	}
	onSuccess := func(jsonFilePath interface{}) {
		log.Infof("Exported %v to %v", zkPath, jsonFilePath)
		enableSpinner(false)
		showNodeExportResultDialog(jsonFilePath.(string))
	}

	createExportTask(zkPath, onStartTask, onError, onSuccess)
}

func showNodeExportResultDialog(jsonFilePath string) {
	exportResultDlg.setResultFile(jsonFilePath)
	exportResultDlg.dlg.Run()
	exportResultDlg.dlg.Hide()
}

func createExportTask(
	zkPath string,
	onStart func(),
	onError func(err error),
	onSuccess func(tree interface{}),
) {
	baseTask := &task.BaseTask{
		OnStart:   onStart,
		OnError:   onError,
		OnSuccess: onSuccess,
	}
	exportTask := &ExportTask{
		ZkPath:   zkPath,
		BaseTask: baseTask,
	}

	exportTask.Process()
}

func (t *ExportTask) Process() {
	log.Infof("Start exporting %s", t.ZkPath)

	t.OnStart()

	node, err := zk.CachingRepo.Export(t.ZkPath)
	t.Tree = node

	// todo replace same pieces with chan error
	if err != nil {
		t.OnError(t.Error)
		return
	}

	if _, err := os.Stat(core.Config.ExportDir); os.IsNotExist(err) {
		log.Tracef("Creating directory %v", core.Config.ExportDir)
		err = os.Mkdir(core.Config.ExportDir, 0775)
		if err != nil {
			t.OnError(err)
			return
		}
	}

	metadata := Metadata{t.ZkPath, time.Now()}
	report := ExportNodeReport{metadata, node}

	data, err := json.Marshal(report)
	if err != nil {
		t.OnError(err)
		return
	}

	jsonFilePath := t.createFilePath(metadata)
	log.Tracef("Save export result to %v", jsonFilePath)

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		t.OnError(err)
		return
	}

	log.Tracef("File %v successfully created", jsonFilePath)
	_, err = jsonFile.Write(data[:])
	if err != nil {
		t.OnError(err)
		return
	}

	log.Tracef("Data has been written to %v", jsonFilePath)
	log.Infof("Finish exporting %s", t.ZkPath)
	t.OnSuccess(jsonFilePath)
}

func (t *ExportTask) createFilePath(meta Metadata) string {
	return path.Join(
		core.Config.ExportDir,
		"export_"+meta.CreatedAt.Format("20060102_150405")+".json",
	)
}

package ui

import (
	"encoding/json"
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/task"
	"github.com/kennygrant/sanitize"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"time"
)

type ExportTask struct {
	ZkPath string
	Tree   *zk.Node
	task.Handler
	*task.BaseTask
}

type Meta struct {
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
	Meta Meta
	Data *zk.Node
}

func exportSelectedNode() {
	treeSelection, _ := getNodesTreeView().GetSelection()
	zkPath, _ := getTreeSelectedZkPath(treeSelection)
	createExportTask(zkPath)
}

func createExportTask(zkPath string) {
	onError := func(err error) {
		log.WithError(err).Infof("Failed to export from %v", zkPath)
		nodeExportDlg.showError(zkPath, err)
	}
	onSuccess := func(jsonFilePath interface{}) {
		log.Infof("Exported %v to %v", zkPath, jsonFilePath)
		nodeExportDlg.showResult(zkPath, jsonFilePath.(string))
	}

	baseTask := &task.BaseTask{
		OnError:   onError,
		OnSuccess: onSuccess,
	}
	exportTask := &ExportTask{
		ZkPath:   zkPath,
		BaseTask: baseTask,
	}

	task.CreateChan <- exportTask
	nodeExportDlg.startExport(zkPath)
}

func (t *ExportTask) Process() {
	log.Infof("Start exporting %s", t.ZkPath)

	node, err := zk.CachingRepo.Export(t.ZkPath)
	t.Tree = node

	// todo replace same pieces with chan error
	if err != nil {
		t.OnError(err)
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

	metadata := Meta{t.ZkPath, time.Now()}
	report := ExportNodeReport{metadata, node}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.OnError(err)
		return
	}

	absJsonFilePath, err := filepath.Abs(t.createFilePath(node, metadata))
	if err != nil {
		t.OnError(err)
		return
	}
	log.Tracef("Save export result to %v", absJsonFilePath)

	jsonFile, err := os.Create(absJsonFilePath)
	if err != nil {
		t.OnError(err)
		return
	}

	log.Tracef("File %v created", absJsonFilePath)
	_, err = jsonFile.Write(data[:])
	if err != nil {
		t.OnError(err)
		return
	}

	log.Tracef("Data has been written to %v", absJsonFilePath)
	log.Infof("Finish exporting %s", t.ZkPath)
	t.OnSuccess(absJsonFilePath)
}

func (t *ExportTask) createFilePath(node *zk.Node, meta Meta) string {
	return path.Join(
		core.Config.ExportDir,
		fmt.Sprintf("%s_%v.json", meta.CreatedAt.Format("20060102_150405"), sanitize.Name(node.Name)),
	)
}

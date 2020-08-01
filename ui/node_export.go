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
	ZkPath       string
	Tree         *zk.Node
	JsonFilePath string
	Error        error
	task.Task
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
	exportTask := &ExportTask{
		ZkPath: zkPath,
	}
	task.CreateChan <- exportTask

	nodeExportDlg.startExport(zkPath)
}

func (t *ExportTask) Fail(task task.Task) {
	exportTask := task.(*ExportTask)
	log.WithError(exportTask.Error).Errorf("Failed to export from %v", exportTask.ZkPath)
	nodeExportDlg.showError(exportTask.ZkPath, exportTask.Error)
}

func (t *ExportTask) Complete(task task.Task) {
	exportTask := task.(*ExportTask)
	jsonFilePath := exportTask.JsonFilePath
	log.Infof("Exported %v to %v", exportTask.ZkPath, jsonFilePath)
	nodeExportDlg.showResult(exportTask.ZkPath, jsonFilePath)
}

func (t *ExportTask) Process() {
	log.Infof("Start exporting %s", t.ZkPath)

	node, err := zk.CachingRepo.Export(t.ZkPath)
	if t.hasError(err) {
		return
	}
	t.Tree = node

	if _, err := os.Stat(core.Config.ExportDir); os.IsNotExist(err) {
		log.Tracef("Creating directory %v", core.Config.ExportDir)
		err = os.Mkdir(core.Config.ExportDir, 0775)
		if t.hasError(err) {
			return
		}
	}

	metadata := Meta{t.ZkPath, time.Now()}
	report := ExportNodeReport{metadata, node}

	data, err := json.MarshalIndent(report, "", "  ")
	if t.hasError(err) {
		return
	}

	absJsonFilePath, err := filepath.Abs(t.createFilePath(node, metadata))
	if t.hasError(err) {
		return
	}
	log.Tracef("Save export result to %v", absJsonFilePath)

	jsonFile, err := os.Create(absJsonFilePath)
	if t.hasError(err) {
		return
	}

	log.Tracef("File %v created", absJsonFilePath)
	_, err = jsonFile.Write(data[:])
	if t.hasError(err) {
		return
	}

	log.Infof("Finish exporting %s to %s", t.ZkPath, absJsonFilePath)
	t.JsonFilePath = absJsonFilePath
	task.CompleteChan <- t
}

func (t *ExportTask) hasError(err error) bool {
	if err != nil {
		t.Error = err
		task.ErrorChan <- t
		return true
	}
	return false
}

func (t *ExportTask) createFilePath(node *zk.Node, meta Meta) string {
	return path.Join(
		core.Config.ExportDir,
		fmt.Sprintf("%s_%v.json", meta.CreatedAt.Format("20060102_150405"), sanitize.Name(node.Name)),
	)
}

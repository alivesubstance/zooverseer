package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

type ExportResultDlg struct {
	dlg         *gtk.Dialog
	statusLabel *gtk.Label
	spinner     *gtk.Spinner
	resultLabel *gtk.Label
	okBtn       *gtk.Button
	mainWindow  *gtk.Window
}

func NewNodeExportDlg(mainWindow *gtk.Window) *ExportResultDlg {
	dlg := &ExportResultDlg{}
	dlg.statusLabel = GetObject("nodeExportStatusLabel").(*gtk.Label)
	dlg.spinner = GetObject("nodeExportSpinner").(*gtk.Spinner)
	dlg.resultLabel = GetObject("nodeExportResultLabel").(*gtk.Label)
	dlg.okBtn = GetObject("exportDlgOkBtn").(*gtk.Button)
	dlg.mainWindow = mainWindow
	dlg.dlg = GetObject("nodeExportDlg").(*gtk.Dialog)
	dlg.dlg.SetTransientFor(mainWindow)

	dlg.okBtn.Connect("clicked", dlg.onOkBtnClicked(dlg))
	return dlg
}

func (d *ExportResultDlg) startExport(path string) {
	d.resetDlg()

	d.setStatus("Start exporting from <b>" + path + "</b>")

	d.spinner.Start()
	d.mainWindow.SetSensitive(false)
	d.okBtn.SetSensitive(false)

	d.dlg.Run()
}

func (d *ExportResultDlg) showError(path string, err error) {
	d.resetDlg()

	d.setStatus("Export failed from <b>" + path + "</b>")
	d.resultLabel.SetMarkup("Error: " + err.Error())

	d.spinner.Stop()
	d.mainWindow.SetSensitive(true)
	d.okBtn.SetSensitive(true)
}

func (d *ExportResultDlg) showResult(path string, filePath string) {
	d.resetDlg()

	html := "Data exported to <a href=\"file://" + filePath + "\">file</a>"
	nodeExportDlg.resultLabel.SetMarkup(html)

	d.setStatus("Export finished from <b>" + path + "</b>")
	d.spinner.Stop()
	d.mainWindow.SetSensitive(true)
	d.okBtn.SetSensitive(true)
}

func (d *ExportResultDlg) setStatus(html string) string {
	d.statusLabel.SetMarkup(html)
	return html
}

func (d *ExportResultDlg) resetDlg() {
	d.statusLabel.SetMarkup("")
	d.resultLabel.SetMarkup("")
	d.spinner.Stop()
}

func (d *ExportResultDlg) onOkBtnClicked(dlg *ExportResultDlg) func() {
	return func() { dlg.dlg.Hide() }
}

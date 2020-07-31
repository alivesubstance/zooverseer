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
	exportResultDlg := &ExportResultDlg{}
	exportResultDlg.statusLabel = GetObject("nodeExportStatusLabel").(*gtk.Label)
	exportResultDlg.spinner = GetObject("nodeExportSpinner").(*gtk.Spinner)
	exportResultDlg.resultLabel = GetObject("nodeExportResultLabel").(*gtk.Label)
	exportResultDlg.okBtn = GetObject("exportDlgOkBtn").(*gtk.Button)
	exportResultDlg.mainWindow = mainWindow
	exportResultDlg.dlg = GetObject("nodeExportDlg").(*gtk.Dialog)
	exportResultDlg.dlg.SetTransientFor(mainWindow)

	exportResultDlg.okBtn.Connect("clicked", exportResultDlg.onOkBtnClicked(exportResultDlg))
	return exportResultDlg
}

func (d *ExportResultDlg) startExport(path string) {
	d.resetDlg()

	d.setStatus("Start exporting from <b>" + path + "</b>")
	d.setSensitive(false)

	d.dlg.Run()
}

func (d *ExportResultDlg) showError(path string, err error) {
	d.resetDlg()

	d.setStatus("Export failed from <b>" + path + "</b>")
	d.resultLabel.SetMarkup("Error: " + err.Error())

	d.setSensitive(true)
}

func (d *ExportResultDlg) showResult(path string, filePath string) {
	d.resetDlg()

	html := "Data exported to <a href=\"file://" + filePath + "\">file</a>"
	nodeExportDlg.resultLabel.SetMarkup(html)

	d.setStatus("Export finished from <b>" + path + "</b>")
	d.setSensitive(true)
}

func (d *ExportResultDlg) setSensitive(value bool) {
	if !value {
		d.spinner.Start()
	} else {
		d.spinner.Stop()
	}
	d.mainWindow.SetSensitive(value)
	d.okBtn.SetSensitive(value)
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

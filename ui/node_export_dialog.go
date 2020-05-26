package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

type ExportResultDlg struct {
	dlg   *gtk.Dialog
	label *gtk.Label
}

func NewExportResultDlg(mainWindow *gtk.Window) *ExportResultDlg {
	exportResultDlg := &ExportResultDlg{}
	exportResultDlg.label = GetObject("exportLabel").(*gtk.Label)
	exportResultDlg.dlg = GetObject("exportDlg").(*gtk.Dialog)
	exportResultDlg.dlg.SetTransientFor(mainWindow)

	GetObject("exportDlgOkBtn").(*gtk.Button).Connect("clicked", func() {
		exportResultDlg.dlg.Hide()
	})
	return exportResultDlg
}

func (d *ExportResultDlg) setResultFile(filePath string) {
	s := "Node exported to <a href=\"file://" + filePath + "\">file</a>"
	exportResultDlg.label.SetMarkup(s)
}

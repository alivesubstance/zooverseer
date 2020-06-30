package ui

import "github.com/gotk3/gotk3/gtk"

func onMenuAboutBtnCloseClicked() {
	GetObject("aboutDlg").(*gtk.Dialog).Hide()
}

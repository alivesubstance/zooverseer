package ui

import "github.com/gotk3/gotk3/gtk"

func onMenuAboutActivate() {
	GetObject("aboutDlg").(*gtk.Dialog).Show()
}

func onMenuAboutBtnCloseClicked() {
	GetObject("aboutDlg").(*gtk.Dialog).Hide()
}

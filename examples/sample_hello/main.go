package main

import (
	"github.com/samuel-jimenez/windigo"
)

func main() {
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetSize(400, 300)
	mainWindow.SetText("Hello World Demo")

	// Main window menu. Context menus on controls also available.
	menu := mainWindow.NewMenu()
	fileMn := menu.AddSubMenu("File")
	fileMn.AddItem("New", windigo.Shortcut{windigo.ModControl, windigo.KeyN})
	editMn := menu.AddSubMenu("Edit")
	cutMn := editMn.AddItem("Cut", windigo.Shortcut{windigo.ModControl, windigo.KeyX})
	copyMn := editMn.AddItem("Copy", windigo.NoShortcut)
	pasteMn := editMn.AddItem("Paste", windigo.NoShortcut)
	menu.Show()
	// Menu items can be disabled and checked.
	copyMn.SetCheckable(true)
	copyMn.SetChecked(true)
	pasteMn.SetEnabled(false)

	cutMn.OnClick().Bind(func(e *windigo.Event) {
		windigo.MsgBoxOk(mainWindow, "Cut", "Click event")
	})

	edt := windigo.NewEdit(mainWindow)
	edt.SetPos(10, 20)
	// Most Controls have default size unless SetSize is called.
	edt.SetText("edit text")

	btn := windigo.NewPushButton(mainWindow)
	btn.SetText("Show or Hide")
	btn.SetPos(40, 50)
	btn.SetSize(100, 40)
	btn.OnClick().Bind(func(e *windigo.Event) {
		if edt.Visible() {
			edt.Hide()
		} else {
			edt.Show()
		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	windigo.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *windigo.Event) {
	windigo.Exit()
}

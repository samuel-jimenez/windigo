package main

import (
	"fmt"

	"github.com/samuel-jimenez/windigo"
)

func btnOnClick(arg *windigo.Event) {
	fmt.Println("Button clicked")
}

func wndOnClose(arg *windigo.Event) {
	windigo.Exit()
}

func main() {
	mainWindow := windigo.NewForm(nil)
	dock := windigo.NewSimpleDock(mainWindow)
	//mainWindow.SetLayout(dock)

	mainWindow.SetSize(700, 600)
	mainWindow.SetText("Controls Demo")

	menu := mainWindow.NewMenu()
	fileMn := menu.AddSubMenu("File")
	fileMn.AddItem("New", windigo.NoShortcut)
	editMn := menu.AddSubMenu("Edit")
	cutMn := editMn.AddItem("Cut", windigo.Shortcut{windigo.ModControl, windigo.KeyX})
	copyMn := editMn.AddItem("Copy", windigo.NoShortcut)
	pasteMn := editMn.AddItem("Paste", windigo.NoShortcut)
	menu.Show()
	copyMn.SetCheckable(true)
	copyMn.SetChecked(true)
	pasteMn.SetEnabled(false)

	cutMn.OnClick().Bind(func(e *windigo.Event) {
		println("cut click")
	})

	imlistTb := windigo.NewImageList(16, 16)
	imlistTb.AddResIcon(10)
	imlistTb.AddResIcon(12)
	imlistTb.AddResIcon(15)

	toolbar := windigo.NewToolbar(mainWindow)
	toolbar.SetImageList(imlistTb)
	addBtn := toolbar.AddButton("Add", 1)
	toolbar.AddSeparator()
	runBtn := toolbar.AddButton("Run Now Fast", 2)
	toolbar.Show()

	runBtn.OnClick().Bind(func(e *windigo.Event) {
		println("runBtn click")
	})

	dock.Dock(toolbar, windigo.Top) // toolbars always dock to the top
	//dock.Dock(tree, windigo.Fill)

	slide := windigo.NewSlider(mainWindow)
	slide.SetPos(10, 50)
	slide.OnScroll().Bind(func(e *windigo.Event) {
		println(slide.Value())
	})

	addBtn.OnClick().Bind(func(e *windigo.Event) {
		println("addBtn click")
		slide.SetValue(30)
	})

	//track.SetRange(0, 100)
	//track.SetValue(20)

	mainWindow.Center()
	mainWindow.Show()
	dock.Update()
	mainWindow.OnClose().Bind(wndOnClose)

	mainWindow.RunMainLoop()
}

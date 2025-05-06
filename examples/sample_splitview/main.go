package main

import (
	"fmt"

	"github.com/samuel-jimenez/windigo"
)

func btnOnClick(arg *windigo.Event) {
	//edt.SetCaption("Got you !!!")
	fmt.Println("Button clicked")
}

func wndOnClose(arg *windigo.Event) {
	windigo.Exit()
}

func main() {
	mainWindow := windigo.NewForm(nil)

	mainWindow.SetSize(700, 600)
	mainWindow.SetText("Controls Demo")

	//none := windigo.Shortcut{}

	menu := mainWindow.NewMenu()
	fileMn := menu.AddSubMenu("File")
	fileMn.AddItem("New", windigo.NoShortcut)
	editMn := menu.AddSubMenu("Edit")
	delMn := editMn.AddItem("Delete", windigo.Shortcut{windigo.ModControl, windigo.KeyX})
	delAllMn := editMn.AddItem("Delete All", windigo.NoShortcut)
	menu.Show()

	delMn.OnClick().Bind(func(e *windigo.Event) {
		dlg := windigo.NewDialog(mainWindow)
		dlg.Center()
		dlg.Show()
	})

	delAllMn.OnClick().Bind(func(e *windigo.Event) {
		dlg := windigo.NewDialog(mainWindow)
		dlg.Center()
		dlg.Show()
	})

	toolbar := windigo.NewPanel(mainWindow)
	toolbar.SetPos(0, 0)
	toolbar.SetSize(100, 40)

	btnRun := windigo.NewIconButton(toolbar)
	btnRun.SetText(" Run")
	btnRun.SetPos(2, 2)
	btnRun.SetSize(98, 38)
	btnRun.SetResIcon(15)

	btnRun.OnClick().Bind(func(e *windigo.Event) {
		println("event OnClick")
	})

	//tipRun := windigo.NewToolTip(mainWindow)
	//tipRun.AddTool(btnRun, "Run project")

	btnEdit := windigo.NewPushButton(toolbar)
	btnEdit.SetText(" Edit")
	btnEdit.SetPos(102, 2)
	btnEdit.SetSize(98, 38)
	btnEdit.SetResIcon(18)

	left := windigo.NewMultiEdit(mainWindow)
	left.SetPos(5, 5)
	right := windigo.NewMultiEdit(mainWindow)
	right.SetPos(50, 100)

	split := windigo.NewVResizer(mainWindow)
	split.SetControl(left, right, windigo.Left, 150)

	// --- Add controls to Dock, LoadStateFile and Show window in this order
	mainWindow.Center()
	mainWindow.Show()

	dock := windigo.NewSimpleDock(mainWindow)
	//mainWindow.SetLayout(dock)
	dock.Dock(toolbar, windigo.Top)
	dock.Dock(left, windigo.Left)
	dock.Dock(split, windigo.Left)
	dock.Dock(right, windigo.Fill)
	// if err := dock.LoadStateFile("layout.json"); err != nil {
	// 	log.Println(err)
	// }

	mainWindow.OnClose().Bind(func(e *windigo.Event) {
		dock.SaveStateFile("layout.json") // error gets ignored
		windigo.Exit()
	})

	dock.Update()

	windigo.RunMainLoop()
	// --- end of Dock and main window management

}

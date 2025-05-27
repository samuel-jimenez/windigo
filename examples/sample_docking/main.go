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

type Item struct {
	T string
}

func (item Item) Text() string    { return item.T }
func (item Item) ImageIndex() int { return 0 }

func main() {
	//windigo.Init()

	mainWindow := windigo.NewForm(nil)
	dock := windigo.NewSimpleDock(mainWindow)

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

	tabs := windigo.NewMultiPanel(mainWindow)
	tabs.SetPos(10, 10)
	tabs.SetSize(100, 92)

	tree := windigo.NewTreeView(mainWindow)
	tree.SetPos(10, 80)
	p := &Item{"First Item"}
	tree.InsertItem(p, nil, nil)
	sec := &Item{"Second"}
	if err := tree.InsertItem(sec, p, nil); err != nil {
		panic(err)
	}
	if err := tree.InsertItem(&Item{"Third"}, p, nil); err != nil {
		panic(err)
	}
	if err := tree.InsertItem(&Item{"Fourth"}, p, nil); err != nil {
		panic(err)
	}
	for i := 0; i < 50; i++ {
		if err := tree.InsertItem(&Item{"after second"}, sec, nil); err != nil {
			panic(err)
		}
	}
	tree.Expand(p)
	tree.OnCollapse().Bind(func(e *windigo.Event) {
		println("collapse")
	})

	cutMn.OnClick().Bind(func(e *windigo.Event) {
		println("cut click")
		ok := tree.EnsureVisible(p)
		fmt.Println("result of EnsureVisible", ok)
	})

	panel := windigo.NewPanel(tabs)
	tabs.AddPanel(panel)

	panelDock := windigo.NewSimpleDock(panel)
	panel.SetLayout(panelDock)
	panel.SetPos(0, 0)

	panelErr := windigo.NewErrorPanel(panel)
	panelErr.SetPos(140, 10)
	panelErr.SetSize(200, 32)
	panelErr.ShowAsError(false)

	edt := windigo.NewEdit(panel)
	edt.SetPos(10, 535)
	edt.SetText("some text")

	btn := windigo.NewPushButton(panel)
	btn.SetText("Button")
	btn.SetSize(100, 40)
	btn.OnClick().Bind(func(e *windigo.Event) {
		if edt.Visible() {
			edt.Hide()
		} else {
			edt.Show()
		}
	})
	btn.SetResIcon(13)

	panelDock.Dock(btn, windigo.Top)
	panelDock.Dock(edt, windigo.Top)
	panelDock.Dock(panelErr, windigo.Top)

	dock.Dock(tree, windigo.Left)
	dock.Dock(tabs, windigo.Top)

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	mainWindow.RunMainLoop()
}

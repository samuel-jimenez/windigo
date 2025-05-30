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

type Item struct {
	T string
}

func (item Item) Text() string    { return item.T }
func (item Item) ImageIndex() int { return 1 }

func main() {
	mainWindow := windigo.NewForm(nil)
	dock := windigo.NewSimpleDock(mainWindow)
	mainWindow.SetLayout(dock)

	mainWindow.SetSize(540, 540)
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

	imlist := windigo.NewImageList(16, 16)
	imlist.AddResIcon(10)
	imlist.AddResIcon(12)
	imlist.AddResIcon(15)

	scroll := windigo.NewScrollView(mainWindow)
	tree := windigo.NewTreeView(scroll)
	scroll.SetChild(&tree.ControlBase)
	//scroll.Show()

	tree.SetImageList(imlist)
	tree.SetPos(10, 80)
	tree.SetSize(800, 800)
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

	imlistTb := windigo.NewImageList(16, 16)
	imlistTb.AddResIcon(10)
	imlistTb.AddResIcon(12)
	imlistTb.AddResIcon(15)

	toolbar := windigo.NewToolbar(mainWindow)
	toolbar.SetImageList(imlistTb)
	addBtn := toolbar.AddButton("Add", 1)
	toolbar.AddSeparator()
	runBtn := toolbar.AddButton("Run", 2)
	toolbar.Show()

	runBtn.OnClick().Bind(func(e *windigo.Event) {
		println("runBtn click")
	})

	addBtn.OnClick().Bind(func(e *windigo.Event) {
		println("addBtn click")
	})

	dock.Dock(toolbar, windigo.Top) // toolbars always dock to the top
	dock.Dock(scroll, windigo.Fill)

	mainWindow.Center()
	dock.Update()
	mainWindow.Show()

	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop()
}

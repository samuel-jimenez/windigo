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
	T       []string
	checked bool
}

func (item Item) Text() []string    { return item.T }
func (item *Item) SetText(s string) { item.T[0] = s }

func (item Item) Checked() bool            { return item.checked }
func (item *Item) SetChecked(checked bool) { item.checked = checked }
func (item Item) ImageIndex() int          { return 0 }

func main() {
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

	cutMn.OnClick().Bind(func(e *windigo.Event) {
		println("cut click")
	})

	imlistTb := windigo.NewImageList(16, 16)
	imlistTb.AddResIcon(10)
	imlistTb.AddResIcon(12)
	imlistTb.AddResIcon(15)

	// --- Toolbar
	toolbar := windigo.NewToolbar(mainWindow)
	toolbar.SetImageList(imlistTb)
	addBtn := toolbar.AddButton("Add", 1)
	toolbar.AddSeparator()
	runBtn := toolbar.AddButton("Run Now Fast", 2)
	toolbar.Show()

	runBtn.OnClick().Bind(func(e *windigo.Event) {
		println("runBtn click")
	})

	addBtn.OnClick().Bind(func(e *windigo.Event) {
		println("addBtn click")
	})

	// --- Tabs
	tabs := windigo.NewTabView(mainWindow)
	panel1 := tabs.AddPanel("First")
	panel2 := tabs.AddPanel("Second")
	panel3 := tabs.AddPanel("Third")
	panel4 := tabs.AddPanel("Fourth")

	edt := windigo.NewEdit(panel1)
	edt.SetPos(100, 10)

	edt2 := windigo.NewEdit(panel2)
	edt2.SetPos(40, 10)

	btn := windigo.NewPushButton(panel3)
	btn.OnClick().Bind(func(e *windigo.Event) {
		println("click")
	})

	imlist := windigo.NewImageList(16, 16)
	imlist.AddResIcon(12)

	ls := windigo.NewListView(panel4)
	ls.SetImageList(imlist)
	ls.EnableEditLabels(false)
	ls.SetCheckBoxes(true)
	//ls.EnableFullRowSelect(true)
	//ls.EnableHotTrack(true)
	//ls.EnableSortHeader(true)
	//ls.EnableSortAscending(true)

	ls.AddColumn("One", 120)
	ls.AddColumn("Two", 120)
	ls.SetPos(10, 180)
	p1 := &Item{[]string{"First Item", "A"}, true}
	ls.AddItem(p1)
	p2 := &Item{[]string{"Second Item", "B"}, true}
	ls.AddItem(p2)
	p3 := &Item{[]string{"Third Item", "C"}, true}
	ls.AddItem(p3)
	for i := 0; i < 200; i++ {
		p4 := &Item{[]string{"Fourth Item", "D"}, false}
		ls.AddItem(p4)
	}

	// --- Dock
	dock2 := windigo.NewSimpleDock(panel4)
	dock2.Dock(ls, windigo.Fill)
	tabs.SetCurrent(0)

	dock.Dock(toolbar, windigo.Top)        // toolbars always dock to the top
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	windigo.RunMainLoop()
}

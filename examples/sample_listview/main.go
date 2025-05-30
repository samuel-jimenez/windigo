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

	none := windigo.Shortcut{}

	imlist := windigo.NewImageList(16, 16)
	imlist.AddResIcon(16)
	imlist.AddResIcon(10)
	imlist.AddResIcon(13)

	ls := windigo.NewListView(mainWindow)
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

	menu := mainWindow.NewMenu()
	fileMn := menu.AddSubMenu("File")
	fileMn.AddItem("New", none)
	editMn := menu.AddSubMenu("Edit")
	delMn := editMn.AddItem("Delete", windigo.Shortcut{windigo.ModControl, windigo.KeyX})
	delAllMn := editMn.AddItem("Delete All", none)
	menu.Show()

	ls.OnEndLabelEdit().Bind(func(e *windigo.Event) {
		println("edited", e)
		// acccept label edit event!
		//d := e.Data.(*windigo.LabelEditEventData)
		//d.Item.SetText(d.Text)
		//fmt.Println(d.Item.Text())
	})

	delMn.OnClick().Bind(func(e *windigo.Event) {
		items := ls.SelectedItems()
		for _, it := range items {
			fmt.Println(it)
		}
	})

	delAllMn.OnClick().Bind(func(e *windigo.Event) {
		ls.DeleteAllItems()
	})

	ls.OnClick().Bind(func(e *windigo.Event) {
		println("onClick listview")
	})

	dock.Dock(ls, windigo.Fill)

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	mainWindow.RunMainLoop()
}

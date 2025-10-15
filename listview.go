/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

// ListItem represents an item in a ListView widget.
type ListItem interface {
	Text() []string  // Text returns the text of the multi-column item.
	ImageIndex() int // ImageIndex is used only if SetImageList is called on the listview
}

// ListItemChecker is used for checkbox support in ListView.
type ListItemChecker interface {
	Checked() bool
	SetChecked(checked bool)
}

// ListItemSetter is used in OnEndLabelEdit event.
type ListItemSetter interface {
	SetText(s string) // set first item in the array via LabelEdit event
}

// StringListItem is helper for basic string lists.
type StringListItem struct {
	ID    int
	Data  string
	Check bool
}

func (s StringListItem) Text() []string          { return []string{s.Data} }
func (s StringListItem) Checked() bool           { return s.Check }
func (s StringListItem) SetChecked(checked bool) { s.Check = checked }
func (s StringListItem) ImageIndex() int         { return 0 }

type ListView struct {
	ControlBase

	iml *ImageList
	lastIndex,
	sortCol,
	cols int // count of columns

	sortAscending bool

	item2Handle map[ListItem]uintptr
	handle2Item map[uintptr]ListItem

	sort func(int, bool)

	onEndLabelEdit,
	onDoubleClick,
	onClick,
	onRClick,
	onKeyDown,
	onItemChanging,
	onItemChanged,
	onCheckChanged,
	onViewChange,
	onEndScroll EventManager
}

func NewListView(parent Controller) *ListView {
	control := new(ListView)

	control.InitControl(w32.WC_LISTVIEW, parent /*w32.WS_EX_CLIENTEDGE*/, 0,
		w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.LVS_REPORT|w32.LVS_EDITLABELS|w32.LVS_SHOWSELALWAYS)

	control.item2Handle = make(map[ListItem]uintptr)
	control.handle2Item = make(map[uintptr]ListItem)

	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetSize(200, 400)

	control.sortCol = -1
	control.sortAscending = true

	if err := control.SetTheme("Explorer"); err != nil {
		// theme error is ignored
	}
	return control
}

// FIXME: Changes the state of an item in a list-view control. Refer LVM_SETITEMSTATE message.
func (control *ListView) setItemState(i int, state, mask uint) {
	var item w32.LVITEM
	item.State, item.StateMask = uint32(state), uint32(mask)
	w32.SendMessage(control.hwnd, w32.LVM_SETITEMSTATE, uintptr(i), uintptr(unsafe.Pointer(&item)))
}

func (control *ListView) EnableSingleSelect(enable bool) {
	ToggleStyle(control.hwnd, enable, w32.LVS_SINGLESEL)
}

func (control *ListView) EnableSortHeader(enable bool,
	sort func(int, bool)) {
	if enable {
		control.sort = sort
	} else {
		control.sort = nil
	}

	ToggleStyle(control.hwnd, !enable, w32.LVS_NOSORTHEADER)
}

func (control *ListView) EnableSortAscending(enable bool) {
	ToggleStyle(control.hwnd, enable, w32.LVS_SORTASCENDING)
}

func (control *ListView) EnableEditLabels(enable bool) {
	ToggleStyle(control.hwnd, enable, w32.LVS_EDITLABELS)
}

func (control *ListView) SetExtendedStyle(style uintptr, enable bool) {
	if enable {
		w32.SendMessage(control.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, style, style)
	} else {
		w32.SendMessage(control.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, style, 0)
	}
}

func (control *ListView) EnableFullRowSelect(enable bool) {
	control.SetExtendedStyle(w32.LVS_EX_FULLROWSELECT, enable)
}

func (control *ListView) EnableDoubleBuffer(enable bool) {
	control.SetExtendedStyle(w32.LVS_EX_DOUBLEBUFFER, enable)
}

func (control *ListView) EnableHotTrack(enable bool) {
	control.SetExtendedStyle(w32.LVS_EX_TRACKSELECT, enable)
}

func (control *ListView) EnableGridlines(enable bool) {
	control.SetExtendedStyle(w32.LVS_EX_GRIDLINES, enable)
}

func (control *ListView) SetItemCount(count int) bool {
	return w32.SendMessage(control.hwnd, w32.LVM_SETITEMCOUNT, uintptr(count), 0) != 0
}

func (control *ListView) ItemCount() int {
	return int(w32.SendMessage(control.hwnd, w32.LVM_GETITEMCOUNT, 0, 0))
}

func (control *ListView) ItemAt(x, y int) ListItem {
	hti := w32.LVHITTESTINFO{Pt: w32.POINT{X: int32(x), Y: int32(y)}}
	w32.SendMessage(control.hwnd, w32.LVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))
	return control.findItemByIndex(int(hti.IItem))
}

func (control *ListView) Items() (list []ListItem) {
	for item := range control.item2Handle {
		list = append(list, item)
	}
	return list
}

func (control *ListView) AddColumn(caption string, width int) {
	var lc w32.LVCOLUMN
	lc.Mask = w32.LVCF_TEXT
	if width != 0 {
		lc.Mask = lc.Mask | w32.LVCF_WIDTH
		lc.Cx = int32(width)
	}
	lc.PszText = syscall.StringToUTF16Ptr(caption)
	control.insertLvColumn(&lc, control.cols)
	control.cols++
}

func (control *ListView) SetColumnWidth(col, width int) {
	if w32.SendMessage(control.hwnd, w32.LVM_SETCOLUMNWIDTH, uintptr(col), uintptr(width)) == 0 {
		//panic("LVM_SETCOLUMNWIDTH failed")
	}
}

func (control *ListView) GetNumColumns() int {
	return control.cols
}

// StretchLastColumn makes the last column take up all remaining horizontal
// space of the *ListView.
// The effect of this is not persistent.
func (control *ListView) StretchLastColumn() error {
	if control.cols == 0 {
		return nil
	}
	if w32.SendMessage(control.hwnd, w32.LVM_SETCOLUMNWIDTH, uintptr(control.cols-1), w32.LVSCW_AUTOSIZE_USEHEADER) == 0 {
		//panic("LVM_SETCOLUMNWIDTH failed")
	}
	return nil
}

// CheckBoxes returns if the *TableView has check boxes.
func (control *ListView) CheckBoxes() bool {
	return w32.SendMessage(control.hwnd, w32.LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)&w32.LVS_EX_CHECKBOXES > 0
}

// SetCheckBoxes sets if the *TableView has check boxes.
func (control *ListView) SetCheckBoxes(value bool) {
	exStyle := w32.SendMessage(control.hwnd, w32.LVM_GETEXTENDEDLISTVIEWSTYLE, 0, 0)
	oldStyle := exStyle
	if value {
		exStyle |= w32.LVS_EX_CHECKBOXES
	} else {
		exStyle &^= w32.LVS_EX_CHECKBOXES
	}
	if exStyle != oldStyle {
		w32.SendMessage(control.hwnd, w32.LVM_SETEXTENDEDLISTVIEWSTYLE, 0, exStyle)
	}

	mask := w32.SendMessage(control.hwnd, w32.LVM_GETCALLBACKMASK, 0, 0)
	if value {
		mask |= w32.LVIS_STATEIMAGEMASK
	} else {
		mask &^= w32.LVIS_STATEIMAGEMASK
	}

	if w32.SendMessage(control.hwnd, w32.LVM_SETCALLBACKMASK, mask, 0) == w32.FALSE {
		panic("SendMessage(LVM_SETCALLBACKMASK)")
	}
}

func (control *ListView) applyImage(lc *w32.LVITEM, imIndex int) {
	if control.iml != nil {
		lc.Mask |= w32.LVIF_IMAGE
		lc.IImage = int32(imIndex)
	}
}

func (control *ListView) AddItem(item ListItem) {
	control.InsertItem(item, control.ItemCount())
}

func (control *ListView) InsertItem(item ListItem, index int) {
	text := item.Text()
	li := &w32.LVITEM{
		Mask:    w32.LVIF_TEXT | w32.LVIF_PARAM,
		PszText: syscall.StringToUTF16Ptr(text[0]),
		IItem:   int32(index),
	}

	control.lastIndex++
	ix := new(int)
	*ix = control.lastIndex
	li.LParam = uintptr(*ix)
	control.handle2Item[li.LParam] = item
	control.item2Handle[item] = li.LParam

	control.applyImage(li, item.ImageIndex())
	control.insertLvItem(li)

	for i := 1; i < len(text); i++ {
		li.Mask = w32.LVIF_TEXT
		li.PszText = syscall.StringToUTF16Ptr(text[i])
		li.ISubItem = int32(i)
		control.setLvItem(li)
	}
}

func (control *ListView) UpdateItem(item ListItem) bool {
	lparam, ok := control.item2Handle[item]
	if !ok {
		return false
	}

	index := control.findIndexByItem(item)
	if index == -1 {
		return false
	}

	text := item.Text()
	li := &w32.LVITEM{
		Mask:    w32.LVIF_TEXT | w32.LVIF_PARAM,
		PszText: syscall.StringToUTF16Ptr(text[0]),
		LParam:  lparam,
		IItem:   int32(index),
	}

	control.applyImage(li, item.ImageIndex())
	control.setLvItem(li)

	for i := 1; i < len(text); i++ {
		li.Mask = w32.LVIF_TEXT
		li.PszText = syscall.StringToUTF16Ptr(text[i])
		li.ISubItem = int32(i)
		control.setLvItem(li)
	}
	return true
}

func (control *ListView) insertLvColumn(controlColumn *w32.LVCOLUMN, iCol int) {
	w32.SendMessage(control.hwnd, w32.LVM_INSERTCOLUMN, uintptr(iCol), uintptr(unsafe.Pointer(controlColumn)))
}

func (control *ListView) insertLvItem(controlItem *w32.LVITEM) {
	w32.SendMessage(control.hwnd, w32.LVM_INSERTITEM, 0, uintptr(unsafe.Pointer(controlItem)))
}

func (control *ListView) setLvItem(controlItem *w32.LVITEM) {
	w32.SendMessage(control.hwnd, w32.LVM_SETITEM, 0, uintptr(unsafe.Pointer(controlItem)))
}

func (control *ListView) DeleteAllItems() bool {
	if w32.SendMessage(control.hwnd, w32.LVM_DELETEALLITEMS, 0, 0) == w32.TRUE {
		control.item2Handle = make(map[ListItem]uintptr)
		control.handle2Item = make(map[uintptr]ListItem)
		return true
	}
	return false
}

func (control *ListView) DeleteItem(item ListItem) error {
	index := control.findIndexByItem(item)
	if index == -1 {
		return errors.New("item not found")
	}

	if w32.SendMessage(control.hwnd, w32.LVM_DELETEITEM, uintptr(index), 0) == 0 {
		return errors.New("SendMessage(TVM_DELETEITEM) failed")
	}

	h := control.item2Handle[item]
	delete(control.item2Handle, item)
	delete(control.handle2Item, h)
	return nil
}

func (control *ListView) findIndexByItem(item ListItem) int {
	lparam, ok := control.item2Handle[item]
	if !ok {
		return -1
	}

	it := &w32.LVFINDINFO{
		Flags:  w32.LVFI_PARAM,
		LParam: lparam,
	}
	var i int = -1
	return int(w32.SendMessage(control.hwnd, w32.LVM_FINDITEM, uintptr(i), uintptr(unsafe.Pointer(it))))
}

func (control *ListView) findItemByIndex(i int) ListItem {
	it := &w32.LVITEM{
		Mask:  w32.LVIF_PARAM,
		IItem: int32(i),
	}

	if w32.SendMessage(control.hwnd, w32.LVM_GETITEM, 0, uintptr(unsafe.Pointer(it))) == w32.TRUE {
		if item, ok := control.handle2Item[it.LParam]; ok {
			return item
		}
	}
	return nil
}

func (control *ListView) EnsureVisible(item ListItem) bool {
	if i := control.findIndexByItem(item); i != -1 {
		return w32.SendMessage(control.hwnd, w32.LVM_ENSUREVISIBLE, uintptr(i), 1) == 0
	}
	return false
}

func (control *ListView) SelectedItem() ListItem {
	if items := control.SelectedItems(); len(items) > 0 {
		return items[0]
	}
	return nil
}

func (control *ListView) SetSelectedItem(item ListItem) bool {
	if i := control.findIndexByItem(item); i > -1 {
		control.SetSelectedIndex(i)
		return true
	}
	return false
}

// mask is used to set the LVITEM.Mask for ListView.GetItem which indicates which attributes you'd like to receive
// of LVITEM.
func (control *ListView) SelectedItems() []ListItem {
	var items []ListItem

	var i int = -1
	for {
		if i = int(w32.SendMessage(control.hwnd, w32.LVM_GETNEXTITEM, uintptr(i), uintptr(w32.LVNI_SELECTED))); i == -1 {
			break
		}

		if item := control.findItemByIndex(i); item != nil {
			items = append(items, item)
		}
	}
	return items
}

func (control *ListView) SelectedCount() uint {
	return uint(w32.SendMessage(control.hwnd, w32.LVM_GETSELECTEDCOUNT, 0, 0))
}

// GetSelectedIndex first selected item index. Returns -1 if no item is selected.
func (control *ListView) SelectedIndex() int {
	var i int = -1
	return int(w32.SendMessage(control.hwnd, w32.LVM_GETNEXTITEM, uintptr(i), uintptr(w32.LVNI_SELECTED)))
}

// Set i to -1 to select all items.
func (control *ListView) SetSelectedIndex(i int) {
	control.setItemState(i, w32.LVIS_SELECTED, w32.LVIS_SELECTED)
}

func (control *ListView) SetImageList(imageList *ImageList) {
	w32.SendMessage(control.hwnd, w32.LVM_SETIMAGELIST, w32.LVSIL_SMALL, uintptr(imageList.Handle()))
	control.iml = imageList
}

// Event publishers
func (control *ListView) OnEndLabelEdit() *EventManager {
	return &control.onEndLabelEdit
}

func (control *ListView) OnDoubleClick() *EventManager {
	return &control.onDoubleClick
}

func (control *ListView) OnClick() *EventManager {
	return &control.onClick
}

func (control *ListView) OnRClick() *EventManager {
	return &control.onRClick
}

func (control *ListView) OnKeyDown() *EventManager {
	return &control.onKeyDown
}

func (control *ListView) OnItemChanging() *EventManager {
	return &control.onItemChanging
}

func (control *ListView) OnItemChanged() *EventManager {
	return &control.onItemChanged
}

func (control *ListView) OnCheckChanged() *EventManager {
	return &control.onCheckChanged
}

func (control *ListView) OnViewChange() *EventManager {
	return &control.onViewChange
}

func (control *ListView) OnEndScroll() *EventManager {
	return &control.onEndScroll
}

// Message processer
func (control *ListView) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	/*case w32.WM_ERASEBKGND:
	control.StretchLastColumn()
	println("case w32.WM_ERASEBKGND")
	return 1*/

	case w32.WM_NOTIFY:
		nm := (*w32.NMHDR)(unsafe.Pointer(lparam))
		code := int32(nm.Code)

		switch code {

		case w32.LVN_COLUMNCLICK:
			nm := (*w32.NMLISTVIEW)(unsafe.Pointer(lparam))
			control.onClick.Fire(NewEvent(control, ListViewEvent{int(nm.IItem), int(nm.ISubItem)}))

			if control.sort != nil {
				sortCol := int(nm.ISubItem)
				if sortCol == control.sortCol {
					control.sortAscending = !control.sortAscending
				} else {
					control.sortAscending = true
					control.sortCol = sortCol
				}
				control.sort(control.sortCol, control.sortAscending)
			}

		case w32.LVN_BEGINLABELEDITW:
			// println("Begin label edit")
		case w32.LVN_ENDLABELEDITW:
			nmdi := (*w32.NMLVDISPINFO)(unsafe.Pointer(lparam))
			if nmdi.Item.PszText != nil {
				if item, ok := control.handle2Item[nmdi.Item.LParam]; ok {
					control.onEndLabelEdit.Fire(NewEvent(control,
						&LabelEditEventData{Item: item,
							Text: w32.UTF16PtrToString(nmdi.Item.PszText)}))
				}
				return w32.TRUE
			}
		case w32.NM_DBLCLK:
			ac := (*w32.NMITEMACTIVATE)(unsafe.Pointer(lparam))
			control.onDoubleClick.Fire(NewEvent(control, ListViewEvent{int(ac.IItem), int(ac.ISubItem)}))

		case w32.NM_RCLICK:
			ac := (*w32.NMITEMACTIVATE)(unsafe.Pointer(lparam))
			control.onRClick.Fire(NewEvent(control, ListViewEvent{int(ac.IItem), int(ac.ISubItem)}))
		case w32.NM_CLICK:

			ac := (*w32.NMITEMACTIVATE)(unsafe.Pointer(lparam))
			var hti w32.LVHITTESTINFO
			hti.Pt = w32.POINT{X: ac.PtAction.X, Y: ac.PtAction.Y}
			w32.SendMessage(control.hwnd, w32.LVM_HITTEST, 0, uintptr(unsafe.Pointer(&hti)))

			if hti.Flags == w32.LVHT_ONITEMSTATEICON {
				if item := control.findItemByIndex(int(hti.IItem)); item != nil {
					if item, ok := item.(ListItemChecker); ok {
						checked := !item.Checked()
						item.SetChecked(checked)
						control.onCheckChanged.Fire(NewEvent(control, item))

						if w32.SendMessage(control.hwnd, w32.LVM_UPDATE, uintptr(hti.IItem), 0) == w32.FALSE {
							panic("SendMessage(LVM_UPDATE)")
						}
					}
				}
			}
			control.onClick.Fire(NewEvent(control, ListViewEvent{int(ac.IItem), int(ac.ISubItem)}))

		case w32.LVN_KEYDOWN:
			nmkey := (*w32.NMLVKEYDOWN)(unsafe.Pointer(lparam))
			if nmkey.WVKey == w32.VK_SPACE && control.CheckBoxes() {
				if item := control.SelectedItem(); item != nil {
					if item, ok := item.(ListItemChecker); ok {
						checked := !item.Checked()
						item.SetChecked(checked)
						control.onCheckChanged.Fire(NewEvent(control, item))
					}

					index := control.findIndexByItem(item)
					if w32.SendMessage(control.hwnd, w32.LVM_UPDATE, uintptr(index), 0) == w32.FALSE {
						panic("SendMessage(LVM_UPDATE)")
					}
				}
			}
			control.onKeyDown.Fire(NewEvent(control, nmkey.WVKey))
			key := nmkey.WVKey
			w32.SendMessage(control.Parent().Handle(), w32.WM_KEYDOWN, uintptr(key), 0)

		case w32.LVN_ITEMCHANGING:
			// This event also fires when listview has changed via code.
			nmcontrol := (*w32.NMLISTVIEW)(unsafe.Pointer(lparam))
			item := control.findItemByIndex(int(nmcontrol.IItem))
			control.onItemChanging.Fire(NewEvent(control, item))

		case w32.LVN_ITEMCHANGED:
			// This event also fires when listview has changed via code.
			nmcontrol := (*w32.NMLISTVIEW)(unsafe.Pointer(lparam))
			item := control.findItemByIndex(int(nmcontrol.IItem))
			control.onItemChanged.Fire(NewEvent(control, item))

		case w32.LVN_GETDISPINFO:
			nmdi := (*w32.NMLVDISPINFO)(unsafe.Pointer(lparam))
			if nmdi.Item.StateMask&w32.LVIS_STATEIMAGEMASK > 0 {
				if item, ok := control.handle2Item[nmdi.Item.LParam]; ok {
					if item, ok := item.(ListItemChecker); ok {

						checked := item.Checked()
						if checked {
							nmdi.Item.State = 0x2000
						} else {
							nmdi.Item.State = 0x1000
						}
					}
				}
			}

			control.onViewChange.Fire(NewEvent(control, nil))

		case w32.LVN_ENDSCROLL:
			control.onEndScroll.Fire(NewEvent(control, nil))
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

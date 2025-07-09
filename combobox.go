/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"syscall"
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

/* DiffComboBoxable
 *
 */
type DiffComboBoxable interface {
	DeleteAllItems() bool
	AddItem(str string) bool
	InsertItem(index int, str string) bool
	DeleteItem(index int) bool
	GetItemLength(index int) int
	GetItem(index int) string
	SelectedItem() int
	GetSelectedItem() string
	SetSelectedItem(value int) bool
	ShowDropdown(show bool)
	OnSelectedChange() *EventManager
	OnChange() *EventManager
	OnUpdate() *EventManager
	OnSelectedEnd() *EventManager
}

/* ComboBoxable
 *
 */
type ComboBoxable interface {
	BaseController
	DiffComboBoxable
}

/* ComboBox
 *
 */
type ComboBox struct {
	ControlBase
	onSelectedChange EventManager
	onChange         EventManager
	onUpdate         EventManager
	onSelectedEnd    EventManager
	*Edit
}

func NewComboBox(parent Controller) *ComboBox {
	return NewComboBoxWithFlags(parent, 0)
}

func NewComboBoxWithFlags(parent Controller, style uint) *ComboBox {
	control := new(ComboBox)

	control.InitControl(w32.WC_COMBOBOX, parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_VSCROLL|w32.CBS_AUTOHSCROLL|w32.CBS_DROPDOWN|style)
	RegMsgHandler(control)

	// register edit control also
	edit_control := new(Edit)
	edit_control.parent = control

	// edit_control.hwnd = w32.ChildWindowFromPoint(control.hwnd, 1, 1)
	edit_control.hwnd = w32.FindWindowEx(control.hwnd, 0, w32.WC_EDIT, "")
	RegMsgHandler(edit_control)
	control.Edit = edit_control

	control.SetFont(DefaultFont)
	control.SetSize(200, 400)
	return control
}

func NewListComboBox(parent Controller) *ComboBox {
	control := new(ComboBox)

	control.InitControl(w32.WC_COMBOBOX, parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_VSCROLL|w32.CBS_DROPDOWNLIST)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetSize(200, 400)
	return control
}

func (control *ComboBox) DeleteAllItems() bool {
	return w32.SendMessage(control.hwnd, w32.CB_RESETCONTENT, 0, 0) == w32.TRUE
}

func (control *ComboBox) AddItem(str string) bool {
	lp := uintptr(unsafe.Pointer(w32.UTF16PtrFromString(str)))
	return w32.SendMessage(control.hwnd, w32.CB_ADDSTRING, 0, lp) != w32.CB_ERR
}

func (control *ComboBox) InsertItem(index int, str string) bool {
	lp := uintptr(unsafe.Pointer(w32.UTF16PtrFromString(str)))
	return w32.SendMessage(control.hwnd, w32.CB_INSERTSTRING, uintptr(index), lp) != w32.CB_ERR
}

func (control *ComboBox) DeleteItem(index int) bool {
	return w32.SendMessage(control.hwnd, w32.CB_DELETESTRING, uintptr(index), 0) != w32.CB_ERR
}

func (control *ComboBox) GetItemLength(index int) int {
	return int(w32.SendMessage(control.hwnd, w32.CB_GETLBTEXTLEN, uintptr(index), 0))
}

func (control *ComboBox) GetItem(index int) string {
	textLen := control.GetItemLength(index) + 1

	buf := make([]uint16, textLen)

	w32.SendMessage(control.hwnd, w32.CB_GETLBTEXT, uintptr(index), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

func (control *ComboBox) SelectedItem() int {
	return int(int32(w32.SendMessage(control.hwnd, w32.CB_GETCURSEL, 0, 0)))
}

func (control *ComboBox) GetSelectedItem() string {
	selected := control.SelectedItem()
	if selected < 0 {
		return ""
	}
	return control.GetItem(selected)
}

func (control *ComboBox) SetSelectedItem(value int) bool {
	return int(int32(w32.SendMessage(control.hwnd, w32.CB_SETCURSEL, uintptr(value), 0))) == value
}

func (control *ComboBox) ShowDropdown(show bool) {
	w32_show := w32.FALSE
	if show {
		w32_show = w32.TRUE
	}
	w32.SendMessage(control.hwnd, w32.CB_SHOWDROPDOWN, uintptr(w32_show), 0)
}

func (control *ComboBox) OnSelectedChange() *EventManager {
	return &control.onSelectedChange
}

func (control *ComboBox) OnChange() *EventManager {
	return &control.onChange
}

func (control *ComboBox) OnUpdate() *EventManager {
	return &control.onUpdate
}

func (control *ComboBox) OnSelectedEnd() *EventManager {
	return &control.onSelectedEnd
}

// Message processor
func (control *ComboBox) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		code := w32.HIWORD(uint32(wparam))

		switch code {
		case w32.CBN_SELCHANGE:
			control.onSelectedChange.Fire(NewEvent(control, nil))
		case w32.CBN_SETFOCUS:
			control.onSetFocus.Fire(NewEvent(control, nil))
		case w32.CBN_KILLFOCUS:
			control.onKillFocus.Fire(NewEvent(control, nil))
		case w32.CBN_EDITUPDATE:
			control.onUpdate.Fire(NewEvent(control, nil))
		case w32.CBN_EDITCHANGE:
			control.onChange.Fire(NewEvent(control, nil))
		case w32.CBN_SELENDOK:
			control.onSelectedEnd.Fire(NewEvent(control, nil))
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
	//return control.W32Control.WndProc(msg, wparam, lparam)
}

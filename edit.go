/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

/* DiffEditable
 *
 */
type DiffEditable interface {
	OnChange() *EventManager
	SetReadOnly(isReadOnly bool)
	Modified() uintptr
	SetModified(modified bool) uintptr
	Selected() (int, int)
	SelectText(start, end int)
	SetPassword(isPassword bool)
}

/* Editable
 *
 */
type Editable interface {
	BaseController
	DiffEditable
}

/* Edit
 *
 */
type Edit struct {
	ControlBase
	onChange EventManager
}

const passwordChar = '*'
const nopasswordChar = ' '

func NewEdit(parent Controller) *Edit {
	control := new(Edit)

	control.InitControl(w32.WC_EDIT, parent, w32.WS_EX_CLIENTEDGE, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.ES_LEFT|
		w32.ES_AUTOHSCROLL)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetSize(200, 22)
	return control
}

// Events.
func (control *Edit) OnChange() *EventManager {
	return &control.onChange
}

// Public methods.
func (control *Edit) SetReadOnly(isReadOnly bool) {
	w32.SendMessage(control.hwnd, w32.EM_SETREADONLY, uintptr(w32.BoolToBOOL(isReadOnly)), 0)
}

func (control *Edit) Modified() uintptr {
	return w32.SendMessage(control.hwnd, w32.EM_GETMODIFY, 0, 0)
}

func (control *Edit) SetModified(modified bool) uintptr {
	w32_modified := w32.FALSE
	if modified {
		w32_modified = w32.TRUE
	}
	return w32.SendMessage(control.hwnd, w32.EM_SETMODIFY, uintptr(w32_modified), 0)
}

func (control *Edit) Selected() (int, int) {
	var start, end int
	w32.SendMessage(control.hwnd, w32.EM_GETSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return start, end
}

func (control *Edit) SelectText(start, end int) {
	w32.SendMessage(control.hwnd, w32.EM_SETSEL, uintptr(start), uintptr(end))
}

func (control *Edit) SetPassword(isPassword bool) {
	if isPassword {
		w32.SendMessage(control.hwnd, w32.EM_SETPASSWORDCHAR, uintptr(passwordChar), 0)
	} else {
		w32.SendMessage(control.hwnd, w32.EM_SETPASSWORDCHAR, 0, 0)
	}
}

func (control *Edit) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		switch w32.HIWORD(uint32(wparam)) {
		case w32.EN_SETFOCUS:
			control.onSetFocus.Fire(NewEvent(control, nil))
		case w32.EN_KILLFOCUS:
			control.onKillFocus.Fire(NewEvent(control, nil))
		case w32.EN_CHANGE:
			if control.Modified() == w32.TRUE {
				control.onChange.Fire(NewEvent(control, nil))
			}
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

// MultiEdit is multiline text edit.
type MultiEdit struct {
	ControlBase
	onChange EventManager
}

func NewMultiEdit(parent Controller) *MultiEdit {
	control := new(MultiEdit)
	control.InitControl(w32.WC_EDIT, parent, w32.WS_EX_CLIENTEDGE, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.ES_LEFT|
		w32.WS_VSCROLL|w32.WS_HSCROLL|w32.ES_MULTILINE|w32.ES_WANTRETURN|w32.ES_AUTOHSCROLL|w32.ES_AUTOVSCROLL)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetSize(200, 400)
	return control
}

// Events
func (control *MultiEdit) OnChange() *EventManager {
	return &control.onChange
}

// Public methods
func (control *MultiEdit) SetReadOnly(isReadOnly bool) {
	w32.SendMessage(control.hwnd, w32.EM_SETREADONLY, uintptr(w32.BoolToBOOL(isReadOnly)), 0)
}

func (control *MultiEdit) AddLine(text string) {
	if len(control.Text()) == 0 {
		control.SetText(text)
	} else {
		control.SetText(control.Text() + "\r\n" + text)
	}
}

func (control *MultiEdit) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {

	case w32.WM_COMMAND:
		switch w32.HIWORD(uint32(wparam)) {
		case w32.EN_CHANGE:
			control.onChange.Fire(NewEvent(control, nil))
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

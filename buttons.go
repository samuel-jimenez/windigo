/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"fmt"

	"github.com/samuel-jimenez/windigo/w32"
)

type Button struct {
	ControlBase
	onClick EventManager
}

func (control *Button) OnClick() *EventManager {
	return &control.onClick
}

func (control *Button) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		control.onClick.Fire(NewEvent(control, nil))
		/*case w32.WM_LBUTTONDOWN:
			w32.SetCapture(control.Handle())
		case w32.WM_LBUTTONUP:
			w32.ReleaseCapture()*/
		/*case win.WM_GETDLGCODE:
		println("GETDLGCODE")*/
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
	//return control.W32Control.WndProc(msg, wparam, lparam)
}

func (control *Button) Checked() bool {
	result := w32.SendMessage(control.hwnd, w32.BM_GETCHECK, 0, 0)
	return result == w32.BST_CHECKED
}

func (control *Button) SetChecked(checked bool) {
	wparam := w32.BST_CHECKED
	if !checked {
		wparam = w32.BST_UNCHECKED
	}
	w32.SendMessage(control.hwnd, w32.BM_SETCHECK, uintptr(wparam), 0)
}

// SetIcon sets icon on the button. Recommended icons are 32x32 with 32bit color depth.
func (control *Button) SetIcon(ico *Icon) {
	w32.SendMessage(control.hwnd, w32.BM_SETIMAGE, w32.IMAGE_ICON, uintptr(ico.handle))
}

func (control *Button) SetResIcon(iconID uint16) {
	if ico, err := NewIconFromResource(GetAppInstance(), iconID); err == nil {
		control.SetIcon(ico)
		return
	}
	panic(fmt.Sprintf("missing icon with icon ID: %d", iconID))
}

type PushButton struct {
	Button
}

func NewPushButton(parent Controller) *PushButton {
	control := new(PushButton)

	control.InitControl("BUTTON", parent, 0, w32.BS_PUSHBUTTON|w32.WS_TABSTOP|w32.WS_VISIBLE|w32.WS_CHILD)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("Button")
	control.SetSize(100, 22)

	return control
}

// SetDefault is used for dialogs to set default button.
func (control *PushButton) SetDefault() {
	control.SetAndClearStyleBits(w32.BS_DEFPUSHBUTTON, w32.BS_PUSHBUTTON)
}

// IconButton does not display text, requires SetResIcon call.
type IconButton struct {
	Button
}

func NewIconButton(parent Controller) *IconButton {
	control := new(IconButton)

	control.InitControl("BUTTON", parent, 0, w32.BS_ICON|w32.WS_TABSTOP|w32.WS_VISIBLE|w32.WS_CHILD)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	// even if text would be set it would not be displayed
	control.SetText("")
	control.SetSize(100, 22)

	return control
}

type CheckBox struct {
	Button
}

func NewCheckBox(parent Controller) *CheckBox {
	control := new(CheckBox)

	control.InitControl("BUTTON", parent, 0, w32.WS_TABSTOP|w32.WS_VISIBLE|w32.WS_CHILD|w32.BS_AUTOCHECKBOX)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("CheckBox")
	control.SetSize(100, 22)

	return control
}

type RadioButton struct {
	Button
}

func NewRadioButton(parent Controller) *RadioButton {
	control := new(RadioButton)

	control.InitControl("BUTTON", parent, 0, w32.WS_TABSTOP|w32.WS_VISIBLE|w32.WS_CHILD|w32.BS_AUTORADIOBUTTON)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("RadioButton")
	control.SetSize(100, 22)

	return control
}

type GroupBox struct {
	Button
}

func NewGroupBox(parent Controller) *GroupBox {
	control := new(GroupBox)

	control.InitControl("BUTTON", parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_GROUP|w32.BS_GROUPBOX)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("GroupBox")
	control.SetSize(100, 100)

	return control

}

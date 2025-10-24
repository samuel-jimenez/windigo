/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"github.com/samuel-jimenez/windigo/w32"
)

// Dialog displayed as z-order top window until closed.
// It also disables parent window so it can not be clicked.
type Dialog struct {
	Form
	isModal,
	isShown bool

	btnOk     *PushButton
	btnCancel *PushButton

	onLoad   EventManager
	onOk     EventManager
	onCancel EventManager
}

func NewDialog(parent Controller) *Dialog {
	control := new(Dialog)

	control.isForm = true
	control.isModal = true
	RegClassOnlyOnce("windigo_Dialog")

	control.hwnd = CreateWindow("windigo_Dialog", parent, w32.WS_EX_CONTROLPARENT, /* IMPORTANT */
		w32.WS_SYSMENU|w32.WS_CAPTION|w32.WS_THICKFRAME /*|w32.WS_BORDER|w32.WS_POPUP*/)
	control.parent = parent
	control.local_shortcuts = make(map[Shortcut]func() bool)

	// this might fail if icon resource is not embedded in the binary
	if ico, err := NewIconFromResource(GetAppInstance(), uint16(AppIconID)); err == nil {
		control.SetIcon(0, ico)
	}

	// Dlg forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(control.hwnd, w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("Form")
	control.SetSize(200, 100)
	return control
}

func (control *Dialog) SetModal(modal bool) {
	if !control.isShown {
		control.isModal = modal
	}
}

// SetButtons wires up dialog events to buttons. btnCancel can be nil.
func (control *Dialog) SetButtons(btnOk *PushButton, btnCancel *PushButton) {
	control.btnOk = btnOk
	control.btnOk.SetDefault()
	control.btnCancel = btnCancel
}

// Events
func (control *Dialog) OnLoad() *EventManager {
	return &control.onLoad
}

func (control *Dialog) OnOk() *EventManager {
	return &control.onOk
}

func (control *Dialog) OnCancel() *EventManager {
	return &control.onCancel
}

// PreTranslateMessage handles dialog specific messages. IMPORTANT.
func (control *Dialog) PreTranslateMessage(msg *w32.MSG) bool {
	switch msg.Message {
	case w32.WM_KEYDOWN:
		// Shortcut support.
		key := Key(msg.WParam)
		if uint32(msg.LParam)>>30 == 0 {
			shortcut := Shortcut{ModifiersDown(), key}
			if action, ok := control.local_shortcuts[shortcut]; ok {
				return action()
			} else if action, ok := shortcut2Action[shortcut]; ok { // Access menu shortcuts. This may not be what we want. TODO check isModal?
				if action.Enabled() {
					action.onClick.Fire(NewEvent(control, msg.Pt))
					return true
				}
			}
		}
	}
	if msg.Message >= w32.WM_KEYFIRST && msg.Message <= w32.WM_KEYLAST && w32.IsDialogMessage(control.hwnd, msg) {
		return true
	}
	return false
}

// Show dialog performs special setup for dialog windows.
func (control *Dialog) Show() {
	if control.isModal {
		control.Parent().SetEnabled(false)
	}
	control.isShown = true
	control.onLoad.Fire(NewEvent(control, nil))
	control.Form.Show()
}

// Close dialog when you done with it.
func (control *Dialog) Close() {
	if control.isModal {
		control.Parent().SetEnabled(true)
	}
	control.ControlBase.Close()
}

func (control *Dialog) cancel() {
	if control.btnCancel != nil {
		control.btnCancel.onClick.Fire(NewEvent(control.btnCancel, nil))
	}
	control.onCancel.Fire(NewEvent(control, nil))
}

func (control *Dialog) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_COMMAND:
		switch w32.LOWORD(uint32(wparam)) {
		case w32.IDOK:
			if control.btnOk != nil {
				control.btnOk.onClick.Fire(NewEvent(control.btnOk, nil))
			}
			control.onOk.Fire(NewEvent(control, nil))
			return w32.TRUE

		case w32.IDCANCEL:
			control.cancel()
			return w32.TRUE
		}
	case w32.WM_CLOSE:
		control.cancel() // use onCancel or control.btnCancel.OnClick to close
		return 0

	case w32.WM_DESTROY:
		if control.isModal {
			control.Parent().SetEnabled(true)
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

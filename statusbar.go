/*
 * Copyright (C) 2025 The Windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"github.com/samuel-jimenez/windigo/w32"
)

type StatusBar struct {
	ControlBase
}

func NewStatusBar(parent Controller) *StatusBar {
	control := new(StatusBar)

	control.InitControl(w32.STATUSCLASSNAME, parent, 0, w32.SBARS_TOOLTIPS|w32.WS_CHILD|w32.WS_VISIBLE|w32.CCS_BOTTOM) // CCS_BOTTOM disables w32.SBARS_SIZEGRIP
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	return control
}

func (control *StatusBar) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

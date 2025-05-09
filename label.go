/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"github.com/samuel-jimenez/windigo/w32"
)

type Label struct {
	ControlBase
}

func NewLabel(parent Controller) *Label {
	lb := new(Label)

	lb.InitControl(w32.WC_STATIC, parent, 0, w32.WS_CHILD|w32.WS_VISIBLE|w32.SS_LEFTNOWORDWRAP)
	RegMsgHandler(lb)

	lb.SetFont(DefaultFont)
	lb.SetText("Label")
	lb.SetSize(100, 25)
	return lb
}

func (lb *Label) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	return w32.DefWindowProc(lb.hwnd, msg, wparam, lparam)
}

/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"syscall"
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

// TabView creates MultiPanel internally and manages tabs as panels.
type TabView struct {
	ControlBase

	panels           *MultiPanel
	onSelectedChange EventManager
}

func NewTabView(parent Controller) *TabView {
	control := new(TabView)

	control.InitControl(w32.WC_TABCONTROL, parent, 0,
		w32.WS_CHILD|w32.WS_VISIBLE|w32.WS_TABSTOP|w32.WS_CLIPSIBLINGS)
	RegMsgHandler(control)

	control.panels = NewMultiPanel(parent)

	control.SetFont(DefaultFont)
	control.SetSize(200, 24)
	return control
}

func (control *TabView) Panels() *MultiPanel {
	return control.panels
}

func (control *TabView) tcitemFromPage(panel Pane) *w32.TCITEM {
	text := syscall.StringToUTF16(panel.Text())
	item := &w32.TCITEM{
		Mask:       w32.TCIF_TEXT,
		PszText:    &text[0],
		CchTextMax: int32(len(text)),
	}
	return item
}

func (control *TabView) AddPanel(text string) *Panel {
	panel := NewPanel(control.panels)
	panel.SetText(text)

	item := control.tcitemFromPage(panel)
	index := control.panels.Count()
	idx := int(w32.SendMessage(control.hwnd, w32.TCM_INSERTITEM, uintptr(index), uintptr(unsafe.Pointer(item))))
	if idx == -1 {
		panic("SendMessage(TCM_INSERTITEM) failed")
	}

	control.panels.AddPanel(panel)
	control.SetCurrent(idx)
	return panel
}

func (control *TabView) AddAutoPanel(text string) *AutoPanel {
	panel := NewAutoPanel(control.panels)
	panel.SetText(text)

	item := control.tcitemFromPage(panel)
	index := control.panels.Count()
	idx := int(w32.SendMessage(control.hwnd, w32.TCM_INSERTITEM, uintptr(index), uintptr(unsafe.Pointer(item))))
	if idx == -1 {
		panic("SendMessage(TCM_INSERTITEM) failed")
	}

	control.panels.AddPanel(panel)
	control.SetCurrent(idx)
	return panel
}

func (control *TabView) DeletePanel(index int) {
	w32.SendMessage(control.hwnd, w32.TCM_DELETEITEM, uintptr(index), 0)
	control.panels.DeletePanel(index)
	switch {
	case control.panels.Count() > index:
		control.SetCurrent(index)
	case control.panels.Count() == 0:
		control.SetCurrent(0)
	}
}

func (control *TabView) Current() int {
	return control.panels.Current()
}

func (control *TabView) SetCurrent(index int) {
	if index < 0 || index >= control.panels.Count() {
		panic("invalid index")
	}
	if ret := int(w32.SendMessage(control.hwnd, w32.TCM_SETCURSEL, uintptr(index), 0)); ret == -1 {
		panic("SendMessage(TCM_SETCURSEL) failed")
	}
	control.panels.SetCurrent(index)

	// status bar disappears otherwise
	control.parent.RefreshStatusBar()
}

func (control *TabView) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_NOTIFY:
		nmhdr := (*w32.NMHDR)(unsafe.Pointer(lparam))

		switch int32(nmhdr.Code) {
		case w32.TCN_SELCHANGE:
			cur := int(w32.SendMessage(control.hwnd, w32.TCM_GETCURSEL, 0, 0))
			control.SetCurrent(cur)

			control.onSelectedChange.Fire(NewEvent(control, nil))
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

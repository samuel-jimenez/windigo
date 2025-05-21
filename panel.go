/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 */

package windigo

import (
	"fmt"

	"slices"

	"github.com/samuel-jimenez/windigo/w32"
)

type Pane interface {
	DockAllow
	Controller
}

type Panel struct {
	ControlBase
	layoutMng LayoutManager
}

func NewPanel(parent Controller) *Panel {
	control := new(Panel)

	RegClassOnlyOnce("windigo_Panel")
	control.hwnd = CreateWindow("windigo_Panel", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	control.parent = parent
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("")
	control.SetSize(200, 65)
	return control
}

// SetLayout panel implements DockAllow interface.
func (control *Panel) SetLayout(mng LayoutManager) {
	control.layoutMng = mng
}

func (control *Panel) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_SIZE, w32.WM_PAINT:
		if control.layoutMng != nil {
			control.layoutMng.Update()
		}
	case w32.WM_ERASEBKGND:
		if control.border_pen != nil {
			canvas := NewCanvasFromHDC(w32.HDC(wparam))
			margin_0 := int32(control.border_pen.width) / 2
			margin_1 := int32(control.border_pen.width+1) / 2

			r := control.Bounds()
			r.rect.Top += margin_0
			r.rect.Left += margin_0
			r.rect.Bottom -= margin_1
			r.rect.Right -= margin_1

			if control.erasure_pen != nil {
				canvas.DrawFillRect(r, control.erasure_pen, NewSystemColorBrush(w32.COLOR_BTNFACE))

				control.erasure_pen.Dispose()
				control.erasure_pen = nil
				w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
			}
			canvas.DrawFillRect(r, control.border_pen, NewSystemColorBrush(w32.COLOR_BTNFACE))
			canvas.Dispose()
			return 1
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

type GroupPanel struct {
	Panel
	groupbox *GroupBox
}

func NewGroupPanel(parent Controller) *GroupPanel {
	control := new(GroupPanel)

	RegClassOnlyOnce("windigo_GroupPanel")
	control.hwnd = CreateWindow("windigo_GroupPanel", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	control.parent = parent
	RegMsgHandler(control)
	control.groupbox = NewGroupBox(control)
	control.SetFont(DefaultFont)
	control.SetText("")
	control.SetSize(200, 65)
	return control
}

func (control *GroupPanel) SetText(caption string) {
	control.groupbox.SetText(caption)
}

func (control *GroupPanel) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_SIZE, w32.WM_PAINT:
		if control.groupbox != nil {
			var margin int
			if control.border_pen != nil {
				margin = int(control.border_pen.width)
				control.groupbox.SetPos(margin, margin)
				margin = 2*margin + 1
			}
			control.groupbox.SetSize(control.ClientWidth()-margin, control.ClientHeight()-margin)
		}
	}
	return control.Panel.WndProc(msg, wparam, lparam)
}

var errorPanelPen = NewPen(w32.PS_GEOMETRIC, 2, NewSolidColorBrush(RGB(255, 128, 128)))
var errorPanelOkPen = NewPen(w32.PS_GEOMETRIC, 2, NewSolidColorBrush(RGB(220, 220, 220)))

// ErrorPanel shows errors or important messages.
// It is meant to stand out of other on screen controls.
type ErrorPanel struct {
	ControlBase
	pen    *Pen
	margin int
}

// NewErrorPanel.
func NewErrorPanel(parent Controller) *ErrorPanel {
	control := new(ErrorPanel)
	control.init(parent)

	control.SetFont(DefaultFont)
	control.SetText("No errors")
	control.SetSize(200, 65)
	control.margin = 5
	control.pen = errorPanelOkPen
	return control
}

func (control *ErrorPanel) init(parent Controller) {
	RegClassOnlyOnce("windigo_ErrorPanel")

	control.hwnd = CreateWindow("windigo_ErrorPanel", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	control.parent = parent

	RegMsgHandler(control)
}

func (control *ErrorPanel) SetMargin(margin int) {
	control.margin = margin
}

func (control *ErrorPanel) Printf(format string, v ...any) {
	control.SetText(fmt.Sprintf(format, v...))
	control.ShowAsError(false)
}

func (control *ErrorPanel) Errorf(format string, v ...any) {
	control.SetText(fmt.Sprintf(format, v...))
	control.ShowAsError(true)
}

func (control *ErrorPanel) ShowAsError(show bool) {
	if show {
		control.pen = errorPanelPen
	} else {
		control.pen = errorPanelOkPen
	}
	control.Invalidate(true)
}

func (control *ErrorPanel) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_ERASEBKGND:
		canvas := NewCanvasFromHDC(w32.HDC(wparam))
		r := control.Bounds()
		r.rect.Left += int32(control.margin)
		r.rect.Right -= int32(control.margin)
		r.rect.Top += int32(control.margin)
		r.rect.Bottom -= int32(control.margin)
		// old code used NewSystemColorBrush(w32.COLOR_BTNFACE)
		canvas.DrawFillRect(r, control.pen, NewSystemColorBrush(w32.COLOR_WINDOW))

		r.rect.Left += 5
		canvas.DrawText(control.Text(), r, 0, control.Font(), RGB(0, 0, 0))
		canvas.Dispose()
		return 1
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

// MultiPanel contains other panes and only makes one of them visible.
type MultiPanel struct {
	ControlBase
	current int
	panels  []Pane
}

func NewMultiPanel(parent Controller) *MultiPanel {
	control := new(MultiPanel)

	RegClassOnlyOnce("windigo_MultiPanel")
	control.hwnd = CreateWindow("windigo_MultiPanel", parent, w32.WS_EX_CONTROLPARENT, w32.WS_CHILD|w32.WS_VISIBLE)
	control.parent = parent
	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("")
	control.SetSize(300, 200)
	control.current = -1
	return control
}

func (control *MultiPanel) Count() int { return len(control.panels) }

// AddPanel adds panels to the internal list, first panel is visible all others are hidden.
func (control *MultiPanel) AddPanel(panel Pane) {
	if len(control.panels) > 0 {
		panel.Hide()
	}
	control.current = 0
	control.panels = append(control.panels, panel)
}

// ReplacePanel replaces panel, useful for refreshing controls on screen.
func (control *MultiPanel) ReplacePanel(index int, panel Pane) {
	control.panels[index] = panel
}

// DeletePanel removed panel.
func (control *MultiPanel) DeletePanel(index int) {
	control.panels = slices.Delete(control.panels, index, index+1)
}

func (control *MultiPanel) Current() int {
	return control.current
}

func (control *MultiPanel) SetCurrent(index int) {
	if index >= len(control.panels) {
		panic("index greater than number of panels")
	}
	if control.current == -1 {
		panic("no current panel, add panels first")
	}
	for i := range control.panels {
		if i != index {
			control.panels[i].Hide()
			control.panels[i].Invalidate(true)
		}
	}
	control.panels[index].Show()
	control.panels[index].Invalidate(true)
	control.current = index
}

func (control *MultiPanel) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_SIZE:
		// resize contained panels
		for _, p := range control.panels {
			p.SetPos(0, 0)
			p.SetSize(control.Size())
		}
	}
	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

func genPoint(p uintptr) (x, y int) {
	x = int(w32.LOWORD(uint32(p)))
	y = int(w32.HIWORD(uint32(p)))
	return
}

func genMouseEventArg(wparam, lparam uintptr) *MouseEventData {
	var data MouseEventData
	data.Button = int(wparam)
	data.X, data.Y = genPoint(lparam)

	return &data
}

func genDropFilesEventArg(wparam uintptr) *DropFilesEventData {
	hDrop := w32.HDROP(wparam)

	var data DropFilesEventData
	_, fileCount := w32.DragQueryFile(hDrop, 0xFFFFFFFF)
	data.Files = make([]string, fileCount)

	var i uint
	for i = range fileCount {
		data.Files[i], _ = w32.DragQueryFile(hDrop, i)
	}

	data.X, data.Y, _ = w32.DragQueryPoint(hDrop)
	w32.DragFinish(hDrop)
	return &data
}

func generalWndProc(hwnd w32.HWND, msg uint32, wparam, lparam uintptr) uintptr {

	switch msg {
	case w32.WM_HSCROLL:
		//println("case w32.WM_HSCROLL")

	case w32.WM_VSCROLL:
		//println("case w32.WM_VSCROLL")
	}

	if controller := GetMsgHandler(hwnd); controller != nil {
		ret := controller.WndProc(msg, wparam, lparam)

		switch msg {
		//TODO Why are we GetMsgHandlering twice?
		// TODO because we call the controller controller.WndProc first?
		case w32.WM_CTLCOLORSTATIC, //Label
			w32.WM_CTLCOLOREDIT, //Edit
			// w32.WM_CTLCOLORBTN, // needs custom draw
			// // stackoverflow.com/questions/75478704/how-to-change-the-color-of-a-button-using-wm-ctlcolorbtn
			w32.WM_CTLCOLORLISTBOX:
			if controller := GetMsgHandler(w32.HWND(lparam)); controller != nil {
				if controller.HasFGColor() {
					w32.SetTextColor(
						w32.HDC(wparam),
						w32.COLORREF(controller.FGColor()),
					)
				}
				if controller.HasHighlightColor() {
					w32.SetBkColor(
						w32.HDC(wparam),
						w32.COLORREF(controller.HighlightColor()),
					)
				}
				if controller.HasBGColor() {
					if !controller.HasHighlightColor() {
						w32.SetBkColor(
							w32.HDC(wparam),
							w32.COLORREF(controller.BGColor()),
						)
					}
					return controller.BGBrush().GetHBRUSH()
				}
			}
		case w32.WM_NOTIFY: //Reflect notification to control
			nm := (*w32.NMHDR)(unsafe.Pointer(lparam))

			// case w32.NM_CUSTOMDRAW:
			// log.Println("NM_CUSTOMDRAW controller", controller.Speak(), msg, ret)
			if controller := GetMsgHandler(nm.HwndFrom); controller != nil {
				ret := controller.WndProc(msg, wparam, lparam)
				if ret != 0 {
					w32.SetWindowLong(hwnd, w32.DWL_MSGRESULT, uint32(ret))
					return w32.TRUE
				}
			}
		case w32.WM_COMMAND:
			if lparam != 0 { //Reflect message to control
				h := w32.HWND(lparam)
				if controller := GetMsgHandler(h); controller != nil {
					ret := controller.WndProc(msg, wparam, lparam)
					if ret != 0 {
						w32.SetWindowLong(hwnd, w32.DWL_MSGRESULT, uint32(ret))
						return w32.TRUE
					}
				}
			}
		case w32.WM_CLOSE:
			controller.OnClose().Fire(NewEvent(controller, nil))
		case w32.WM_KILLFOCUS:
			controller.OnKillFocus().Fire(NewEvent(controller, nil))
		case w32.WM_SETFOCUS:
			controller.OnSetFocus().Fire(NewEvent(controller, nil))
		case w32.WM_DROPFILES:
			controller.OnDropFiles().Fire(NewEvent(controller, genDropFilesEventArg(wparam)))
		case w32.WM_CONTEXTMENU:
			if wparam != 0 { //Reflect message to control
				h := w32.HWND(wparam)
				if controller := GetMsgHandler(h); controller != nil {
					contextMenu := controller.ContextMenu()
					x, y := genPoint(lparam)

					if contextMenu != nil {
						id := w32.TrackPopupMenuEx(
							contextMenu.hMenu,
							w32.TPM_NOANIMATION|w32.TPM_RETURNCMD,
							int32(x),
							int32(y),
							controller.Handle(),
							nil)

						item := findMenuItemByID(int(id))
						if item != nil {
							item.OnClick().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
						}
						return 0
					}
				}
			}

		case w32.WM_LBUTTONDOWN:
			controller.OnLBDown().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_LBUTTONUP:
			controller.OnLBUp().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_LBUTTONDBLCLK:
			controller.OnLBDbl().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MBUTTONDOWN:
			controller.OnMBDown().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MBUTTONUP:
			controller.OnMBUp().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONDOWN:
			controller.OnRBDown().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONUP:
			controller.OnRBUp().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONDBLCLK:
			controller.OnRBDbl().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MOUSEMOVE:
			controller.OnMouseMove().Fire(NewEvent(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_PAINT:
			canvas := NewCanvasFromHwnd(hwnd)
			defer canvas.Dispose()
			controller.drawBorder(canvas)
			controller.OnPaint().Fire(NewEvent(controller, &PaintEventData{Canvas: canvas}))
			return 0

		case w32.WM_KEYUP:
			controller.OnKeyUp().Fire(NewEvent(controller, &KeyUpEventData{int(wparam), int(lparam)}))
		case w32.WM_SIZE:
			x, y := genPoint(lparam)
			controller.OnSize().Fire(NewEvent(controller, &SizeEventData{uint(wparam), x, y}))
		case w32.WM_ERASEBKGND:
			if controller.HasBGColor() {
				canvas := NewCanvasFromHDC(w32.HDC(wparam))
				defer canvas.Dispose()
				bounding_rect := controller.WindowBounds()
				canvas.FillRect(bounding_rect, controller.BGBrush())
			}
			return 1
		}
		return ret
	}

	return w32.DefWindowProc(hwnd, uint32(msg), wparam, lparam)
}

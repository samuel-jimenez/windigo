/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

type LayoutManager interface {
	Update()
}

// A Form is main window of the application.
type Form struct {
	ControlBase

	layoutMng LayoutManager

	// Fullscreen / Unfullscreen
	isFullscreen            bool
	previousWindowStyle     uint32
	previousWindowPlacement w32.WINDOWPLACEMENT
}

func NewCustomForm(parent Controller, exStyle int, dwStyle uint) *Form {
	control := new(Form)

	RegClassOnlyOnce("windigo_Form")

	control.isForm = true

	if exStyle == 0 {
		exStyle = w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
	}

	if dwStyle == 0 {
		dwStyle = w32.WS_OVERLAPPEDWINDOW
	}

	control.hwnd = CreateWindow("windigo_Form", parent, uint(exStyle), dwStyle)
	control.parent = parent

	// this might fail if icon resource is not embedded in the binary
	if ico, err := NewIconFromResource(GetAppInstance(), uint16(AppIconID)); err == nil {
		control.SetIcon(0, ico)
	}

	// This forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(control.hwnd, w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)

	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("Form")
	return control
}

func NewForm(parent Controller) *Form {
	control := new(Form)

	RegClassOnlyOnce("windigo_Form")

	control.isForm = true
	control.hwnd = CreateWindow("windigo_Form", parent, w32.WS_EX_CONTROLPARENT|w32.WS_EX_APPWINDOW, w32.WS_OVERLAPPEDWINDOW)
	control.parent = parent

	// this might fail if icon resource is not embedded in the binary
	if ico, err := NewIconFromResource(GetAppInstance(), uint16(AppIconID)); err == nil {
		control.SetIcon(0, ico)
	}

	// This forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(control.hwnd, w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)

	RegMsgHandler(control)

	control.SetFont(DefaultFont)
	control.SetText("Form")
	return control
}

func (control *Form) SetLayout(mng LayoutManager) {
	control.layoutMng = mng
}

// UpdateLayout refresh layout.
func (control *Form) UpdateLayout() {
	if control.layoutMng != nil {
		control.layoutMng.Update()
	}
}

func (control *Form) NewMenu() *Menu {
	hMenu := w32.CreateMenu()
	if hMenu == 0 {
		panic("failed CreateMenu")
	}
	m := &Menu{hMenu: hMenu, hwnd: control.hwnd}
	if !w32.SetMenu(control.hwnd, hMenu) {
		panic("failed SetMenu")
	}
	return m
}

func (control *Form) DisableIcon() {
	windowInfo := getWindowInfo(control.hwnd)
	frameless := windowInfo.IsPopup()
	if frameless {
		return
	}
	exStyle := w32.GetWindowLong(control.hwnd, w32.GWL_EXSTYLE)
	w32.SetWindowLong(control.hwnd, w32.GWL_EXSTYLE, uint32(exStyle|w32.WS_EX_DLGMODALFRAME))
	w32.SetWindowPos(control.hwnd, 0, 0, 0, 0, 0,
		uint(
			w32.SWP_FRAMECHANGED|
				w32.SWP_NOMOVE|
				w32.SWP_NOSIZE|
				w32.SWP_NOZORDER),
	)
}

func (control *Form) Maximise() {
	w32.ShowWindow(control.hwnd, w32.SW_MAXIMIZE)
}

func (control *Form) Minimise() {
	w32.ShowWindow(control.hwnd, w32.SW_MINIMIZE)
}

func (control *Form) Restore() {
	w32.ShowWindow(control.hwnd, w32.SW_RESTORE)
}

// Public methods
func (control *Form) Center() {

	windowInfo := getWindowInfo(control.hwnd)
	frameless := windowInfo.IsPopup()

	info := getMonitorInfo(control.hwnd)
	workRect := info.RcWork
	screenMiddleW := workRect.Left + (workRect.Right-workRect.Left)/2
	screenMiddleH := workRect.Top + (workRect.Bottom-workRect.Top)/2
	var winRect *w32.RECT
	if !frameless {
		winRect = w32.GetWindowRect(control.hwnd)
	} else {
		winRect = w32.GetClientRect(control.hwnd)
	}
	winWidth := winRect.Right - winRect.Left
	winHeight := winRect.Bottom - winRect.Top
	windowX := screenMiddleW - (winWidth / 2)
	windowY := screenMiddleH - (winHeight / 2)
	w32.SetWindowPos(control.hwnd, w32.HWND_TOP, int(windowX), int(windowY), int(winWidth), int(winHeight), w32.SWP_NOSIZE)
}

func (control *Form) Fullscreen() {
	if control.isFullscreen {
		return
	}

	control.previousWindowStyle = uint32(w32.GetWindowLongPtr(control.hwnd, w32.GWL_STYLE))
	monitor := w32.MonitorFromWindow(control.hwnd, w32.MONITOR_DEFAULTTOPRIMARY)
	var monitorInfo w32.MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	if !w32.GetMonitorInfo(monitor, &monitorInfo) {
		return
	}
	if !w32.GetWindowPlacement(control.hwnd, &control.previousWindowPlacement) {
		return
	}
	w32.SetWindowLong(control.hwnd, w32.GWL_STYLE, control.previousWindowStyle & ^uint32(w32.WS_OVERLAPPEDWINDOW))
	w32.SetWindowPos(control.hwnd, w32.HWND_TOP,
		int(monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Top),
		int(monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left),
		int(monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top),
		w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
	control.isFullscreen = true
}

func (control *Form) UnFullscreen() {
	if !control.isFullscreen {
		return
	}
	w32.SetWindowLong(control.hwnd, w32.GWL_STYLE, control.previousWindowStyle)
	w32.SetWindowPlacement(control.hwnd, &control.previousWindowPlacement)
	w32.SetWindowPos(control.hwnd, 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOZORDER|w32.SWP_NOOWNERZORDER|w32.SWP_FRAMECHANGED)
	control.isFullscreen = false
}

// IconType: 1 - ICON_BIG; 0 - ICON_SMALL
func (control *Form) SetIcon(iconType int, icon *Icon) {
	if iconType > 1 {
		panic("IconType is invalid")
	}
	w32.SendMessage(control.hwnd, w32.WM_SETICON, uintptr(iconType), uintptr(icon.Handle()))
}

func (control *Form) EnableMaxButton(b bool) {
	ToggleStyle(control.hwnd, b, w32.WS_MAXIMIZEBOX)
}

func (control *Form) EnableMinButton(b bool) {
	ToggleStyle(control.hwnd, b, w32.WS_MINIMIZEBOX)
}

func (control *Form) EnableSizable(b bool) {
	ToggleStyle(control.hwnd, b, w32.WS_THICKFRAME)
}

func (control *Form) EnableDragMove(_ bool) {
	//control.isDragMove = b
}

func (control *Form) EnableTopMost(b bool) {
	tag := w32.HWND_NOTOPMOST
	if b {
		tag = w32.HWND_TOPMOST
	}
	w32.SetWindowPos(control.hwnd, tag, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE)
}

func (control *Form) WndProc(msg uint32, wparam, lparam uintptr) uintptr {

	switch msg {
	case w32.WM_COMMAND:
		if lparam == 0 && w32.HIWORD(uint32(wparam)) == 0 {
			// Menu support.
			actionID := uint16(w32.LOWORD(uint32(wparam)))
			if action, ok := actionsByID[actionID]; ok {
				action.onClick.Fire(NewEvent(control, nil))
			}
		}
	case w32.WM_KEYDOWN:
		// Accelerator support.
		key := Key(wparam)
		if uint32(lparam)>>30 == 0 {
			// Using TranslateAccelerators refused to work, so we handle them
			// ourselves, at least for now.
			shortcut := Shortcut{ModifiersDown(), key}
			if action, ok := shortcut2Action[shortcut]; ok {
				if action.Enabled() {
					action.onClick.Fire(NewEvent(control, nil))
				}
			}
		}

	case w32.WM_CLOSE:
		return 0
	case w32.WM_DESTROY:
		w32.PostQuitMessage(0)
		return 0

	case w32.WM_SIZE, w32.WM_PAINT:
		if control.layoutMng != nil {
			control.layoutMng.Update()
		}
	case w32.WM_GETMINMAXINFO:
		if control.minWidth != 0 || control.maxWidth != 0 || control.minHeight != 0 || control.maxHeight != 0 {
			dpix, dpiy := control.GetWindowDPI()

			DPIScaleX := dpix / 96.0
			DPIScaleY := dpiy / 96.0

			mmi := (*w32.MINMAXINFO)(unsafe.Pointer(lparam))
			if control.minWidth > 0 && control.minHeight > 0 {
				mmi.PtMinTrackSize.X = int32(control.minWidth * int(DPIScaleX))
				mmi.PtMinTrackSize.Y = int32(control.minHeight * int(DPIScaleY))
			}
			if control.maxWidth > 0 && control.maxHeight > 0 {
				mmi.PtMaxSize.X = int32(control.maxWidth * int(DPIScaleX))
				mmi.PtMaxSize.Y = int32(control.maxHeight * int(DPIScaleY))
				mmi.PtMaxTrackSize.X = int32(control.maxWidth * int(DPIScaleX))
				mmi.PtMaxTrackSize.Y = int32(control.maxHeight * int(DPIScaleY))
			}
			return 0
		}
	}

	return w32.DefWindowProc(control.hwnd, msg, wparam, lparam)
}

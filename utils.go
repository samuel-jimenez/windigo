/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/samuel-jimenez/windigo/w32"
)

func internalTrackMouseEvent(hwnd w32.HWND) {
	var tme w32.TRACKMOUSEEVENT
	tme.CbSize = uint32(unsafe.Sizeof(tme))
	tme.DwFlags = w32.TME_LEAVE
	tme.HwndTrack = hwnd
	tme.DwHoverTime = w32.HOVER_DEFAULT

	w32.TrackMouseEvent(&tme)
}

func ToggleStyle(hwnd w32.HWND, b bool, style int) {
	originalStyle := int(w32.GetWindowLongPtr(hwnd, w32.GWL_STYLE))
	if originalStyle != 0 {
		if b {
			originalStyle |= style
		} else {
			originalStyle ^= style
		}
		w32.SetWindowLongPtr(hwnd, w32.GWL_STYLE, uintptr(originalStyle))
	}
}

func ToggleExStyle(hwnd w32.HWND, b bool, style int) {
	originalStyle := int(w32.GetWindowLongPtr(hwnd, w32.GWL_EXSTYLE))
	if originalStyle != 0 {
		if b {
			originalStyle |= style
		} else {
			originalStyle ^= style
		}
		w32.SetWindowLongPtr(hwnd, w32.GWL_EXSTYLE, uintptr(originalStyle))
	}
}

func CreateWindow(className string, parent Controller, exStyle, style uint) w32.HWND {
	instance := GetAppInstance()
	var parentHwnd w32.HWND
	if parent != nil {
		parentHwnd = parent.Handle()
	}
	var hwnd w32.HWND
	hwnd = w32.CreateWindowEx(
		exStyle,
		syscall.StringToUTF16Ptr(className),
		nil,
		style,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		w32.CW_USEDEFAULT,
		parentHwnd,
		0,
		instance,
		nil)

	if hwnd == 0 {
		errStr := fmt.Sprintf("Error occurred in CreateWindow(%s, %v, %d, %d)", className, parent, exStyle, style)
		panic(errStr)
	}

	return hwnd
}

func RegisterClass(className string, wndproc uintptr) {
	instance := GetAppInstance()
	icon := w32.LoadIcon(instance, w32.MakeIntResource(w32.IDI_APPLICATION))

	var wc w32.WNDCLASSEX
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	wc.WndProc = wndproc
	wc.Instance = instance
	wc.Background = w32.COLOR_BTNFACE + 1
	wc.Icon = icon
	wc.Cursor = w32.LoadCursor(0, w32.MakeIntResource(w32.IDC_ARROW))
	wc.ClassName = syscall.StringToUTF16Ptr(className)
	wc.MenuName = nil
	wc.IconSm = icon

	if ret := w32.RegisterClassEx(&wc); ret == 0 {
		panic(syscall.GetLastError())
	}
}

func getMonitorInfo(hwnd w32.HWND) *w32.MONITORINFO {
	currentMonitor := w32.MonitorFromWindow(hwnd, w32.MONITOR_DEFAULTTONEAREST)
	var info w32.MONITORINFO
	info.CbSize = uint32(unsafe.Sizeof(info))
	w32.GetMonitorInfo(currentMonitor, &info)
	return &info
}
func getWindowInfo(hwnd w32.HWND) *w32.WINDOWINFO {
	var info w32.WINDOWINFO
	info.CbSize = uint32(unsafe.Sizeof(info))
	w32.GetWindowInfo(hwnd, &info)
	return &info
}

func RegClassOnlyOnce(className string) {
	isExists := false
	for _, class := range gRegisteredClasses {
		if class == className {
			isExists = true
			break
		}
	}

	if !isExists {
		RegisterClass(className, GeneralWndprocCallBack)
		gRegisteredClasses = append(gRegisteredClasses, className)
	}
}

func ScreenToClientRect(hwnd w32.HWND, rect *w32.RECT) *Rect {
	l, t, r, b := rect.Left, rect.Top, rect.Right, rect.Bottom
	l1, t1, _ := w32.ScreenToClient(hwnd, int(l), int(t))
	r1, b1, _ := w32.ScreenToClient(hwnd, int(r), int(b))
	return NewRect(l1, t1, r1, b1)
}

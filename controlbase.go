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

type ControlBase struct {
	hwnd                                               w32.HWND
	font                                               *Font
	color_fg, color_highlight, color_bg                Color
	draw_color_fg, draw_color_highlight, draw_color_bg bool
	brush_bg                                           *Brush
	parent                                             Controller
	contextMenu                                        *MenuItem

	isForm bool

	statusbar *StatusBar

	minWidth, minHeight int
	maxWidth, maxHeight int

	margin_top, margin_btm,
	margin_left, margin_right int

	border_pen  *Pen
	erasure_pen *Pen

	// General events
	onCreate EventManager
	onClose  EventManager

	// Focus events
	onKillFocus EventManager
	onSetFocus  EventManager

	// Drag and drop events
	onDropFiles EventManager

	// Mouse events
	onLBDown    EventManager
	onLBUp      EventManager
	onLBDbl     EventManager
	onMBDown    EventManager
	onMBUp      EventManager
	onRBDown    EventManager
	onRBUp      EventManager
	onRBDbl     EventManager
	onMouseMove EventManager

	// use MouseControl to capture onMouseHover and onMouseLeave events.
	onMouseHover EventManager
	onMouseLeave EventManager

	// Keyboard events
	onKeyUp EventManager

	// Paint events
	onPaint EventManager
	onSize  EventManager
}

func (control *ControlBase) ClipboardCopyText(text string) bool {
	if !w32.OpenClipboard(control.hwnd) {
		return false
	}

	w32.EmptyClipboard()
	if len(text) == 0 {
		w32.CloseClipboard()
		return true
	}

	//allocate and fill buffer
	bufferSize := uint32(len(text)) * 16 //sizeof(TCHAR))
	globalBuffer := w32.GlobalAlloc(w32.GMEM_MOVEABLE, bufferSize)
	pGlobalBuffer := w32.GlobalLock(globalBuffer)
	w32.MoveMemory(pGlobalBuffer, unsafe.Pointer(syscall.StringToUTF16Ptr(text)), bufferSize)
	w32.GlobalUnlock(globalBuffer)

	// Place the handle on the clipboard.
	w32.SetClipboardData(w32.CF_UNICODETEXT, globalBuffer)

	// Close the clipboard.
	w32.CloseClipboard()
	return true
}

// initControl is called by controls: edit, button, treeview, listview, and so on.
func (control *ControlBase) InitControl(className string, parent Controller, exstyle, style uint) {
	control.hwnd = CreateWindow(className, parent, exstyle, style)
	if control.hwnd == 0 {
		panic("cannot create window for " + className)
	}
	control.parent = parent
}

// InitWindow is called by custom window based controls such as split, panel, etc.
func (control *ControlBase) InitWindow(className string, parent Controller, exstyle, style uint) {
	RegClassOnlyOnce(className)
	control.hwnd = CreateWindow(className, parent, exstyle, style)
	if control.hwnd == 0 {
		panic("cannot create window for " + className)
	}
	control.parent = parent
}

// SetTheme for TreeView and ListView controls.
func (control *ControlBase) SetTheme(appName string) error {
	if hr := w32.SetWindowTheme(control.hwnd, syscall.StringToUTF16Ptr(appName), nil); w32.FAILED(hr) {
		return fmt.Errorf("SetWindowTheme %d", hr)
	}
	return nil
}

func (control *ControlBase) Handle() w32.HWND {
	return control.hwnd
}

func (control *ControlBase) SetHandle(hwnd w32.HWND) {
	control.hwnd = hwnd
}

func (control *ControlBase) GetWindowDPI() (w32.UINT, w32.UINT) {
	monitor := w32.MonitorFromWindow(control.hwnd, w32.MONITOR_DEFAULTTOPRIMARY)
	var dpiX, dpiY w32.UINT
	w32.GetDPIForMonitor(monitor, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
	return dpiX, dpiY
}

func (control *ControlBase) SetAndClearStyleBits(set, clear uint32) error {
	style := uint32(w32.GetWindowLong(control.hwnd, w32.GWL_STYLE))
	if style == 0 {
		return fmt.Errorf("GetWindowLong")
	}

	if newStyle := style&^clear | set; newStyle != style {
		if w32.SetWindowLong(control.hwnd, w32.GWL_STYLE, newStyle) == 0 {
			return fmt.Errorf("SetWindowLong")
		}
	}
	return nil
}

func (control *ControlBase) SetIsForm(isform bool) {
	control.isForm = isform
}

func (control *ControlBase) SetText(caption string) {
	w32.SetWindowText(control.hwnd, caption)
}

func (control *ControlBase) Text() string {
	return w32.GetWindowText(control.hwnd)
}

func (control *ControlBase) Close() {
	UnRegMsgHandler(control.hwnd)
	w32.DestroyWindow(control.hwnd)
}

func (control *ControlBase) SetTranslucentBackground() {
	var accent = w32.ACCENT_POLICY{
		AccentState: w32.ACCENT_ENABLE_BLURBEHIND,
	}
	var data w32.WINDOWCOMPOSITIONATTRIBDATA
	data.Attrib = w32.WCA_ACCENT_POLICY
	data.PvData = unsafe.Pointer(&accent)
	data.CbData = unsafe.Sizeof(accent)

	w32.SetWindowCompositionAttribute(control.hwnd, &data)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (control *ControlBase) clampSize(width, height int) (int, int) {
	if control.minWidth != 0 {
		width = max(width, control.minWidth)
	}
	if control.maxWidth != 0 {
		width = min(width, control.maxWidth)
	}
	if control.minHeight != 0 {
		height = max(height, control.minHeight)
	}
	if control.maxHeight != 0 {
		height = min(height, control.maxHeight)
	}
	return width, height
}

func (control *ControlBase) SetSize(width, height int) {
	x, y := control.Pos()
	width, height = control.clampSize(width, height)
	w32.MoveWindow(control.hwnd, x, y, width, height, true)
}

func (control *ControlBase) SetMinSize(width, height int) {
	control.minWidth = width
	control.minHeight = height

	// Ensure we set max if min > max
	if control.maxWidth > 0 {
		control.maxWidth = max(control.minWidth, control.maxWidth)
	}
	if control.maxHeight > 0 {
		control.maxHeight = max(control.minHeight, control.maxHeight)
	}

	x, y := control.Pos()
	currentWidth, currentHeight := control.Size()
	clampedWidth, clampedHeight := control.clampSize(currentWidth, currentHeight)
	if clampedWidth != currentWidth || clampedHeight != currentHeight {
		w32.MoveWindow(control.hwnd, x, y, clampedWidth, clampedHeight, true)
	}
}
func (control *ControlBase) SetMaxSize(width, height int) {
	control.maxWidth = width
	control.maxHeight = height

	// Ensure we set min if max > min
	if control.minWidth > 0 {
		control.minWidth = min(control.maxWidth, control.minWidth)
	}
	if control.maxHeight > 0 {
		control.minHeight = min(control.maxHeight, control.minHeight)
	}

	x, y := control.Pos()
	currentWidth, currentHeight := control.Size()
	clampedWidth, clampedHeight := control.clampSize(currentWidth, currentHeight)
	if clampedWidth != currentWidth || clampedHeight != currentHeight {
		w32.MoveWindow(control.hwnd, x, y, clampedWidth, clampedHeight, true)
	}
}

func (control *ControlBase) Size() (width, height int) {
	rect := w32.GetWindowRect(control.hwnd)
	width = int(rect.Right - rect.Left)
	height = int(rect.Bottom - rect.Top)
	return
}

func (control *ControlBase) Width() int {
	rect := w32.GetWindowRect(control.hwnd)
	return int(rect.Right - rect.Left)
}

func (control *ControlBase) Height() int {
	rect := w32.GetWindowRect(control.hwnd)
	return int(rect.Bottom - rect.Top)
}

func (control *ControlBase) SetPos(x, y int) {
	w32.SetWindowPos(control.hwnd, 0, x, y, 0, 0, w32.SWP_NOSIZE|w32.SWP_NOZORDER)
}

func (control *ControlBase) SetPosAfter(x, y int, after Dockable) {
	w32.SetWindowPos(control.hwnd, after.Handle(), x, y, 0, 0, w32.SWP_NOSIZE)
}

func (control *ControlBase) SetZAfter(after Dockable) {
	w32.SetWindowPos(control.hwnd, after.Handle(), 0, 0, 0, 0, w32.SWP_NOSIZE|w32.SWP_NOMOVE)
}

func (control *ControlBase) SetZTop() {
	w32.SetWindowPos(control.hwnd, w32.HWND_TOP, 0, 0, 0, 0, w32.SWP_NOSIZE|w32.SWP_NOMOVE)
}

func (control *ControlBase) SetZBottom() {
	w32.SetWindowPos(control.hwnd, w32.HWND_BOTTOM, 0, 0, 0, 0, w32.SWP_NOSIZE|w32.SWP_NOMOVE)
}

func (control *ControlBase) Pos() (x, y int) {
	rect := w32.GetWindowRect(control.hwnd)
	x = int(rect.Left)
	y = int(rect.Top)
	if !control.isForm && control.parent != nil {
		x, y, _ = w32.ScreenToClient(control.parent.Handle(), x, y)
	}
	return
}

// Marginal
func (control *ControlBase) MarginTop() int {
	return control.margin_top
}

func (control *ControlBase) MarginBtm() int {
	return control.margin_btm
}
func (control *ControlBase) MarginLeft() int {
	return control.margin_left
}
func (control *ControlBase) MarginRight() int {
	return control.margin_right
}

func (control *ControlBase) SetMarginsAll(margin int) {
	control.margin_top = margin
	control.margin_right = margin
	control.margin_btm = margin
	control.margin_left = margin
}

func (control *ControlBase) SetMarginsHV(margin_horizontal, margin_vertical int) {
	control.margin_top = margin_vertical
	control.margin_right = margin_horizontal
	control.margin_btm = margin_vertical
	control.margin_left = margin_horizontal
}

func (control *ControlBase) SetMargins(margin_left, margin_top, margin_right, margin_btm int) {
	control.margin_top = margin_top
	control.margin_right = margin_right
	control.margin_btm = margin_btm
	control.margin_left = margin_left
}

func (control *ControlBase) SetMarginTop(margin int) {
	control.margin_top = margin
}

func (control *ControlBase) SetMarginBtm(margin int) {
	control.margin_btm = margin
}

func (control *ControlBase) SetMarginLeft(margin int) {
	control.margin_left = margin
}

func (control *ControlBase) SetMarginRight(margin int) {
	control.margin_right = margin
}

func (control *ControlBase) Visible() bool {
	return w32.IsWindowVisible(control.hwnd)
}

func (control *ControlBase) ToggleVisible() bool {
	visible := w32.IsWindowVisible(control.hwnd)
	if visible {
		control.Hide()
	} else {
		control.Show()
	}
	return !visible
}

func (control *ControlBase) Border() *Pen {
	return control.border_pen
}

func (control *ControlBase) SetBorder(pen *Pen) {
	if control.brush_bg == nil {
		control.brush_bg = DefaultBackgroundBrush
	}
	if control.border_pen != nil {
		//erase old pen marks
		control.erasure_pen = NewPen(w32.PS_GEOMETRIC|w32.PS_SOLID|w32.PS_ENDCAP_SQUARE, control.border_pen.Width(), control.brush_bg)
	}
	control.border_pen = pen
	control.Invalidate(true)
}

func (control *ControlBase) drawBorder(canvas *Canvas) {
	if control.border_pen == nil {
		return
	}
	margin_0 := int32(control.border_pen.width) / 2
	margin_1 := int32(control.border_pen.width+1) / 2

	r := control.Bounds()

	if control.erasure_pen != nil {
		canvas.DrawFillRect(r, control.erasure_pen, control.brush_bg)
		control.erasure_pen.Dispose()
		control.erasure_pen = nil
	}
	r.rect.Top += margin_0
	r.rect.Left += margin_0
	r.rect.Bottom -= margin_1
	r.rect.Right -= margin_1
	canvas.DrawFillRect(r, control.border_pen, control.brush_bg)

}

func (control *ControlBase) ContextMenu() *MenuItem {
	return control.contextMenu
}

func (control *ControlBase) SetContextMenu(menu *MenuItem) {
	control.contextMenu = menu
}

func (control *ControlBase) Bounds() *Rect {
	rect := w32.GetWindowRect(control.hwnd)
	if control.isForm {
		return &Rect{*rect}
	}

	return ScreenToClientRect(control.hwnd, rect)
}

func (control *ControlBase) WindowBounds() *Rect {
	return ScreenToClientRect(control.hwnd, w32.GetWindowRect(control.hwnd))
}

func (control *ControlBase) ClientRect() *Rect {
	rect := w32.GetClientRect(control.hwnd)
	return ScreenToClientRect(control.hwnd, rect)
}
func (control *ControlBase) ClientWidth() int {
	rect := w32.GetClientRect(control.hwnd)
	return int(rect.Right - rect.Left)
}

func (control *ControlBase) ClientHeight() int {
	rect := w32.GetClientRect(control.hwnd)
	return int(rect.Bottom - rect.Top)
}

func (control *ControlBase) Show() {
	w32.ShowWindow(control.hwnd, w32.SW_SHOWDEFAULT)
}

func (control *ControlBase) Hide() {
	w32.ShowWindow(control.hwnd, w32.SW_HIDE)
}

func (control *ControlBase) Enabled() bool {
	return w32.IsWindowEnabled(control.hwnd)
}

func (control *ControlBase) SetEnabled(b bool) {
	w32.EnableWindow(control.hwnd, b)
}

func (control *ControlBase) SetFocus() {
	w32.SetFocus(control.hwnd)
}

func (control *ControlBase) Invalidate(erase bool) {
	// pRect := w32.GetClientRect(control.hwnd)
	// if control.isForm {
	// 	w32.InvalidateRect(control.hwnd, pRect, erase)
	// } else {
	// 	rc := ScreenToClientRect(control.hwnd, pRect)
	// 	w32.InvalidateRect(control.hwnd, rc.GetW32Rect(), erase)
	// }
	w32.InvalidateRect(control.hwnd, nil, erase)
}

func (control *ControlBase) Parent() Controller {
	return control.parent
}

func (control *ControlBase) SetParent(parent Controller) {
	control.parent = parent
}

func (control *ControlBase) StatusBar() *StatusBar {
	return control.statusbar
}

func (control *ControlBase) SetStatusBar(statusbar *StatusBar) {
	control.statusbar = statusbar
}

func (control *ControlBase) RefreshStatusBar() {
	if control.statusbar != nil {
		control.statusbar.SetSize(0, 0)
	}
}

func (control *ControlBase) Font() *Font {
	return control.font
}

func (control *ControlBase) SetFont(font *Font) {
	w32.SendMessage(control.hwnd, w32.WM_SETFONT, uintptr(font.hfont), 1)
	control.font = font
}

func (control *ControlBase) SetFGColor(color Color) {
	control.color_fg = color
	control.draw_color_fg = true
}

func (control *ControlBase) ClearFGColor() {
	control.draw_color_fg = false
}

func (control *ControlBase) HasFGColor() bool {
	return control.draw_color_fg
}

func (control *ControlBase) FGColor() Color {
	return control.color_fg
}

func (control *ControlBase) SetHighlightColor(color Color) {
	control.color_highlight = color
	control.draw_color_highlight = true
}

func (control *ControlBase) ClearHighlightColor() {
	control.draw_color_highlight = false
}

func (control *ControlBase) HasHighlightColor() bool {
	return control.draw_color_highlight
}

func (control *ControlBase) HighlightColor() Color {
	return control.color_highlight
}

func (control *ControlBase) SetBGColor(color Color) {
	control.color_bg = color
	control.draw_color_bg = true
	if control.brush_bg != nil && control.brush_bg != DefaultBackgroundBrush {
		control.brush_bg.Dispose()
	}
	control.brush_bg = NewSolidColorBrush(color)
}

func (control *ControlBase) ClearBGColor() {
	control.draw_color_bg = false
}

func (control *ControlBase) HasBGColor() bool {
	return control.draw_color_bg
}

func (control *ControlBase) BGColor() Color {
	return control.color_bg
}

func (control *ControlBase) BGBrush() *Brush {
	return control.brush_bg
}

func (control *ControlBase) EnableDragAcceptFiles(b bool) {
	w32.DragAcceptFiles(control.hwnd, b)
}

func (control *ControlBase) InvokeRequired() bool {
	if control.hwnd == 0 {
		return false
	}

	windowThreadId, _ := w32.GetWindowThreadProcessId(control.hwnd)
	currentThreadId := w32.GetCurrentThread()

	return windowThreadId != currentThreadId
}

func (control *ControlBase) PreTranslateMessage(msg *w32.MSG) bool {
	if msg.Message == w32.WM_GETDLGCODE {
		println("pretranslate, WM_GETDLGCODE")
	}
	return false
}

// Events
func (control *ControlBase) OnCreate() *EventManager {
	return &control.onCreate
}

func (control *ControlBase) OnClose() *EventManager {
	return &control.onClose
}

func (control *ControlBase) OnKillFocus() *EventManager {
	return &control.onKillFocus
}

func (control *ControlBase) OnSetFocus() *EventManager {
	return &control.onSetFocus
}

func (control *ControlBase) OnDropFiles() *EventManager {
	return &control.onDropFiles
}

func (control *ControlBase) OnLBDown() *EventManager {
	return &control.onLBDown
}

func (control *ControlBase) OnLBUp() *EventManager {
	return &control.onLBUp
}

func (control *ControlBase) OnLBDbl() *EventManager {
	return &control.onLBDbl
}

func (control *ControlBase) OnMBDown() *EventManager {
	return &control.onMBDown
}

func (control *ControlBase) OnMBUp() *EventManager {
	return &control.onMBUp
}

func (control *ControlBase) OnRBDown() *EventManager {
	return &control.onRBDown
}

func (control *ControlBase) OnRBUp() *EventManager {
	return &control.onRBUp
}

func (control *ControlBase) OnRBDbl() *EventManager {
	return &control.onRBDbl
}

func (control *ControlBase) OnMouseMove() *EventManager {
	return &control.onMouseMove
}

func (control *ControlBase) OnMouseHover() *EventManager {
	return &control.onMouseHover
}

func (control *ControlBase) OnMouseLeave() *EventManager {
	return &control.onMouseLeave
}

func (control *ControlBase) OnPaint() *EventManager {
	return &control.onPaint
}

func (control *ControlBase) OnSize() *EventManager {
	return &control.onSize
}

func (control *ControlBase) OnKeyUp() *EventManager {
	return &control.onKeyUp
}

// RunMainLoop processes messages in main application loop.
func (control *ControlBase) RunMainLoop() int {
	var m w32.MSG

	for w32.GetMessage(&m, 0, 0, 0) != 0 {
		if !PreTranslateMessage(&m) {
			if !w32.IsDialogMessage(control.hwnd, &m) {
				w32.TranslateMessage(&m)
				w32.DispatchMessage(&m)
			}
		}
	}

	w32.GdiplusShutdown()
	return int(m.WParam)
}

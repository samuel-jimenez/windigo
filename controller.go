/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"github.com/samuel-jimenez/windigo/w32"
)

// BaseController needs to be mutually exclusive with ComponentFrame
type BaseController interface {
	Text() string

	Enabled() bool
	SetFocus()

	Invalidate(erase bool)
	Parent() Controller

	Size() (w, h int)
	Bounds() *Rect
	WindowBounds() *Rect
	ClientRect() *Rect

	SetText(s string)
	SetEnabled(b bool)
	SetPosAfter(x, y int, after Controller)
	EnableDragAcceptFiles(b bool)

	ContextMenu() *MenuItem
	SetContextMenu(menu *MenuItem)

	StatusBar() *StatusBar
	SetStatusBar(statusbar *StatusBar)
	RefreshStatusBar()

	Font() *Font
	SetFont(font *Font)

	SetFGColor(color Color)
	ClearFGColor()
	HasFGColor() bool
	FGColor() Color

	SetHighlightColor(color Color)
	ClearHighlightColor()
	HasHighlightColor() bool
	HighlightColor() Color

	SetBGColor(color Color)
	ClearBGColor()
	HasBGColor() bool
	BGColor() Color
	BGBrush() *Brush

	InvokeRequired() bool
	PreTranslateMessage(msg *w32.MSG) bool
	WndProc(msg uint32, wparam, lparam uintptr) uintptr

	//General events
	OnCreate() *EventManager
	OnClose() *EventManager

	// Focus events
	OnKillFocus() *EventManager
	OnSetFocus() *EventManager

	//Drag and drop events
	OnDropFiles() *EventManager

	//Mouse events
	OnLBDown() *EventManager
	// OnLBUp() *EventManager
	OnLBDbl() *EventManager
	OnMBDown() *EventManager
	OnMBUp() *EventManager
	OnRBDown() *EventManager
	OnRBUp() *EventManager
	OnRBDbl() *EventManager
	// OnMouseMove() *EventManager

	// OnMouseLeave and OnMouseHover does not fire unless control called internalTrackMouseEvent.
	// Use MouseControl for a how to example.
	OnMouseHover() *EventManager
	OnMouseLeave() *EventManager

	//Keyboard events
	OnKeyUp() *EventManager

	//Paint events
	OnPaint() *EventManager
	OnSize() *EventManager
}

type Controller interface {
	BaseController
	Dockable
	Bordered
}

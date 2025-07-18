/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"github.com/samuel-jimenez/windigo/w32"
)

type RawMsg struct {
	Hwnd           w32.HWND
	Msg            uint32
	WParam, LParam uintptr
}

type MouseEventData struct {
	X, Y   int
	Button int
	Wheel  int
}

type DropFilesEventData struct {
	X, Y  int
	Files []string
}

type PaintEventData struct {
	Canvas *Canvas
}

type LabelEditEventData struct {
	Item ListItem
	Text string
	//PszText *uint16
}

/*type LVDBLClickEventData struct {
	NmItem *w32.NMITEMACTIVATE
}*/

type KeyUpEventData struct {
	VKey, Code int
}

type SizeEventData struct {
	Type uint
	X, Y int
}

type ListViewEvent struct {
	Row, Column int
}

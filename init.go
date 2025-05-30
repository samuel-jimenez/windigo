/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

import (
	"github.com/samuel-jimenez/windigo/w32"
)

func init() {
	gControllerRegistry = make(map[w32.HWND]Controller)
	gRegisteredClasses = make([]string, 0)

	var si w32.GdiplusStartupInput
	si.GdiplusVersion = 1
	w32.GdiplusStartup(&si, nil)
}

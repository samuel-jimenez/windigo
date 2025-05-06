package main

import (
	"github.com/samuel-jimenez/windigo"
	"github.com/samuel-jimenez/windigo/w32"
)

// TopForm displayed as topmost window until closed.
// By itself this is not very useful since Form has function EnableTopMost() making form topmost.
// This is just an example showing how custom window Form can be implemented inside your package.
type TopForm struct {
	windigo.Form

	onLoad windigo.EventManager
}

func NewTopForm(parent windigo.Controller) *TopForm {
	dlg := new(TopForm)
	dlg.SetIsForm(true)

	windigo.RegClassOnlyOnce("my_TopForm")
	dlg.SetHandle(windigo.CreateWindow("my_TopForm", parent, w32.WS_EX_DLGMODALFRAME|w32.WS_EX_TOPMOST,
		w32.WS_VISIBLE|w32.WS_SYSMENU|w32.WS_CAPTION))
	dlg.SetParent(parent)

	// dlg might fail if icon resource is not embedded in the binary
	if ico, err := windigo.NewIconFromResource(windigo.GetAppInstance(), uint16(windigo.AppIconID)); err == nil {
		dlg.SetIcon(0, ico)
	}

	// Dlg forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(dlg.Handle(), w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)
	windigo.RegMsgHandler(dlg)

	dlg.SetFont(windigo.DefaultFont)
	dlg.SetText("Form")
	return dlg
}

// Events
func (dlg *TopForm) OnLoad() *windigo.EventManager {
	return &dlg.onLoad
}

func (dlg *TopForm) Show() {
	dlg.onLoad.Fire(windigo.NewEvent(dlg, nil))
	dlg.Form.Show()
}

func (dlg *TopForm) WndProc(msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case w32.WM_CLOSE:
		dlg.Close()
	case w32.WM_DESTROY:
		if dlg.Parent() == nil {
			windigo.Exit()
		}
	}
	return w32.DefWindowProc(dlg.Handle(), msg, wparam, lparam)
}

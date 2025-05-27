package main

import (
	"github.com/samuel-jimenez/windigo"
)

func main() {
	mainWindow := NewTopForm(nil) // Our TopForm control gets created here.
	mainWindow.SetSize(400, 300)
	mainWindow.SetText("Hello World Demo")

	edt := windigo.NewEdit(mainWindow)
	edt.SetPos(10, 20)
	// Most Controls have default size unless SetSize is called.
	edt.SetText("edit text")

	btn := windigo.NewPushButton(mainWindow)
	btn.SetText("Show or Hide")
	btn.SetPos(40, 50)
	btn.SetSize(100, 40)
	btn.OnClick().Bind(func(e *windigo.Event) {
		if edt.Visible() {
			edt.Hide()
		} else {
			edt.Show()
		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	mainWindow.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *windigo.Event) {
	windigo.Exit()
}

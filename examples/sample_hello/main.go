package main

import (
	"github.com/samuel-jimenez/windigo"
)

func main() {
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetSize(800, 600)
	mainWindow.SetText("Hello World Demo")

	dock := windigo.NewSimpleDock(mainWindow)

	group_panel := windigo.NewAutoPanel(mainWindow)
	//When docking, controls are expanded to fill available space
	//e.g., since this will dock to Top, its width will be ignored
	group_panel.SetSize(400, 105)
	group_panel.SetMarginsAll(5)

	edit_field := windigo.NewLabeledEdit(group_panel, 110, 200, 22, "Edit Control")
	edit_field.SetMarginsAll(15)
	optional_field := windigo.NewLabeledEdit(group_panel, 110, 200, 22, "Optional Control")
	optional_field.SetMarginsAll(25)

	group_panel.Dock(edit_field, windigo.Top)
	group_panel.Dock(optional_field, windigo.Top)

	button := windigo.NewPushButton(mainWindow)
	button.SetText("Show or Hide")
	button.SetSize(100, 40)
	button.OnClick().Bind(func(e *windigo.Event) {
		if optional_field.Visible() {
			optional_field.Hide()
		} else {
			optional_field.Show()
		}
	})
	button.SetMarginsAll(5)
	dock.Dock(button, windigo.Left)

	dock.Dock(group_panel, windigo.Top)

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	windigo.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *windigo.Event) {
	windigo.Exit()
}

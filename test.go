package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	app2 := app.New()

	w := app2.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app2.Quit()
		}),
	))

	//a := app2.NewWindow("Title")


	w.ShowAndRun()
}

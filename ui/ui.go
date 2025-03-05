package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SendToApplication(){

}

func StartUi(){
	UI := app.New()
	Window := UI.NewWindow("Calculation Service")
	options := []string{"Calculate Expression", "Get Expressions List", "Get Expression By Id"}

	selectWidget := widget.NewSelect(options, func(selected string) {
		log.Println("User choosed " + selected)
	})

	btn := widget.NewButton("Send Expression", SendToApplication)
	Window.SetContent(container.NewVBox(selectWidget, btn))
	Window.SetContent(container.NewVBox(btn))
	Window.Resize(fyne.NewSize(300, 200))
}

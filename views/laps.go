package views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"lapcounter/models"
	"log"
	"time"
)

var (
	clock *canvas.Text
	lapCounter *canvas.Text
	distance *canvas.Text
	doRunClock bool = true
)

func BuildView() *fyne.Container {
	clock = canvas.NewText("00:00 MM", color.White)
	lapCounter = canvas.NewText("0", color.White)
	distance = canvas.NewText(fmt.Sprintf("%0.2f", 0.00), color.White)

	go updateValues()

	button := widget.NewButton("Reset", func() {
		log.Println("Reset")
	})

	clockBar := container.NewHBox(
			layout.NewSpacer(),
			clock,
			layout.NewSpacer(),
		)

	buttonBar := container.NewHBox(
		layout.NewSpacer(),
		button,
		layout.NewSpacer(),
		)

	lapCard := widget.NewCard(
		"Laps",
		"",
		lapCounter,
		)

	distanceCard := widget.NewCard(
		"Distance",
		"in miles",
		distance,
		)

	cardGrid := container.NewGridWithColumns(2, lapCard, distanceCard)

	border := container.NewBorder(
		clockBar,
			buttonBar,
			nil,
			nil,
			cardGrid,
		)

	go runClock()
	return border
}
func updateValues() {
	for update := range models.UpdateCh {
		models.DefaultLapModel().LapCount += update
		log.Printf("%d", models.DefaultLapModel().LapCount)

		lapCounter.Text = fmt.Sprintf("%d", models.DefaultLapModel().LapCount)
		lapCounter.Refresh()

		distance.Text = fmt.Sprintf("%0.2f", models.DefaultLapModel().Distance)
		distance.Refresh()
	}
}

func StopClock() {
	doRunClock = false
}

func runClock() {
	tickChannel := time.Tick(time.Second * 1)
	for _ = range tickChannel {
		if !doRunClock {
			break
		}

		clock.Text = time.Now().Format("03:04:05 PM")
		clock.Refresh()
	}
}

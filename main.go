package main

import (
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"lapcounter/models"
	"lapcounter/views"
	"log"
	"strings"
)

var (
	adapter string
	firmataAdapter *firmata.Adaptor
	sensor *gpio.PIRMotionDriver
)

func initPir() {
	firmataAdapter = firmata.NewAdaptor(adapter)
	if strings.Compare(adapter, "") == 0 {
		log.Println("No adapter was set.")
		return
	}

	firmataAdapter.Connect()
	sensor = gpio.NewPIRMotionDriver(firmataAdapter, "2")

	sensor.On(gpio.MotionDetected, func(data interface{}) {
		models.DefaultLapModel().LapCount++
		log.Printf("%s, lapcount: %d", gpio.MotionDetected, models.DefaultLapModel().LapCount)
	})
	sensor.On(gpio.MotionStopped, func(data interface{}) {
		log.Println(gpio.MotionStopped)
	})

	err := sensor.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PIR connected!")
}

func main() {
	flag.StringVar(&adapter, "a", "", "The usb adapter string.")
	flag.Parse()

	lapApp := app.New()
	window := lapApp.NewWindow("Lap Counter")

	window.SetContent(views.BuildView())
	window.Resize(fyne.Size{Height: 300, Width: 300})

	go initPir()

	defer func() {
		views.StopClock()
		sensor.Halt()
		firmataAdapter.Disconnect()
	}()
	window.ShowAndRun()
}

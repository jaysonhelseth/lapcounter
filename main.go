package main

import (
	"flag"
	"fyne.io/fyne/v2/app"
	"gobot.io/x/gobot"
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
	robot *gobot.Robot
)

func initPir() {
	if strings.Compare(adapter, "") == 0 {
		log.Println("No adapter was set.")
		return
	}

	firmataAdapter = firmata.NewAdaptor(adapter)
	sensor = gpio.NewPIRMotionDriver(firmataAdapter, "2")

	work := func() {
		sensor.On(gpio.MotionDetected, func(data interface{}) {
			log.Println(gpio.MotionDetected)
			models.UpdateCh <- 1
		})
		sensor.On(gpio.MotionStopped, func(data interface{}) {
			log.Println(gpio.MotionStopped)
		})
	}

	robot = gobot.NewRobot("PIR",
		[]gobot.Connection{firmataAdapter},
		[]gobot.Device{sensor},
		work,
	)

	log.Println("PIR connected!")
	robot.Start()
}

func main() {
	flag.StringVar(&adapter, "a", "", "The usb adapter string.")
	flag.Parse()

	lapApp := app.New()
	window := lapApp.NewWindow("Lap Counter")

	lapApp.Settings().SetTheme(&myTheme{})
	window.SetContent(views.BuildView())
	window.SetFullScreen(true)

	go initPir()

	defer func() {
		close(models.UpdateCh)
		views.StopClock()
		robot.Stop()
		firmataAdapter.Disconnect()
	}()
	window.ShowAndRun()
}

package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
	"gopkg.in/yaml.v2"
)

const bfSizeChLED = 4096

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("A configuration YAML file path should specified as a command argument.\n")
		os.Exit(1)
	}

	path := os.Args[1]

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	var cnf Config

	if err := yaml.Unmarshal(b, &cnf); err != nil {
		log.Panic(err)
	}

	if err := rpio.Open(); err != nil {
		log.Panic(err)
	}

	defer rpio.Close()

	chLED, chLEDDone := ledOn(cnf)

	defer func() {
		close(chLED)
		<-chLEDDone
	}()

	pin := rpio.Pin(cnf.MotionPinNo)
	pin.Input()

	log.Println("Ready")

	for {
		if pin.Read() == rpio.High {
			log.Println("Motion detected")
			chLED <- struct{}{}
		}

		time.Sleep(1 * time.Second)
	}
}

func ledOn(cnf Config) (chan<- struct{}, <-chan struct{}) {
	chLED := make(chan struct{}, bfSizeChLED)
	chLEDDone := make(chan struct{})

	go func() {
		pin := rpio.Pin(cnf.LEDPinNo)
		pin.Output()

		for _ = range chLED {
			log.Println("LED on")

			pin.High()

			time.Sleep(1 * time.Second)

			pin.Low()
		}

		chLEDDone <- struct{}{}
	}()

	return chLED, chLEDDone
}

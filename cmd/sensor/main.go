package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/stianeikeland/go-rpio"
	"gopkg.in/yaml.v2"
)

const bfSizeChLED = 4096
const bfSizeChSend = 4096

func main() {
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, os.Interrupt, os.Kill)

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

	chSend, chSendDone := send(cnf)

	defer func() {
		close(chSend)
		<-chSendDone
	}()

	pin := rpio.Pin(cnf.MotionPinNo)
	pin.Input()

	log.Println("Ready")

	ticker := time.NewTicker(1 * time.Second)

	defer ticker.Stop()

	for {
		if pin.Read() == rpio.High {
			t := time.Now()
			log.Println("Motion detected")
			chLED <- struct{}{}
			chSend <- t
		}

		select {
		case <-ticker.C:
		case sig := <-chSig:
			log.Println("Got signal:", sig)
			return
		}
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

func send(cnf Config) (chan<- time.Time, <-chan struct{}) {
	chSend := make(chan time.Time, bfSizeChSend)
	chSendDone := make(chan struct{})

	go func() {
		for t := range chSend {
			log.Printf("Send motion data [time: %+v]\n", t)
		}
	}()

	return chSend, chSendDone
}

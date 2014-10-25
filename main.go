package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
	"gopkg.in/yaml.v2"
)

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

	var conf Config

	if err := yaml.Unmarshal(b, &conf); err != nil {
		log.Panic(err)
	}

	log.Printf("%s\n", string(b))
	log.Printf("%+v\n", conf)

	if err := rpio.Open(); err != nil {
		log.Panic(err)
	}

	defer rpio.Close()

	pin := rpio.Pin(conf.MotionPinNo)
	pin.Input()

	log.Println("Ready")

	for {
		if pin.Read() == rpio.High {
			log.Println("Motion detected")
		}

		time.Sleep(1 * time.Second)
	}
}

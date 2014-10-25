package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v2"
)

func topIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello PIR HTTP Server")
}

func motionsCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusCreated)
}

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

	rtr := httprouter.New()
	rtr.GET("/", topIndex)
	rtr.POST("/motions", motionsCreate)

	log.Panic(http.ListenAndServe(":"+cnf.Port, rtr))
}

package main

import (
	"log"
	"net/http"

	"github.com/pmalek/proton_task/handler"
	"github.com/pmalek/proton_task/update"
)

const (
	Version = 2
)

func main() {
	updateProvider, err := update.NewFileSystemProvider(".")
	if err != nil {
		log.Fatal(err)
	}
	h, err := handler.New(Version, updateProvider)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", h.Index)
	http.HandleFunc("/check", h.Check)
	http.HandleFunc("/install", h.Install)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package terjang

import (
	"embed"
	"fmt"
	"log"
	"net/http"
)

//go:embed web/dist/*
var f embed.FS
var contents, err = f.ReadFile("readme.txt")

func static(w http.ResponseWriter, r *http.Request) {

	contents, err := f.ReadFile("readme.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", contents)
}

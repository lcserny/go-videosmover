package main

import (
	"flag"
	"fmt"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"log"
	"os"
)

var screenWidth = flag.Int("screenWidth", 0, "screen width")
var screenHeight = flag.Int("screenHeight", 0, "screen height")

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: Please provide `screenWidth` and `screenHeight` flags")
		return
	}

	flag.Parse()

	windowWidth, windowHeight := 400, 400
	top, left := (*screenHeight/2)-(windowHeight/2), (*screenWidth/2)-(windowWidth/2)

	rect := sciter.NewRect(top, left, windowWidth, windowHeight)
	w, err := window.New(sciter.DefaultWindowCreateFlag, rect)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("handle: %v", w.Handle)
	w.LoadFile("static/html/simple.html")
	w.SetTitle("Example")
	w.Show()
	w.Run()
}

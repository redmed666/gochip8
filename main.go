package main

import (
	"flag"
	"log"

	"github.com/redmed666/gochip8/chip8"
)

var (
	gamePath string
)

func init() {
	flag.StringVar(&gamePath, "path", "", "Path to the game")
}

func main() {
	flag.Parse()
	if gamePath == "" {
		log.Fatal("You must set the path to the game you want to load")
	}
	/*
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		go func() {
			log.Fatal(http.Serve(ln, nil))
		}()
		webview.Open("Test motherfucker", "https://google.com", 800, 600, true)
	*/
	chip8 := chip8.Chip8{}
	chip8.Initialize()
	chip8.LoadGame(gamePath)

	for {
		chip8.EmulateCycle()
	}
}

package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/rthornton128/goncurses"
)

func playAudio(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Error opening audio file: %v", err)
		return
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Printf("Error creating MP3 decoder: %v", err)
		return
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		log.Printf("Error creating audio context: %v", err)
		return
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := p.Write([]byte{0}); err != nil {
		log.Printf("Error playing audio: %v", err)
		return
	}

	_, err = f.Seek(0, 0) // Rewind the audio file
	if err != nil {
		log.Printf("Error rewinding audio: %v", err)
		return
	}

	if _, err := io.Copy(p, d); err != nil {
		log.Printf("Error playing audio: %v", err)
		return
	}
}

func main() {
	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer goncurses.End()

	goncurses.Cursor(0) // Hide the cursor

	// Set up colors (if supported)
	if goncurses.HasColors() {
		goncurses.StartColor()
		goncurses.InitPair(1, goncurses.C_CYAN, goncurses.C_BLACK)
	}

	stdscr.Clear()
	stdscr.Print("Welcome to GoNCurses!\n")
	stdscr.Refresh()

	// Create a new window
	win, err := goncurses.NewWindow(10, 40, 2, 2)
	if err != nil {
		log.Fatal(err)
	}
	defer win.Delete()

	win.Box(0, 0)
	win.MovePrint(1, 1, "This is a GoNCurses window")
	win.Refresh()

	// Start playing audio in a Goroutine
	go playAudio("../../audio/Enter.mp3")

	// Wait for a few seconds to play audio before accepting user input
	time.Sleep(5 * time.Second)

	// Wait for a keypress
	stdscr.Print("Press any key to exit...")
	stdscr.Refresh()
	stdscr.GetChar()
}

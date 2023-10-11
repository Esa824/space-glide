package main

import (
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/oto"
	"github.com/rthornton128/goncurses"

	"github.com/hajimehoshi/go-mp3"
)

func playAudio(file string) error {
    f, err := os.Open(file)
    if err != nil {
        log.Printf("Error opening audio file: %v", err)
        return err
    }
    defer f.Close()

    d, err := mp3.NewDecoder(f)
    if err != nil {
        log.Printf("Error creating MP3 decoder: %v", err)
        return err
    }

    c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
    if err != nil {
        log.Printf("Error creating audio context: %v", err)
        return err
    }
    defer c.Close()

    p := c.NewPlayer()
    defer p.Close()

    if _, err := io.Copy(p, d); err != nil {
        log.Printf("Error playing audio: %v", err)
        return err
    }

    return nil
}

func main() {
	goncurses.Init()
	for {
		if err := playAudio("../../audio/Enter.mp3"); err != nil {
			log.Fatal(err)
		}
	}
}



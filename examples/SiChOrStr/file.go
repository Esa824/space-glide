
package main

import (
    "log"
    gc "github.com/rthornton128/goncurses"
)

func main() {
    stdscr, err := gc.Init()
    if err != nil {
        log.Fatal(err)
    }
    defer gc.End()

    // Simulate pressing the 'A' key
    stdscr.Print("Hello I am Simulating a string")

    stdscr.Refresh()

    stdscr.GetChar() // Wait for user input to see the effect
}

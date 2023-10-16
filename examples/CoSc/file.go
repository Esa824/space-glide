package main

import (
	"encoding/json"
	gc "github.com/rthornton128/goncurses"
	"io"
	"log"
	"os"
)

// Define a struct to represent the control settings
type Controls struct {
	Up    string `json:"up"`
	Down  string `json:"down"`
	Left  string `json:"left"`
	Right string `json:"right"`
	Shoot string `json:"shoot"`
}

type Settings struct {
	Controls Controls `json:"controls"`
}

func main() {
	// Initialize gc
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal("Init:", err)
	}
	defer gc.End()

	gc.StartColor()
	gc.Cursor(0)
	gc.Echo(false)
	stdscr.Keypad(true)
	stdscr.Timeout(0)

	// Open the JSON file for reading
	file, err := os.Open("settings.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the JSON data from the file
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Create a variable to hold the control settings
	var settings Settings

	// Unmarshal the JSON data into the settings variable
	if err := json.Unmarshal(data, &settings); err != nil {
		log.Fatal(err)
	}

	// Display control settings
	stdscr.MovePrint(0, 10, "Up: "+settings.Controls.Up)
	stdscr.MovePrint(1, 10, "Down: "+settings.Controls.Down)
	stdscr.MovePrint(2, 10, "Left: "+settings.Controls.Left)
	stdscr.MovePrint(3, 10, "Right: "+settings.Controls.Right)
	stdscr.MovePrint(4, 10, "Shoot: "+settings.Controls.Shoot)

	// Refresh the screen
	stdscr.Refresh()

	// Wait for user input
	stdscr.GetChar()
	for {

	}
}

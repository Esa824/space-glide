package main

import (
	"encoding/json"
	gc "github.com/rthornton128/goncurses"
	"io"
	"log"
	"os"
)

// Define a struct to represent the control settings
type Level struct {
	Number  int `json:"number"`
	Enemies int `json:"enemies"`
	Time    int `json:"time"`
	Score   int `json:"score"`
}

type Levels struct {
	Levels []Level `json:"levels"`
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
	file, err := os.Open("levels.json")
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
	var levels Levels

	// Unmarshal the JSON data into the settings variable
	if err := json.Unmarshal(data, &levels); err != nil {
		log.Fatal(err)
	}

	// Display control settings
	stdscr.MovePrintf(0, 10, "%d", levels.Levels[0].Number)
	stdscr.MovePrintf(1, 10, "%d", levels.Levels[0].Time)
	stdscr.MovePrintf(2, 10, "%d", levels.Levels[0].Enemies)
	stdscr.MovePrintf(3, 10, "%d", levels.Levels[0].Score)

	// Refresh the screen
	stdscr.Refresh()

	// Wait for user input
	stdscr.GetChar()
	for {

	}
}

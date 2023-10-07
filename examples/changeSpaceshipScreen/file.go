package main

import (
	"encoding/json"
	"github.com/rthornton128/goncurses"
	"io/ioutil"
	"log"
	"os"
)

// Define a struct to represent the character data
type Character struct {
	Name       string   `json:"name"`
	AsciiArt   []string `json:"ascii_art"`
	Attributes struct {
		Speed  int    `json:"speed"`
		Damage int    `json:"damage"`
		Color  string `json:"color"`
	} `json:"attributes"`
}

type Characters struct {
	Characters []Character `json:"characters"`
}

func main() {
	// Open the JSON file for reading
	file, err := os.Open("characters.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the JSON data from the file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Create a variable to hold the character data
	var characters Characters

	// Unmarshal the JSON data into the characters variable
	if err := json.Unmarshal(data, &characters); err != nil {
		log.Fatal(err)
	}

	// Initialize goncurses
	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("Init:", err)
	}
	defer goncurses.End()

	goncurses.StartColor()
	goncurses.Cursor(0)
	goncurses.Echo(false)
	stdscr.Keypad(true)
	stdscr.Timeout(0)

	// Calculate the width of the characters' display area
	_, maxX := stdscr.MaxYX()
	displayWidth := maxX / 3

	// Initialize the character index to display in the middle
	currentCharacterIndex := 1

	for {
		// Clear the screen
		stdscr.Clear()

		// Display characters in the left, middle, and right
		for i, character := range characters.Characters {
			x := (i % 3) * displayWidth + 12
			if i == currentCharacterIndex {
				stdscr.AttrOn(goncurses.A_BOLD)
			}

			// Display character name
			stdscr.MovePrint(1, x, character.Name)

			// Display character ASCII art
			for j, line := range character.AsciiArt {
				stdscr.MovePrint(j+2, x, line)
			}

			// Display character attributes
			stdscr.MovePrintf(len(character.AsciiArt)+2, x, "Speed: %d", character.Attributes.Speed)
			stdscr.MovePrintf(len(character.AsciiArt)+3, x, "Damage: %d", character.Attributes.Damage)
			stdscr.MovePrintf(len(character.AsciiArt)+4, x, "Color: %s", character.Attributes.Color)

			if i == currentCharacterIndex {
				stdscr.AttrOff(goncurses.A_BOLD)
			}
		}

		// Refresh the screen
		stdscr.Refresh()

		// Listen for user input
		key := stdscr.GetChar()

		// Handle navigation
		switch key {
		case 'q':
			// Quit the program
			return
		case goncurses.KEY_RIGHT:
			// Move to the next character on the right
			if currentCharacterIndex < len(characters.Characters)-1 {
				currentCharacterIndex++
			}
		case goncurses.KEY_LEFT:
			// Move to the next character on the left
			if currentCharacterIndex > 0 {
				currentCharacterIndex--
			}
		}
	}
}

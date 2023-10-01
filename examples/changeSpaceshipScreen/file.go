package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

// Character represents a character with a name, ASCII art, and attributes.
type Character struct {
	Name      string            `json:"name"`
	ASCIIArt  []string          `json:"ascii_art"`
	Attributes map[string]interface{} `json:"attributes"`
}

func main() {
	// Read the characters.json file
	charactersJSON, err := ioutil.ReadFile("characters.json")
	if err != nil {
		log.Fatal("Error reading characters.json:", err)
	}

	// Parse the JSON data into a slice of Character
	var characters []Character
	if err := json.Unmarshal(charactersJSON, &characters); err != nil {
		log.Fatal(err)
	}

	// Initialize ncurses
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal("gc.Init:", err)
	}
	defer gc.End()

	// Create a window for character selection
	lines, cols := stdscr.MaxYX()
	charSelectionWin, err := gc.NewWindow(lines-2, cols-2, 1, 1)
	if err != nil {
		log.Fatal("gc.NewWindow:", err)
	}
	charSelectionWin.Keypad(true)
	charSelectionWin.Box(0, 0)
	charSelectionWin.Refresh()

	// Display characters in the window
	for i, character := range characters {
		charSelectionWin.MovePrint(i+1, 1, fmt.Sprintf("%d. %s", i+1, character.Name))
	}

	// Wait for user input to select a character
	for {
		key := charSelectionWin.GetChar()
		if key == 27 { // ESC key to exit
			break
		}
		if key >= '1' && key <= '9' {
			index := int(key - '1')
			if index < len(characters) {
				selectedCharacter := characters[index]
				displayCharacterInfo(selectedCharacter, stdscr)
			}
		}
	}
}

func displayCharacterInfo(character Character, stdscr *gc.Window) {
	// Clear the screen
	stdscr.Clear()
	stdscr.Refresh()

	// Create a window for displaying character info
	lines, cols := stdscr.MaxYX()
	charInfoWin, err := gc.NewWindow(lines-2, cols-2, 1, 1)
	if err != nil {
		log.Fatal("gc.NewWindow:", err)
	}
	charInfoWin.Keypad(true)
	charInfoWin.Box(0, 0)
	charInfoWin.Refresh()

	// Display character name
	charInfoWin.MovePrint(1, 1, "Name: "+character.Name)

	// Display ASCII art
	for i, line := range character.ASCIIArt {
		charInfoWin.MovePrint(i+3, 1, line)
	}

	// Display attributes
	attrLines := make([]string, 0)
	for key, value := range character.Attributes {
		attrLines = append(attrLines, fmt.Sprintf("%s: %v", key, value))
	}
	attrText := strings.Join(attrLines, ", ")
	charInfoWin.MovePrint(len(character.ASCIIArt)+3, 1, "Attributes: "+attrText)

	// Wait for user input to go back
	for {
		key := charInfoWin.GetChar()
		if key == 27 { // ESC key to go back
			break
		}
	}
}

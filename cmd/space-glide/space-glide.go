package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	gc "github.com/rthornton128/goncurses"
	log "github.com/sirupsen/logrus"
)

/* star_density is how dense the stars are the higher the density the lower the density */
const star_density = 0.005

/* planet_density is how dense the planets are the higher the density the lower the density */
const planet_density = 0.0005

/* ship_ascii is the ascii art for the spaceship */
var ship_ascii = []string{
	` ,     `,
	` |\-   `,
	`>|^===0`,
	` |/-    `,
	` '      `,
}

/* enemy_ascii is the ascii art for the enemies */
var enemy_ascii = []string{
	`  ^^^  `,
	`{(-+-)}`,
	`{(-+-)}`,
	`{(-+-)}`,
	`  ***  `,
}

/* explosion_ascii is the ascii art for an explosion */
var explosion_ascii = []string{
	`.** *. *`,
	`*. *.*. `,
	`.** **.*`,
	`*.*.***.`,
}

/* skipMainMenu is if you want to skip the main menu if it is true you skip the main menu otherwise you don't  */
var skipMainMenu bool

/* numberOfLevel is the level that you chose used for skipMainMenu */
var numberOfLevel int

/* A json structure for a Character */
type Character struct {
	Name       string   `json:"name"`      /* This is the name of the character */
	AsciiArt   []string `json:"ascii_art"` /* This is the ascii art of the character */
	Attributes struct { /* These is the attributes of the character */
		Speed  int    `json:"speed"`  /* The speed of the character */
		Damage int    `json:"damage"` /* The amount of damage the character can take (health) */
		Color  string `json:"color"`  /* The colour of the character */
	} `json:"attributes"`
}

/* A json structure for the controls */
type Controls struct {
	Up    string `json:"up"`    /* This is the control for moving up */
	Down  string `json:"down"`  /* This is the control for moving down */
	Left  string `json:"left"`  /* This is the control for moving left */
	Right string `json:"right"` /* This is the control for moving right */
	Shoot string `json:"shoot"` /* This is the control for shooting */
}

/* A json structure for a level */
type Level struct {
	Number  int `json:"number"`
	Enemies int `json:"enemies"`
	Time    int `json:"time"`
	Score   int `json:"score"`
}

/* A json structure for all of the levels */
type Levels struct {
	Levels []Level `json:"levels"`
}

/* A json structure for the all of the settings */
type Settings struct {
	Controls Controls `json:"controls"`
}

/* A json structure for the all of the characters */
type Characters struct {
	Characters []Character `json:"characters"`
}

/* A function that allows you to change or select a level */
func SelectLevel(stdscr *gc.Window) Level {
	file, err := os.Open("json/levels.json")
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
	stdscr.Clear()
	contents, err := readFile("design/levels_menu.txt")
	if err != nil {
		log.Fatal(err)
	}
	stdscr.MovePrint(0, 0, contents)
	stdscr.Timeout(-1)
	if skipMainMenu {
		skipMainMenu = !skipMainMenu
		stdscr.Timeout(0)
		return levels.Levels[numberOfLevel-1]
	}
	for {
		input := stdscr.GetChar()
		input2 := stdscr.GetChar()
		if err != nil {
			log.Fatal(err)
		}
		var inputAsNumber int
		var err error
		log.Infof("Number selected as string: %s%s", string(input), string(input2))
		if string(input) == "" {
			inputAsNumber, err = strconv.Atoi(string(input2))
		} else {
			inputAsNumber, err = strconv.Atoi(string(input) + string(input2))
		}
		log.Infof("Number selected as integer: %d", inputAsNumber)
		if err != nil {
			log.Fatal(err)
		}
		if inputAsNumber >= 1 && inputAsNumber <= 40 {
			log.Infof("Passed Check with number: %d", inputAsNumber)
			numberOfLevel = inputAsNumber
			stdscr.Timeout(0)
			return levels.Levels[inputAsNumber-1]
		}
	}
}

/* A function that allows you to change between spaceships */
func changeShip(stdscr *gc.Window) Character {
	file, err := os.Open("json/characters.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the JSON data from the file
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Create a variable to hold the character data
	var characters Characters

	// Unmarshal the JSON data into the characters variable
	if err := json.Unmarshal(data, &characters); err != nil {
		log.Fatal(err)
	}
	_, maxX := stdscr.MaxYX()
	displayWidth := maxX / 3

	// Initialize the character index to display in the middle
	currentCharacterIndex := 1
	for {
		// Clear the screen
		stdscr.Clear()

		// Display characters in the left, middle, and right
		for i, character := range characters.Characters {
			x := (i%3)*displayWidth + 12
			if i == currentCharacterIndex {
				stdscr.AttrOn(gc.A_BOLD)
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
				stdscr.AttrOff(gc.A_BOLD)
			}
		}

		// Refresh the screen
		stdscr.Refresh()

		// Listen for user input
		key := stdscr.GetChar()

		// Handle navigation
		switch key {
		case gc.KEY_RIGHT:
			// Move to the next character on the right
			if currentCharacterIndex < len(characters.Characters)-1 {
				currentCharacterIndex++
			}
		case gc.KEY_LEFT:
			// Move to the next character on the left
			if currentCharacterIndex > 0 {
				currentCharacterIndex--
			}
		case gc.KEY_RETURN:
			ship_ascii = characters.Characters[currentCharacterIndex].AsciiArt
			return characters.Characters[currentCharacterIndex]
		}
	}
}

/* A function that generates a field (*gc.Pad) and returns it */
func genStarfield(pl, pc int) *gc.Pad {
	pad, err := gc.NewPad(pl, pc)
	if err != nil {
		log.Fatal(err)
	}
	stars := int(float64(pc*pl) * star_density)
	planets := int(float64(pc*pl) * planet_density)
	for i := 0; i < stars; i++ {
		y, x := rand.Intn(pl), rand.Intn(pc)
		c := int16(rand.Intn(4) + 1)
		pad.AttrOn(gc.A_BOLD | gc.ColorPair(c))
		pad.MovePrint(y, x, ".")
		pad.AttrOff(gc.A_BOLD | gc.ColorPair(c))
	}
	for i := 0; i < planets; i++ {
		y, x := rand.Intn(pl), rand.Intn(pc)
		c := int16(rand.Intn(2) + 5)
		pad.ColorOn(c)
		if i%2 == 0 {
			pad.MoveAddChar(y, x, 'O')
		}
		pad.MoveAddChar(y, x, 'o')
		pad.ColorOff(c)
	}
	return pad
}

/* A function that populates the settings doesn't show the controls menu */
func NewControls() Controls {
	// Open the JSON file for reading
	file, err := os.Open("json/settings.json")
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
	return settings.Controls
}

/* A function that allows you to see and change the users controls and populates the controls */
func controls(stdscr *gc.Window) Controls {
	stdscr.Clear()
	lines, cols := stdscr.MaxYX()
	centerX := (cols - 120) / 2
	centerY := (lines - 40) / 2

	// Open the JSON file for reading
	file, err := os.Open("json/settings.json")
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
	controlsArt := []string{
		"000000000 0 0 0 0 0 0 0 0 0 0 0 0           0 00000000000 000000000 0 0 0 0 0 0 0 0 0 0 0 0           000000000",
		"0         0                   0 0 0         0      0      0       0 0                   0 0           0        ",
		"0         0                   0 0  0        0      0      0       0 0                   0 0           0        ",
		"0         0                   0 0   0       0      0      000000000 0                   0 0           0        ",
		"0         0                   0 0    0      0      0      0         0                   0 0           0        ",
		"0         0                   0 0     0     0      0      00        0                   0 0           000000000",
		"0         0                   0 0      0    0      0      0 0       0                   0 0                   0",
		"0         0                   0 0       0   0      0      0  0      0                   0 0                   0",
		"0         0                   0 0        0  0      0      0   0     0                   0 0                   0",
		"0         0                   0 0         0 0      0      0    0    0                   0 0                   0",
		"000000000 0 0 0 0 0 0 0 0 0 0 0 0          0       0      0     0   0 0 0 0 0 0 0 0 0 0 0 00000000000 000000000",
	}
	buttonControlsArt := []string{
		"     ___  ",
		"    | . | ",
		"   -|___|-",
		"    | . | ",
		"    |___| ",
	}
	for i, line := range controlsArt {
		stdscr.MovePrint(centerY+i, centerX, line)
	}
	for i := range buttonControlsArt {
		arg1, arg2 := settings.Controls.ReturnControlForNumber(i)
		for j, line := range buttonControlsArt {
			if j == 2 {
				if arg1 == "shoot" && arg2 == " " {
					arg2 = "space"
				}
				stdscr.MovePrintf(centerY+j+10+i*5, centerX+45-len(arg1), strconv.Itoa(i+1)+". "+"%s-|___|-%s", arg1, arg2)
				continue
			}
			stdscr.MovePrint(centerY+j+10+i*5, centerX+45, line)
		}
	}
	stdscr.Refresh()
	for {
		stdscr.Timeout(-1)

		controlNumber := stdscr.GetChar()

		if int(controlNumber) == 27 {
			break
		}
		dataForControl := stdscr.GetChar()
		control, _ := settings.Controls.ReturnControlForNumber(int(controlNumber) - 49) // Subtract 49 to get the correct index
		if control != "" && dataForControl != 0 {
			settings.Controls.SetControlForString(string(dataForControl), control)
			settings.Controls.ChangeControlButtonArtUpdated(string(dataForControl), control, stdscr, centerY, centerX, int(controlNumber)-49)

			// Marshal the updated controls back to JSON
			newData, err := json.Marshal(settings)
			if err != nil {
				log.Fatal(err)
			}

			// Write the JSON data back to the file
			if err := os.WriteFile("json/settings.json", newData, 0644); err != nil {
				log.Fatal(err)
			}
		}
	}
	stdscr.Timeout(0)
	return settings.Controls
}

/* A method that Changes the control printed if a control is changed after the controls are printed */
func (c *Controls) ChangeControlButtonArtUpdated(dataForControl string, control string, stdscr *gc.Window, centerY, centerX int, controlNumber int) {
	arg1, arg2 := c.ReturnControlForNumber(controlNumber)
	if arg1 == "shoot" && arg2 == " " {
		arg2 = "space"
	}
	if arg1 == "shoot" && arg2 != " " {
		arg2 = arg2 + "    "
	}
	stdscr.MovePrintf(centerY+2+10+controlNumber*5, centerX+45-len(arg1), strconv.Itoa(controlNumber+1)+". "+"%s-|___|-%s", arg1, arg2)
}

/* A method that changes the controls with 'dataForControl' and 'control' which are 'what you want to put for the control' and 'which control'  */
func (c *Controls) SetControlForString(dataForControl string, control string) {
	switch control {
	case "up":
		{
			c.Up = dataForControl
		}
	case "down":
		{
			c.Down = dataForControl

		}
	case "left":
		{
			c.Left = dataForControl

		}
	case "right":
		{

			c.Right = dataForControl

		}
	case "shoot":
		{
			c.Shoot = dataForControl
		}
	}
}

/* A method that takes a number and returns the data for that number */
func (c *Controls) ReturnControlForNumber(n int) (string, string) {
	switch n {
	case 0:
		{
			return "up", c.Up
		}
	case 1:
		{
			return "down", c.Down
		}
	case 2:
		{
			return "left", c.Left
		}
	case 3:
		{
			return "right", c.Right
		}
	case 4:
		{
			return "shoot", c.Shoot
		}
	default:
		{
			return "", ""
		}
	}
}

/* A function that prints the game over menu */
func gameOverMenu(stdscr *gc.Window) bool {
	lines, cols := stdscr.MaxYX()
	centerX := (cols + 158) / 2
	centerY := (lines - 40) / 2
	content, err := readFile("design/death_menu.txt")
	if err != nil {
		log.Fatal(err)
	}
	stdscr.MovePrint(centerY, centerX, content)
	stdscr.Refresh()
	for {
		input := stdscr.GetChar()
		switch int(input) {
		case '1':
			{
				return true
			}
		case '2':
			{
				return false
			}
		}
	}
}

/* A function that prints the main menu */
func showMenu(stdscr *gc.Window) rune {
	if skipMainMenu {
		return '1'
	}
	stdscr.Clear()
	stdscr.Erase()
	stdscr.Refresh()

	leftBullet := newBullet(19, 19, 1)
	rightBullet := newBullet(19, 123, -1)
	objects = append(objects, leftBullet)
	objects = append(objects, rightBullet)
	contents, err := readFile("design/main_menu.txt")
	if err != nil {
		log.Fatal(err)
	}
	stdscr.MovePrint(0, 0, contents)
	stdscr.Refresh()
	bulletTicker := time.NewTicker(time.Second / 16)
	for {
		select {
		case <-bulletTicker.C:
			lefty, leftx := leftBullet.YX()
			righty, rightx := rightBullet.YX()
			log.Infof("Printing a bullet on the screen: y: %d x: %d", lefty, leftx)
			stdscr.MovePrint(lefty, leftx, " ")
			stdscr.MovePrint(righty, rightx, " ")
			if leftx == rightx {
				leftBullet.Erase()
				rightBullet.Erase()
				leftBullet = newBullet(19, 19, 1)
				rightBullet = newBullet(19, 123, -1)
				objects = append(objects, leftBullet)
				objects = append(objects, rightBullet)
			}
			leftBullet.Update()
			rightBullet.Update()
			drawObjects(stdscr)
		default:
			key := stdscr.GetChar()
			if key >= '1' && key <= '9' {
				leftBullet.Erase()
				rightBullet.Erase()
				return rune(key)
			}
		}
	}
}

/* A method that handles input for the spaceship */
func (s *Ship) handleInput(stdscr *gc.Window, settings Settings, character *Character) {
	lines, cols := stdscr.MaxYX()
	y, x := s.YX()
	k := stdscr.GetChar()
	if settings.Controls.Shoot == "space" {
		settings.Controls.Shoot = " "
	}
	switch byte(k) {
	case 0:
		break
	case byte([]rune(settings.Controls.Left)[0]):
		x -= character.Attributes.Speed
		if x < 2 {
			x = 2
		}
	case byte([]rune(settings.Controls.Right)[0]):
		x += character.Attributes.Speed
		if x > cols-3 {
			x = cols - 3
		}
	case byte([]rune(settings.Controls.Down)[0]):
		y += character.Attributes.Speed
		if y > lines-4 {
			y = lines - 4
		}
	case byte([]rune(settings.Controls.Up)[0]):
		y -= character.Attributes.Speed
		if y < 2 {
			y = 2
		}
	case byte([]rune(settings.Controls.Shoot)[0]):
		objects = append(objects, newBullet(y+1, x+4, 1))
		objects = append(objects, newBullet(y+3, x+4, 1))
	default:
	}
	s.MoveWindow(y, x)
	for _, ob := range objects {
		if b, ok := ob.(*Bullet); ok {
			bullety, bulletx := b.YX()
			if bullety >= y && bullety <= y+5 && bulletx >= x && bulletx <= x+6 && b.dirX == -1 {
				b.handleEnemyBullet(s, y, x)
				break
			} else {
				b.handleSpaceshipBullet(s)
			}
		}
	}
}

/* Handles the enemy's bullets */
func (b *Bullet) handleEnemyBullet(s *Ship, y, x int) {
	objects = append(objects, newExplosion(y, x))
	b.alive = false
	s.life--
}

/* Handles the spaceship's bullets */
func (b *Bullet) handleSpaceshipBullet(s *Ship) {
	bullety, bulletx := b.YX()
	for _, ob := range objects {
		if enemy, ok := ob.(*EnemyShip); ok {
			y, x := enemy.YX()
			if bullety >= y && bullety <= y+5 && bulletx >= x && bulletx <= x+6 {
				objects = append(objects, newExplosion(y, x))
				b.alive = false
				enemy.Clear()
				enemy.alive = false
				s.Score++
				break
			}
		}
	}
}

/* An interface for any object such as the spaceship, the enemies, the stars, the explosion, and the bullets */
type Object interface {
	Cleanup()
	Draw(*gc.Window)
	Expired(int, int) bool
	Update()
}

/* A struct for the bullets */
type Bullet struct {
	*gc.Window
	alive bool
	dirX  int
}

/* A function that creates a new bullet */
func newBullet(y, x int, dirX int) *Bullet {
	w, err := gc.NewWindow(1, 1, y, x)
	if err != nil {
		log.Println("newBullet:", err)
	}
	w.AttrOn(gc.A_BOLD | gc.ColorPair(4))
	w.Print("-")
	return &Bullet{w, true, dirX}
}

/* A function that deletes a bullet */
func (b *Bullet) Cleanup() {
	b.Delete()
}

/* A function that draws a bullet */
func (b *Bullet) Draw(w *gc.Window) {
	w.Overlay(b.Window)
}

/* A function that checks if a bullet has expired/died/offTheScreen */
func (b *Bullet) Expired(my, mx int) bool {
	_, x := b.YX()
	if x >= mx-1 || !b.alive {
		return true
	}
	return false
}

/* A function that updates the bullet */
func (b *Bullet) Update() {
	y, x := b.YX()
	if x == 255 || x == 0 {
		b.Erase()
		return
	}
	b.MoveWindow(y, x+b.dirX) // Update the bullet's x-coordinate based on direction
}

/* A struct for the spaceship */
type Ship struct {
	*gc.Window
	life  int
	Score int
}

/* A struct for the Explosions animation */
type Explosion struct {
	*gc.Window
	life int
}

/* A function that makes a new Explosion animation */
func newExplosion(y, x int) *Explosion {
	w, err := gc.NewWindow(4, 8, y-1, x-1)
	if err != nil {
		log.Println("newExplosion:", err)
	}
	w.ColorOn(4)
	for i, line := range explosion_ascii {
		w.MovePrint(i, 0, line)
	}
	w.ColorOff(4)
	return &Explosion{w, 5}
}

/* A function that deletes the Explosion animation */
func (e *Explosion) Cleanup() {
	e.Delete()
}

/* An empty function just to make it so that a explosion can fit into the object interface */
func (e *Explosion) Collide(i int) {}

/* A function that draws the Explosion animation */
func (e *Explosion) Draw(w *gc.Window) {
	w.Overlay(e.Window)
}

/* A function that is used just to make it so that a explosion can fit into the object interface and if the explosion animation is gone */
func (e *Explosion) Expired(y, x int) bool {
	return e.life <= 0
}

/* A function that is used just to make it so that a explosion can fit into the object interface */
func (e *Explosion) Update() {
	e.life--
}

/* A function that is used to delete the spaceship */
func (s *Ship) Cleanup() {
	s.Delete()
}

/* A function that makes the new spaceship */
func newShip(y, x int, character *Character) *Ship {
	w, err := gc.NewWindow(5, 7, y, x)
	if err != nil {
		log.Fatal("newShip:", err)
	}

	// Determine the color pair based on the character's color attribute
	var colorPair gc.Char
	switch character.Attributes.Color {
	case "blue":
		colorPair = gc.ColorPair(5) // Assuming color pair 1 is blue
	case "red":
		colorPair = gc.ColorPair(4) // Assuming color pair 2 is red
	case "green":
		colorPair = gc.ColorPair(6) // Assuming color pair 3 is green
	default:
		colorPair = gc.ColorPair(0) // Default to the default color pair (usually white text on black background)
	}

	// Set the color pair for the ship's window
	w.AttrOn(colorPair)

	for i := 0; i < len(ship_ascii); i++ {
		w.MovePrint(i, 0, ship_ascii[i])
	}

	// Turn off color pair to avoid affecting subsequent output
	w.AttrOff(colorPair)
	if character.Attributes.Damage == 0 {
		character.Attributes.Damage = 5
	}
	if character.Attributes.Speed == 0 {
		character.Attributes.Speed = 1
	}

	return &Ship{w, character.Attributes.Damage, 0}
}

/* A function that reads a file and returns the contents and an error if there is one while reading a file */
func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

/* A function that draws the spaceship */
func (s *Ship) Draw(w *gc.Window) {
	w.Overlay(s.Window)
}

/* A function that checks if the spaceship has died */
func (s *Ship) Expired(y, x int) bool {
	return s.life <= 0
}

/* A function just to make it so that the spaceship can be part of the object */
func (s *Ship) Update() {}

/* A struct for the ememies spaceships */
type EnemyShip struct {
	*gc.Window
	alive        bool
	shootTicker  *time.Ticker // Ticker to control enemy ship shooting
	bulletSymbol string       // Symbol for enemy ship bullets
	bulletDirX   int          // X-direction for enemy ship bullets (-1 for left, 1 for right)
}

/* A function that makes a new enemy ship */
func newEnemyShip(y, x int) *EnemyShip {
	w, err := gc.NewWindow(5, 7, y, x)
	if err != nil {
		log.Fatal("newEnemyShip:", err)
	}
	for i := 0; i < len(enemy_ascii); i++ {
		w.MovePrint(i, 0, enemy_ascii[i])
	}
	// Create a ticker for enemy ship shooting
	shootTicker := time.NewTicker(time.Second * 2)
	return &EnemyShip{w, true, shootTicker, "-", -1}
}

/* A function that deletes the enemy ship */
func (e *EnemyShip) Cleanup() {
	e.Delete()
}

/* A function that draws the ememy ship */
func (e *EnemyShip) Draw(w *gc.Window) {
	w.Overlay(e.Window)
}

/* A function that checks if the ememy ship has expired/died/goneOffTheScreen */
func (e *EnemyShip) Expired(my, mx int) bool {
	_, x := e.YX()
	if x <= mx-1 || !e.alive {
		return true
	}
	return false
}

/* A function that updates a enemy ship */
func (e *EnemyShip) Update() {
	y, x := e.YX()
	e.MoveWindow(y, x-1)

	select {
	case <-e.shootTicker.C:
		// Create bullets for enemy ships when they shoot, but in the opposite direction
		objects = append(objects, newBullet(y+1, x-1, e.bulletDirX)) // Adjust the direction here
		objects = append(objects, newBullet(y+3, x-1, e.bulletDirX)) // Adjust the direction here
	default:
		// Do nothing if it's not time to shoot
	}
}

/* A variable for all objects */
var objects = make([]Object, 0, 16)

/* A function used to update objects */
func updateObjects(my, mx int) {
	end := len(objects)
	tmp := make([]Object, 0, end)
	for _, ob := range objects {
		ob.Update()
	}
	for _, ob := range objects {
		if ob.Expired(my, mx) {
			ob.Cleanup()
		} else {
			tmp = append(tmp, ob)
		}
	}
	if len(objects) > end {
		objects = append(tmp, objects[end:]...)
	} else {
		objects = tmp
	}
}

/* A function used to draw objects  */
func drawObjects(s *gc.Window) {
	for _, ob := range objects {
		ob.Draw(s)
	}
}

/* A function that takes a number and returns a string of '*' as hearts */
func lifeToText(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "*"
	}
	return s
}

func signalHandler(signals chan os.Signal) {
	<-signals
	gc.End()
	os.Exit(0)
}

/* The main function where everything starts */
func main() {
	// Logging
	logFile, err := os.OpenFile("/tmp/space-glide.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		print("could not open log file")
		os.Exit(1)
	}
	log.SetOutput(logFile)

	var stdscr *gc.Window
	stdscr, err = gc.Init()
	if err != nil {
		log.Println("Init:", err)
	}
	defer gc.End()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(sigCh)

	gc.StartColor()
	gc.Cursor(0)
	gc.Echo(false)
	stdscr.Keypad(true)
	stdscr.Timeout(0)

	lines, cols := stdscr.MaxYX()
	pl, pc := lines, cols*3
	field := genStarfield(pl, pc)

	character := Character{}
	settings := Settings{}
	settings.Controls = NewControls()
	level := Level{}

	gc.InitPair(1, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(2, gc.C_YELLOW, gc.C_BLACK)
	gc.InitPair(3, gc.C_MAGENTA, gc.C_BLACK)
	gc.InitPair(4, gc.C_RED, gc.C_BLACK)

	gc.InitPair(5, gc.C_BLUE, gc.C_BLACK)
	gc.InitPair(6, gc.C_GREEN, gc.C_BLACK)

	stdscr.Clear()
	for {
		key := showMenu(stdscr)
		if key == '2' {
			character = changeShip(stdscr)
		}
		if key == '3' {
			settings.Controls = controls(stdscr)
		}
		if key == '4' {
			break
		} else if key != '1' {
			continue
		}
		if key == '1' {
			level = SelectLevel(stdscr)
		}

		lines, cols := stdscr.MaxYX()

		ship := newShip(lines/2, 5, &character)
		objects = append(objects, ship)

		text := stdscr.Duplicate()

		c := time.NewTicker(time.Second / 16)
		c2 := time.NewTicker(time.Second / 16)
		px := 0

		enemyTicker := time.NewTicker(time.Second * 2) // Create a ticker for spawning enemy ships
		timeLeftTicker := time.NewTicker(time.Second)  // Create a ticker for spawning enemy ships

	loop:
		for {
			text.Erase()
			stdscr.Erase()
			text.MovePrintf(0, 0, "Life: [%-"+strconv.Itoa(character.Attributes.Damage)+"s]", lifeToText(ship.life))
			text.MovePrintf(0, 20, "Score: %d", ship.Score)
			text.MovePrintf(0, 40, "TimeLeft: %ds", level.Time)
			stdscr.Copy(field.Window, 0, px, 0, 0, lines-1, cols-1, true)
			drawObjects(stdscr)
			stdscr.Overlay(text)
			stdscr.Refresh()
			ship.handleInput(stdscr, settings, &character)
			select {
			case <-c.C:
				if px+cols >= pc {
					break loop
				}
				px++
			case <-c2.C:
				updateObjects(lines, cols)
				drawObjects(stdscr)
			case <-enemyTicker.C: // Spawn enemy ships periodically
				ey := rand.Intn(lines-4) + 2 // Randomly select the y position for the enemy ship
				ex := cols - 10              // Set the x position to the right edge of the screen
				objects = append(objects, newEnemyShip(ey, ex))
			case <-timeLeftTicker.C:
				level.Time -= 1
			default:
				if ship.Expired(-1, -1) {
					break loop
				}
			}
		}
		skipMainMenu = gameOverMenu(stdscr)
	}
}

package main

import (
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"
	"os"
	"time"
)

const density = 0.05

var ship_ascii = []string{
	` ,`,
	` |\-`,
	`>|^===0`,
	` |/-`,
	` '`,
}

var enemy_ascii = []string{
	` `,
	`  /`,
	`---E`,
	`  \`,
	` `,
}

func genStarfield(pl, pc int) *gc.Pad {
	pad, err := gc.NewPad(pl, pc)
	if err != nil {
		log.Fatal(err)
	}
	stars := int(float64(pc*pl) * density)
	for i := 0; i < stars; i++ {
		y, x := rand.Intn(pl), rand.Intn(pc)
		c := int16(rand.Intn(4) + 1)
		pad.AttrOn(gc.A_BOLD | gc.ColorPair(c))
		pad.MovePrint(y, x, ".")
		pad.AttrOff(gc.A_BOLD | gc.ColorPair(c))
	}

	return pad
}

func gameOverMenu(stdscr *gc.Window) bool {

	stdscr.AttrOn(gc.A_BOLD | gc.ColorPair(4))
	lines, cols := stdscr.MaxYX()
	centerX := (cols - 56) / 2
	centerY := (lines - 20) / 2

	// Game Over ASCII Art
	gameOverArt := []string{
		"000000000     000000000           0      00000000000 0        0        ",
		"0        0    0                  0 0          0      0        0        ",
		"0         0   0                 0   0         0      0        0        ",
		"0          0  0                0     0        0      0        0        ",
		"0           0 0               0       0       0      0        0        ",
		"0           0 000000000      00000000000      0      0000000000        ",
		"0          0  0             0           0     0      0        0        ",
		"0         0   0            0             0    0      0        0        ",
		"0        0    0           0               0   0      0        0        ",
		"0       0     0          0                 0  0      0        0        ",
		"00000000      000000000 0                   0 0      0        0        ",
	}

	// Game Over Menu
	menu := []string{
		"                          __________",
		"                         |          |",
		"                         | Mainmenu |",
		"                         |__________|",
	}

	// Print Game Over ASCII Art
	for i, line := range gameOverArt {
		stdscr.MovePrint(centerY+i, centerX, line)
	}

	// Print Game Over Menu
	for i, line := range menu {
		stdscr.MovePrint(centerY+len(gameOverArt)+i, centerX, line)
	}

	stdscr.Refresh()
	stdscr.AttrOff(gc.A_BOLD | gc.ColorPair(4))
	for {
		key := stdscr.GetChar()
		if key == 'r' || key == 'R' {
			return true // Restart
		} else if key == 'm' || key == 'M' {
			return false // Mainmenu
		}
	}
}

func showMenu(stdscr *gc.Window) bool {
	stdscr.Clear()
	stdscr.MovePrint(2, 2, "Main Menu")
	stdscr.MovePrint(4, 2, "1. Start Game")
	stdscr.MovePrint(5, 2, "2. Quit")
	stdscr.Refresh()

	for {
		key := stdscr.GetChar()
		if key == '1' {
			return true
		} else if key == '2' {
			return false
		}
	}
}

func handleInput(stdscr *gc.Window, ship *Ship) bool {
	lines, cols := stdscr.MaxYX()
	y, x := ship.YX()
	k := stdscr.GetChar()

	switch byte(k) {
	case 0:
		break
	case 'a':
		x--
		if x < 2 {
			x = 2
		}
	case 'd':
		x++
		if x > cols-3 {
			x = cols - 3
		}
	case 's':
		y++
		if y > lines-4 {
			y = lines - 4
		}
	case 'w':
		y--
		if y < 2 {
			y = 2
		}
	case ' ':
		objects = append(objects, newBullet(y+1, x+4, 1))
		objects = append(objects, newBullet(y+3, x+4, 1))
	default:
		return true
	}
	ship.MoveWindow(y, x)
	for i, ob := range objects {
		if b, ok := ob.(*Bullet); ok && b.dirX == -1 {
			bullety, bulletx := b.YX()
			if bullety == y && bulletx >= x && bulletx <= x+6 {
				objects = append(objects, newExplosion(y, x))
				b.alive = false
				ship.Collide(i)
				break
			}
		}
	}
	return true
}

type Object interface {
	Cleanup()
	Draw(*gc.Window)
	Expired(int, int) bool
	Update()
}

type Bullet struct {
	*gc.Window
	alive bool
	dirX  int // Direction of the bullet (-1 for left, 1 for right)
}

func newBullet(y, x int, dirX int) *Bullet {
	w, err := gc.NewWindow(1, 1, y, x)
	if err != nil {
		log.Println("newBullet:", err)
	}
	w.AttrOn(gc.A_BOLD | gc.ColorPair(4))
	w.Print("-")
	return &Bullet{w, true, dirX}
}

func (b *Bullet) Cleanup() {
	b.Delete()
}

func (b *Bullet) Draw(w *gc.Window) {
	w.Overlay(b.Window)
}

func (b *Bullet) Expired(my, mx int) bool {
	_, x := b.YX()
	if x >= mx-1 || !b.alive {
		return true
	}
	return false
}

func (b *Bullet) Update() {
	y, x := b.YX()
	b.MoveWindow(y, x+b.dirX) // Update the bullet's x-coordinate based on direction
}

type Ship struct {
	*gc.Window
	life int
}

type Explosion struct {
	*gc.Window
	life int
}

func newExplosion(y, x int) *Explosion {
	w, err := gc.NewWindow(3, 3, y-1, x-1)
	if err != nil {
		log.Println("newExplosion:", err)
	}
	w.ColorOn(4)
	w.MovePrint(0, 0, `\ /`)
	w.AttrOn(gc.A_BOLD)
	w.MovePrint(1, 0, ` X `)
	w.AttrOn(gc.A_DIM)
	w.MovePrint(2, 0, `/ \`)
	return &Explosion{w, 5}
}

func (e *Explosion) Cleanup() {
	e.Delete()
}

func (e *Explosion) Collide(i int) {}

func (e *Explosion) Draw(w *gc.Window) {
	w.Overlay(e.Window)
}

func (e *Explosion) Expired(y, x int) bool {
	return e.life <= 0
}

func (e *Explosion) Update() {
	e.life--
}

func (s *Ship) Cleanup() {
	s.Delete()
}

func (s *Ship) Collide(i int) {
	ty, tx := s.YX()
	by, bx := s.MaxYX()
	for _, ob := range objects {
		if b, ok := ob.(*Bullet); ok && b.dirX == -1 {
			y, x := b.YX()
			if y >= ty && y <= ty+by && x >= tx && x <= tx+bx {
				objects = append(objects, newExplosion(s.YX()))
				b.alive = false
				s.life--
				break
			}
		}
	}
}

func newShip(y, x int) *Ship {
	w, err := gc.NewWindow(5, 7, y, x)
	if err != nil {
		log.Fatal("newShip:", err)
	}
	for i := 0; i < len(ship_ascii); i++ {
		w.MovePrint(i, 0, ship_ascii[i])
	}
	return &Ship{w, 5}
}

func (s *Ship) Draw(w *gc.Window) {
	w.Overlay(s.Window)
}

func (s *Ship) Expired(y, x int) bool {
	return s.life <= 0
}

func (s *Ship) Update() {}

type EnemyShip struct {
	*gc.Window
	alive        bool
	shootTicker  *time.Ticker // Ticker to control enemy ship shooting
	bulletSymbol string       // Symbol for enemy ship bullets
	bulletDirX   int          // X-direction for enemy ship bullets (-1 for left, 1 for right)
}

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

func (e *EnemyShip) Cleanup() {
	e.Delete()
}

func (e *EnemyShip) Draw(w *gc.Window) {
	w.Overlay(e.Window)
}

func (e *EnemyShip) Expired(my, mx int) bool {
	_, x := e.YX()
	if x <= 0 || !e.alive {
		return true
	}
	return false
}

func (e *EnemyShip) Collide(i int) {
	ty, tx := e.YX()
	by, bx := e.MaxYX()
	for _, ob := range objects {
		if b, ok := ob.(*Bullet); ok && b.dirX == 1 {
			y, x := b.YX()
			if y <= ty && y >= ty+by && x <= tx && x >= tx+bx {
				objects = append(objects, newExplosion(e.YX()))
				b.alive = false
				e.Cleanup()
				break
			}
		}
	}
}

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

var objects = make([]Object, 0, 16)

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

func drawObjects(s *gc.Window) {
	for _, ob := range objects {
		ob.Draw(s)
	}
}

func lifeToText(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "*"
	}
	return s
}

func main() {
	f, err := os.Create("err.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

	var stdscr *gc.Window
	stdscr, err = gc.Init()
	if err != nil {
		log.Println("Init:", err)
	}
	defer gc.End()

	rand.Seed(time.Now().Unix())
	gc.StartColor()
	gc.Cursor(0)
	gc.Echo(false)
	stdscr.Timeout(0)

	lines, cols := stdscr.MaxYX()
	pl, pc := lines, cols*3

	field := genStarfield(pl, pc)

	gc.InitPair(1, gc.C_WHITE, gc.C_BLACK)
	gc.InitPair(2, gc.C_YELLOW, gc.C_BLACK)
	gc.InitPair(3, gc.C_MAGENTA, gc.C_BLACK)
	gc.InitPair(4, gc.C_RED, gc.C_BLACK)

	gc.InitPair(5, gc.C_BLUE, gc.C_BLACK)
	gc.InitPair(6, gc.C_GREEN, gc.C_BLACK)

	for {
		if !showMenu(stdscr) {
			break
		}

		lines, cols := stdscr.MaxYX()

		ship := newShip(lines/2, 5)
		objects = append(objects, ship)

		text := stdscr.Duplicate()

		c := time.NewTicker(time.Second / 2)
		c2 := time.NewTicker(time.Second / 16)
		px := 0

		enemyTicker := time.NewTicker(time.Second * 2) // Create a ticker for spawning enemy ships

	loop:
		for {
			text.Erase()
			text.MovePrintf(0, 0, "Life: [%-5s]", lifeToText(ship.life))
			stdscr.Copy(field.Window, 0, px, 0, 0, lines-1, cols-1, true)
			stdscr.Erase()
			drawObjects(stdscr)
			stdscr.Overlay(text)
			stdscr.Refresh()
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
			default:
				if !handleInput(stdscr, ship) || ship.Expired(-1, -1) {
					break loop
				}
			}
		}
		gameOverMenu(stdscr)
		gc.Nap(2000)
	}
}

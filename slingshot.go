package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

type SlingshotGame struct {
	xSize   int
	ySize   int
	planets []Planet
	players []SlingshotPlayer
	win     *pixelgl.Window
	turn    int
	cam     *SlingshotCamera
	atlas   *text.Atlas
}

func (sg *SlingshotGame) Update() {
	sg.draw()
	sg.getInput()
	time.Sleep(100 * time.Millisecond)

}

func (sg *SlingshotGame) getInput() {
	// Control the camera
	dt := time.Since(sg.cam.last).Seconds()
	sg.cam.last = time.Now()

	sg.cam.cam = pixel.IM.Scaled(sg.cam.camPos, sg.cam.camZoom).Moved(sg.win.Bounds().Center().Sub(sg.cam.camPos))
	sg.win.SetMatrix(sg.cam.cam)

	if sg.win.Pressed(pixelgl.KeyLeft) {
		sg.cam.camPos.X -= sg.cam.camSpeed * dt
	}
	if sg.win.Pressed(pixelgl.KeyRight) {
		sg.cam.camPos.X += sg.cam.camSpeed * dt
	}
	if sg.win.Pressed(pixelgl.KeyDown) {
		sg.cam.camPos.Y -= sg.cam.camSpeed * dt
	}
	if sg.win.Pressed(pixelgl.KeyUp) {
		sg.cam.camPos.Y += sg.cam.camSpeed * dt
	}
	sg.cam.camZoom *= math.Pow(sg.cam.camZoomSpeed, sg.win.MouseScroll().Y)

	// Turn ship left
	if sg.win.Pressed(pixelgl.Key1) {
		sg.players[sg.turn].ship.angle += 10
	}

	// Turn ship right
	if sg.win.Pressed(pixelgl.Key2) {
		sg.players[sg.turn].ship.angle -= 10
	}

	// More power
	if sg.win.Pressed(pixelgl.Key3) {
		sg.players[sg.turn].ship.power += 10
	}

	//Less power
	if sg.win.Pressed(pixelgl.Key4) {
		sg.players[sg.turn].ship.power -= 10
	}

	if sg.win.Pressed(pixelgl.KeySpace) {
		sg.turn = (sg.turn + 1) % len(sg.players)
	}
}

func (sg *SlingshotGame) drawPicture(xPos, yPos, angle float64, path string) {
	pic, err := loadPicture(path)
	if err != nil {
		panic(err)
	}

	mat := pixel.IM
	mat = mat.Moved((sg.win.Bounds().Min))
	mat = mat.Moved(pixel.V(xPos, yPos))
	mat = mat.Rotated(pixel.V(xPos, yPos), 360*angle/(math.Pi))
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(sg.win, mat)

}

func NewSlingshotGame(numPlanets, numPlayers, xSize, ySize int) *SlingshotGame {

	//Load planet images
	var planetImages []string
	files, err := ioutil.ReadDir("img/planets")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		planetImages = append(planetImages, "./img/planets/"+f.Name())
	}

	//Load ship images
	var shipImages []string
	files, err = ioutil.ReadDir("img/ships")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		shipImages = append(shipImages, "./img/ships/"+f.Name())
	}

	//Create window
	cfg := pixelgl.WindowConfig{
		Title:  "Slingshot",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//Create text for info
	face, err := loadTTF("font.ttf", 30)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)

	var planets []Planet
	var players []SlingshotPlayer

	cam := NewSlingshotCamera()
	sg := &SlingshotGame{xSize, ySize, planets, players, win, 0, cam, atlas}

	// Add Planets
	for i := 0; i < numPlanets; i++ {
		xPos := rand.Float64() * float64(xSize)
		yPos := rand.Float64() * float64(ySize)
		diam := rand.Float64() * float64(100)
		sg.addPlanet(xPos, yPos, diam, planetImages[i%len(planetImages)])
	}

	// Add Players
	for i := 0; i < numPlayers; i++ {
		xPos := rand.Float64() * float64(xSize)
		yPos := rand.Float64() * float64(ySize)
		sg.addPlayer(xPos, yPos, float64((90*i)%360), shipImages[i%len(shipImages)])
	}

	return sg
}

// Add a planet to the game
func (sg *SlingshotGame) addPlanet(xPos, yPos, diameter float64, image string) {
	p := Planet{xPos, yPos, diameter, image}
	sg.planets = append(sg.planets, p)
}

// Add a player to the game
func (sg *SlingshotGame) addPlayer(xPos, yPos, angle float64, image string) {
	ship := SpaceShip{xPos, yPos, angle, 10, image}
	player := SlingshotPlayer{ship, 0}
	sg.players = append(sg.players, player)
}

// Draw all images on the screen
func (sg SlingshotGame) draw() {
	sg.win.Clear(colornames.Blue)
	sg.drawPlanets()
	sg.drawShips()
	sg.drawScore()
}

// Draw the planets
func (sg *SlingshotGame) drawPlanets() {
	for _, v := range sg.planets {
		sg.drawPicture(v.xPos, v.yPos, 0, v.image)
	}
}

// Draw the ships
func (sg *SlingshotGame) drawShips() {
	for _, v := range sg.players {
		sg.drawPicture(v.ship.xPos, v.ship.yPos, v.ship.angle, v.ship.image)
	}
}

// Draw the players score
func (sg SlingshotGame) drawScore() {
	txt := text.New(pixel.V(50, 500), sg.atlas)

	// Print Players info in respective color
	for k, v := range sg.players {
		if sg.turn == k {
			txt.Color = colornames.Red
		} else {
			txt.Color = colornames.Grey
		}
		txt.WriteString("Payer: " + string(k+1) + ": " + string(v.score) + "\n")
	}

	txt.Draw(sg.win, pixel.IM.Moved(sg.win.Bounds().Max.Sub(txt.Bounds().Max)))
}

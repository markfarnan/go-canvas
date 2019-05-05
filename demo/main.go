package main

import (
	"image/color"
	"time"

	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
)

type gameState struct{ laserX, laserY, directionX, directionY, laserSize float64 }

var done chan struct{}

var cvs *canvas.Canvas2d
var width float64
var height float64

var gs = gameState{laserSize: 35, directionX: 13.7, directionY: -13.7, laserX: 40, laserY: 40}

// This specifies how long a delay between calls to 'render'.     To get Frame Rate,   1s / renderDelay
var renderDelay time.Duration = 20 * time.Millisecond

func main() {

	FrameRate := time.Second / renderDelay
	println("Hello Browser FPS:", FrameRate)
	cvs, _ = canvas.NewCanvas2d(true)
	height = float64(cvs.Height())
	width = float64(cvs.Width())

	go doEvery(renderDelay, Render) // Kick off the Render function as go routine as it never returns
	<-done
}

// Helper function which calls the required func (in this case 'render') every time.Duration,  Call as a go-routine to prevent blocking, as this never returns
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// This is called on a timer from 'doEvery'  to render the image / screen.
// The render rate is totaly under the users control, and is seperated from the Browsers frame Rendering (which is handled inside the canvas2d class)
func Render(t time.Time) {
	// Get the painter drawing Context
	gc := cvs.Gc()

	if gs.laserX+gs.directionX > width-gs.laserSize || gs.laserX+gs.directionX < gs.laserSize {
		gs.directionX = -gs.directionX
	}
	if gs.laserY+gs.directionY > height-gs.laserSize || gs.laserY+gs.directionY < gs.laserSize {
		gs.directionY = -gs.directionY
	}

	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.Clear()
	// move red laser
	gs.laserX += gs.directionX
	gs.laserY += gs.directionY

	// draws red 🔴 laser
	gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})

	gc.BeginPath()
	//gc.ArcTo(gs.laserX, gs.laserY, gs.laserSize, gs.laserSize, 0, math.Pi*2)
	draw2dkit.Circle(gc, gs.laserX, gs.laserY, gs.laserSize)
	gc.FillStroke()
	gc.Close()

}

// Copyright [2019] [Mark Farnan]

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// NOTICE:  Much of this demo is a re-write of the 'Moving red Laser' demo by Martin Olsansky https://medium.freecodecamp.org/webassembly-with-golang-is-fun-b243c0e34f02
// It has been re-written to make use of the go-canvas library,  and avoid context calls for drawing.

package main

import (
	"github.com/markfarnan/go-canvas/canvas"
)

type gameState struct{ laserX, laserY, directionX, directionY, laserSize float64 }

var done chan struct{}

var cvs *canvas.CanvasWebGL
var width float64
var height float64

//var gs = gameState{laserSize: 35, directionX: 50, directionY: -41, laserX: 50, laserY: 50}

var gsall [200]gameState

var colourfactor uint8

func main() {

	colourfactor = 255 / uint8(len(gsall))
	offsetFactor := 600 / float64(len(gsall))
	directionFactor := 50 / float64(len(gsall))

	//factor = 200.0 / float64(len(gsall))

	for k := range gsall {
		gsall[k] = gameState{laserSize: 20, directionX: 50, directionY: -41 + float64(k)*directionFactor, laserX: float64(k) * offsetFactor, laserY: 50}
	}

	cvs, _ = canvas.NewCanvasWebgl(false)
	//cvs.Create(int(js.Global().Get("innerWidth").Float()*0.9), int(js.Global().Get("innerHeight").Float()*0.9)) // Make Canvas 90% of window size.  For testing rendering canvas smaller than full windows
	cvs.Create(600, 600) // Make Canvas 90% of window size.  For testing rendering canvas smaller than full windows

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(60, Render)

	//go doEvery(renderDelay, Render) // Kick off the Render function as go routine as it never returns
	<-done
}

// Called from the 'requestAnnimationFrame' function.   It may also be called seperatly from a 'doEvery' function, if the user prefers drawing to be seperate from the annimationFrame callback
func Render() bool {

	return true
}

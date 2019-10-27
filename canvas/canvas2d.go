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

package canvas

import (
	"image"
	"syscall/js"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

type RenderFunc func(gc *draw2dimg.GraphicContext) bool

type Canvas2d struct {
	done chan struct{} // Used as part of 'run forever' in the render handler

	// DOM properties
	window js.Value
	doc    js.Value
	body   js.Value

	// Canvas properties
	canvas js.Value
	ctx    js.Value
	im     js.Value
	width  int
	height int

	// Drawing Context
	gctx     *draw2dimg.GraphicContext // Graphic Context
	image    *image.RGBA               // The Shadow frame we actually draw on
	font     *truetype.Font
	fontData draw2d.FontData

	reqID    js.Value // Storage of the current annimationFrame requestID - For Cancel
	timeStep float64  // Min Time delay between frames. - Calculated as   maxFPS/1000
}

func NewCanvas2d(create bool) (*Canvas2d, error) {

	var c Canvas2d

	c.window = js.Global()
	c.doc = c.window.Get("document")
	c.body = c.doc.Get("body")

	// If create, make a canvas that fills the windows
	if create {
		c.Create(int(c.window.Get("innerWidth").Int()), int(c.window.Get("innerHeight").Int()))
	}

	return &c, nil
}

// Create a new Canvas in the DOM, and append it to the Body.
// This also calls Set to create relevant shadow Buffer etc

// TODO suspect this needs to be fleshed out with more options
func (c *Canvas2d) Create(width int, height int) {

	// Make the Canvas
	canvas := c.doc.Call("createElement", "canvas")

	canvas.Set("height", height)
	canvas.Set("width", width)
	c.body.Call("appendChild", canvas)

	c.Set(canvas, width, height)
}

// Used to setup with an existing Canvas element which was obtained from JS
func (c *Canvas2d) Set(canvas js.Value, width int, height int) {
	c.canvas = canvas
	c.height = height
	c.width = width

	// Setup the 2D Drawing context
	c.ctx = c.canvas.Call("getContext", "2d")
	c.im = c.ctx.Call("createImageData", width, height) // Note Width, then Height

	c.image = image.NewRGBA(image.Rect(0, 0, width, height))
	c.gctx = draw2dimg.NewGraphicContext(c.image)

	// init font
	c.font, _ = truetype.Parse(FontData["font.ttf"])

	c.fontData = draw2d.FontData{
		Name:   "roboto",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	}
	fontCache := &FontCache{}
	fontCache.Store(c.fontData, c.font)

	c.gctx.FontCache = fontCache
}

// Starts the annimationFrame callbacks running.   (Recently seperated from Create / Set to give better control for when things start / stop)
func (c *Canvas2d) Start(maxFPS float64, rf RenderFunc) {
	c.SetFPS(maxFPS)
	c.initFrameUpdate(rf)
}

// This needs to be called on an 'beforeUnload' trigger, to properly close out the render callback, and prevent browser errors on page Refresh
func (c *Canvas2d) Stop() {
	c.window.Call("cancelAnimationFrame", c.reqID)
	c.done <- struct{}{}
	close(c.done)
}

// Sets the maximum FPS (Frames per Second).  This can be changed on the fly and will take affect next frame.
func (c *Canvas2d) SetFPS(maxFPS float64) {
	c.timeStep = 1000 / maxFPS
}

// Get the Drawing context for the Canvas
func (c *Canvas2d) Gc() *draw2dimg.GraphicContext {
	return c.gctx
}

func (c *Canvas2d) Height() int {
	return c.height
}

func (c *Canvas2d) Width() int {
	return c.width
}

// handles calls from Render, and copies the image over.
func (c *Canvas2d) initFrameUpdate(rf RenderFunc) {
	// Hold the callbacks without blocking
	go func() {
		var renderFrame js.Func
		var lastTimestamp float64

		renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			timestamp := args[0].Float()
			if timestamp-lastTimestamp >= c.timeStep { // Constrain FPS
				if rf != nil { // If required, call the requested render function, before copying the frame
					if rf(c.gctx) { // Only copy the image back if RenderFunction returns TRUE. (i.e. stuff has changed.)  This allows Render to return false, saving time this cycle if nothing changed.  (Keep frame as before)
						c.imgCopy()
					}
				} else { // Just do the copy, rendering must be being done elsewhere
					c.imgCopy()
				}
				lastTimestamp = timestamp
			}

			c.reqID = js.Global().Call("requestAnimationFrame", renderFrame) // Captures the requestID to be used in Close / Cancel
			return nil
		})
		defer renderFrame.Release()
		js.Global().Call("requestAnimationFrame", renderFrame)
		<-c.done
	}()
}

// Does the actuall copy over of the image data for the 'render' call.
func (c *Canvas2d) imgCopy() {
	// golang buffer
	ta := js.TypedArrayOf(c.image.Pix)
	c.im.Get("data").Call("set", ta)
	ta.Release()
	c.ctx.Call("putImageData", c.im, 0, 0)
}

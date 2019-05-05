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

type Canvas2d struct {
	done chan struct{} // Used as part of 'run forever' in the render handler

	// DOM properties
	window     js.Value
	doc        js.Value
	body       js.Value
	windowSize struct{ w, h float64 }

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
}

func NewCanvas2d(create bool) (*Canvas2d, error) {

	var c Canvas2d

	c.window = js.Global()
	c.doc = c.window.Get("document")
	c.body = c.doc.Get("body")

	c.windowSize.h = c.window.Get("innerHeight").Float()
	c.windowSize.w = c.window.Get("innerWidth").Float()

	// If create, make a canvas that fills the windows
	if create {
		c.Create(int(c.windowSize.w), int(c.windowSize.h))
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
	c.im = c.ctx.Call("createImageData", c.windowSize.w, c.windowSize.h) // Note Width, then Height

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

	// Kick off the render callback routine.
	c.initFrameUpdate()
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
func (c *Canvas2d) initFrameUpdate() {
	// Hold the callbacks without blocking
	go func() {
		var renderFrame js.Func
		renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.imgCopy()
			js.Global().Call("requestAnimationFrame", renderFrame)
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

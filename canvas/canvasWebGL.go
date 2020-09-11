package canvas

import "syscall/js"

type WebGLRenderFunc func(c *CanvasWebGL, gl js.Value)

type CanvasWebGL struct {
	done chan struct{} // Used as part of 'run forever' in the render handler

	// DOM properties
	window js.Value
	doc    js.Value
	body   js.Value

	// Canvas properties
	canvas js.Value
	ctx    js.Value // Graphics context for the WebGL
	width  int
	height int

	reqID    js.Value // Storage of the current annimationFrame requestID - For Cancel
	timeStep float64  // Min Time delay between frames. - Calculated as   maxFPS/1000

	FPS float64 // populated by the Render loop.   most recent calculated FPS
}

func NewCanvasWebgl(create bool) (*CanvasWebGL, error) {

	var c CanvasWebGL

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
func (c *CanvasWebGL) Create(width int, height int) {

	// Make the Canvas
	canvas := c.doc.Call("createElement", "canvas")

	canvas.Set("height", height)
	canvas.Set("width", width)
	c.body.Call("appendChild", canvas)

	c.Set(canvas, width, height)
}

// Used to setup with an existing Canvas element which was obtained from JS
func (c *CanvasWebGL) Set(canvas js.Value, width int, height int) {
	c.canvas = canvas
	c.height = height
	c.width = width

	// Setup the WebGL Drawing context
	c.ctx = c.canvas.Call("getContext", "webgl")

}

// Starts the annimationFrame callbacks running.   (Recently seperated from Create / Set to give better control for when things start / stop)
func (c *CanvasWebGL) Start(rf WebGLRenderFunc) {
	c.initFrameUpdate(rf)
}

func (c *CanvasWebGL) Height() int {
	return c.height
}

func (c *CanvasWebGL) Width() int {
	return c.width
}

func (c *CanvasWebGL) Document() js.Value {
	return c.doc
}

func (c *CanvasWebGL) GraphicContext() js.Value {
	return c.ctx
}

func (c *CanvasWebGL) Window() js.Value {
	return c.window
}

func (c *CanvasWebGL) Canvas() js.Value {
	return c.canvas
}

// handles calls from requestAnimationFrame
func (c *CanvasWebGL) initFrameUpdate(rf WebGLRenderFunc) {
	// Hold the callbacks without blocking

	go func() {
		var renderFrame js.Func
		var lastTS float64

		renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			tsnow := c.window.Get("performance").Call("now").Float()
			delta := tsnow - lastTS
			fps := 1 / delta * 1000
			lastTS = tsnow
			c.FPS = fps

			rf(c, c.ctx)

			c.reqID = js.Global().Call("requestAnimationFrame", renderFrame) // Captures the requestID to be used in Close / Cancel
			return nil
		})

		defer renderFrame.Release()
		js.Global().Call("requestAnimationFrame", renderFrame)
		<-c.done
	}()
}

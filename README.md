# go-canvas

go-canvas is a pure **go**+**webassembly** Library for efficiently drawing on a html5 `canvas` element within the browser from go, without requiring calls back to JS to utilise canvas drawing functions.  

The library provides the following features:
- Abstracts away the initial DOM interactions to setup the canvas. 
- Creates the shadow image frame, and graphical Context to draw on it. 
- Initializes basic font cache for text using truetype font. 
- Sets up and handles `requestAnimationFrame` callback from the browser. 

## Concept 
go-canvas takes an alternate approach to the current common methods for using canvas, allowing all drawing primitives to be done totaly with go code, without calling JS. 

### standard syscall way
In a standard WASM application for canvas, the go code must create a function that responds to `requestAnimationFrame` callbacks, and has to do complete rendering within that call.  Also It to interacts with the canvas drawing primitives via the syscall/js functions and context switches,  i.e. 
```go
laserCtx.Call("beginPath")
laserCtx.Call("arc", gs.laserX, gs.laserY, gs.laserSize, 0, math.Pi*2, false)
laserCtx.Call("fill")
laserCtx.Call("closePath")
```

Apart from messy JS calls, which couldn't easily be checked at compile time, one other downside of this I didn't like, is it forces a full redraw every frame, even if nothing changed on that canvas.  

### go native way
go-canvas seperates the drawing, from the `requestAnimationFrame`, and does all drawing with go.  It does this by creating an entirley seperate image buffer, which is drawn to using a 2D drawing library.  I'm currently using the one from  https://github.com/llgcode/draw2d which provides most of the standard canvas primites, and more.    This shadow Image buffer can be updated at whatever rate the developer deems appropriate, which may very well be slower than the browsers annimation rate. 

This shadow Image buffer is then copied over to the browser canvas buffer, each `requestAnimationFrame`, at whatever rate the browser requests.  The handling of the callback, and copy is done automatically within the library. 

Drawing therefore, is pure **go**  i.e. 

```go
gc := canvas.Gc()  // Grab the graphic context for drawing to shadow image frame
// draws red 🔴 laser
gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})

gc.BeginPath()
gc.ArcTo(gs.laserX, gs.laserY, gs.laserSize, gs.laserSize, 0, math.Pi*2)
gc.FillStroke()
gc.Close()
```
A simple way to cause the code to draw the frame on schedule, independant from the browsers callbacks, is to use `time.Tick`.  An example is in the demo app below. 

If however your image is only updated from either user input, or some network activity, then it would be straightforward to fire the redraw only when required from these inputs.  For all other cycles of the `requestAnimationFrame` it just copies the buffer over, and nothing changes. 

### Known issues !
There is currently a likley race condition for long draw functions, where the `requestAnimationFrame` may get a partially completed image buffer.  This is more likley the longer the user render operation takes.    Currently think how best to handle this, ideally without locks. 


# Demo
A simple demo can be found in  ./demo directory.  
This is a shamless rewrite of the 'Moving red Laser' demo by Martin Olsansky https://medium.freecodecamp.org/webassembly-with-golang-is-fun-b243c0e34f02


Compile with  `GOOS=js GOARCH=wasm go build -o main.wasm`

Includes a Caddy configuration file to support WASM,  so will serve by just running 'caddy' in the demo directory and opening browser to http://localhost:8080

## Live
Live Demo available at : https://markfarnan.github.io/go-canvas


# Future
This library was written after a weekend of investigation and posted on request for the folks on #webassembly on Gophers Slack.  Right now it is very v0.001, user beware !

I intend to extend it further, time permitting, into fully fledged support package for all things go-canvas-wasm related, using this image frame method.   

Several of the ideas i'm considering are: 
- [ ] Support for layered canvas, at least 3 for 'background', 'action'  and 'user interaction'
- [ ] Traps & helper functions for mouse interactions over the canvas
- [ ] Unit tests - soon as I figure out how to do tests for WASM work. 
- [ ] Performance improvments in the image buffer copy - https://github.com/agnivade/shimmer/blob/c073303a81ab9a90b6fc14eb6d90c3a1b930025e/load_image_cb.go#L40 has been suggested as a place to start. 
- [ ] Detect if nothing has changed for the frame, and if so, don't even recopy the buffer, saving yet more time.   May be usefull for layers that change less frequently. 
- [ ] Multiple draw / render frames to fix the 'incomplete image' problem. 

Others ? Feedback, suggestions etc welcome.  I can be found on Gophers Slack, #Webassembly channel. 

Mark Farnan, May 2019

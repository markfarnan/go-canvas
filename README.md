# go-canvas

go-canvas is a pure **go**+**webassembly** Library for efficiently drawing on a html5 `canvas` element within the browser from go without requiring calls back to JS to utilise canvas drawing functions.

The library provides the following features:
- Abstracts away the initial DOM interactions to setup the canvas.
- Creates the shadow image frame, and graphical Context to draw on it.
- Initializes basic font cache for text using truetype font.
- Sets up and handles `requestAnimationFrame` callback from the browser.

## Concept 
go-canvas takes an alternate approach to the current common methods for using canvas, allowing all drawing primitives to be done totaly with go code, without calling JS. 

### standard syscall way
In a standard WASM application for canvas, the go code must create a function that responds to `requestAnimationFrame` callbacks and renders the frame within that call. It interacts with the canvas drawing primitives via the syscall/js functions and context switches, i.e. 

```go
laserCtx.Call("beginPath")
laserCtx.Call("arc", gs.laserX, gs.laserY, gs.laserSize, 0, math.Pi*2, false)
laserCtx.Call("fill")
laserCtx.Call("closePath")
```

Downsides of this approach (for me at least), are messy JS calls which can't easily be checked at compile time, forcing a full redraw every frame, even if nothing changed on that canvas, or changes being much slower than the requested frame rate. 

### go native way
go-canvas allows all drawing to be done nativley using Go by creating an entirley seperate image buffer which is drawn to using a 2D drawing library. I'm currently using one from https://github.com/llgcode/draw2d which provides most of the standard canvas primites and more. This shadow Image buffer can be updated at whatever rate the developer deems appropriate, which may very well be slower than the browsers annimation rate. 

This shadow Image buffer is then copied over to the browser canvas buffer during each `requestAnimationFrame` callback, at whatever rate the browser requests. The handling of the callback and copy is done automatically within the library.

Secondly, this also allows the option of drawing to the imageBuffer, outside of the `requestAnimationFrame` callback if required. After some testing it appears that it is still best to do the drawing within the `requestAnimationFrame` callback.

go-canvas provides serveral options to controll all this, and take care of the browser/dom interactions
 - User specifes the go render/draw callback method when calling the START function. This callback passes the graphical context to the render routine.
 - Render routine can choose to return whether any drawing actually took place. If it returns false, then the `requestAnimationFrame` callback does nothing, just returns immediatly, saving CPU cycles. (No point to copy buffers, and redraw if nothing has changed) This allows the drawing to be adaptive to the rate of data changes. 
 - The 'start' function accepts a maxFPS parameter. The library will automatically throttle the `requestAnimationFrame` callback to only do redraws or imagebuffer copies to the this max rate. Note it MAY be slower depending ont he Render time, and the requirments of the browser doing other work. When a tab is hidden, the browser regularly reduces and may even stop call to the animation callback. No critical timing should be done in the render/draw routings. 
 - You may pass 'nil' for the render function. In this case all drawing happens totaly under the users control, outside of the library. This may be more usefull in future when WASM supports proper threading. Right now however, testing shows it is slower as all work is in the one thread, and you loose the scheduling benefits of the `requestAnimationFrame` call. 

Drawing therefore, is pure **go** i.e. 

```go
func Render(gc *draw2dimg.GraphicContext) bool {
    // {some movement code removed for clarity, see the demo code for full function}
    // draws red 🔴 laser
    gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
    gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})

    gc.BeginPath()
    gc.ArcTo(gs.laserX, gs.laserY, gs.laserSize, gs.laserSize, 0, math.Pi*2)
    gc.FillStroke()
    gc.Close()
return true  // Yes, we drew something, copy it over to the browser
```
If you do want to render outside the animation loop, a simple way to cause the code to draw the frame on schedule, independant from the browsers callbacks, is to use `time.Tick`. An example is in the demo app below. 

If however your image is only updated from user input or some network activity, then it would be straightforward to fire the redraw only when required from these inputs. This can be controlled within the Render function, by just returning FALSE at the start. Nothing is draw, nor copied (saving CPU time) and the previous frames data remains.

### Known issues !
~~There is currently a likley race condition for long draw functions, where the `requestAnimationFrame` may get a partially completed image buffer. This is more likley the longer the user render operation takes. Currently think how best to handle this, ideally without locks.~~ Turns out this is not an issue, due to the single threaded nature. Eventually if drawing is in a seperate thread, this will have to be handled. 


# Demo
A simple demo can be found in ./demo directory. 
This is a shamless rewrite of the 'Moving red Laser' demo by Martin Olsansky https://medium.freecodecamp.org/webassembly-with-golang-is-fun-b243c0e34f02


Compile with `GOOS=js GOARCH=wasm go build -o main.wasm`

Includes a Caddy configuration file to support WASM, so will serve by just running 'caddy' in the demo directory and opening browser to http://localhost:8080

## Live
Live Demo available at : https://markfarnan.github.io/go-canvas

# Future
This library was written after a weekend of investigation and posted on request for the folks on #webassembly on Gophers Slack. Right now it is very v0.001, user beware !

I intend to extend it further, time permitting, into fully fledged support package for all things go-canvas-wasm related, using this image frame method. 

Several of the ideas I'm considering are: 
- [ ] Support for layered canvas, at least 3 for 'background', 'action' and 'user interaction'
- [ ] Traps & helper functions for mouse interactions over the canvas
- [ ] Unit tests - soon as I figure out how to do tests for WASM work. 
- [ ] Performance improvements in the image buffer copy - https://github.com/agnivade/shimmer/blob/c073303a81ab9a90b6fc14eb6d90c3a1b930025e/load_image_cb.go#L40 has been suggested as a place to start. 
- [X] Detect if nothing has changed for the frame, and if so, don't even recopy the buffer, saving yet more time. May be useful for layers that change less frequently. 
- [X] Multiple draw / render frames to fix the 'incomplete image' problem. -- Not actually a problem
- [X] Tidy up the close/end frame functionality to properly release resources on page unload, and prevent 'broweser reload errors' due to missing annimation callback function. 
- [ ] Add FPS Calculater metric

Others ? Feedback, suggestions etc welcome. I can be found on Gophers Slack, #Webassembly channel. 

Mark Farnan, May 2019

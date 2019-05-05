# gocanvas
Library to use HTML5 Canvas  from Go-WASM, with all drawing within go code.

This version 0.001 and was part of some experiments over the weekend.   I intend to flesh it out into more fully fledged library for generic Canvas support via Go-WASM

Rather than use context calls to Javascript for the drawing,  this methos does all drawing on a shaddow 2D image buffer, which is then copied into the Canvas buffer on each 'render' callback. 
The library sets up the callback and handles the JS interop.  

This approach has two advantages for my work. 
 - All drawing is done in Go.  No context calls to javascript.   Means I can use other drawing libs, and render elsewhere
 - Seperates the drawing  of the image buffer, from the browsers rendering.   Allows better decision making as to what to draw, when. (i.e. don't redraw if nothing changed !)

2D drawing is done using the GO 2D library from  https://github.com/llgcode/draw2d

Simple demo can be found in  ./demo directory 
This is a shamless rewrite of the 'Moving red Laser' demo by Martin Olsansky https://medium.freecodecamp.org/webassembly-with-golang-is-fun-b243c0e34f02

Compile with  GOOS=js GOARCH=wasm go build -o main.wasm

Includes a Caddy configuration file to support WASM,  so will serve by just running 'caddy' in the demo directory and opening browser to http://localhost:8080


Live Demo available at : https://markfarnan.github.io/go-canvas

Feedback, suggestions etc welcome.  I can be found on Gophers Slack, #Webassembly channel. 

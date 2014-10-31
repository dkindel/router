
package main

import (
	"fmt"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glu"
	"github.com/go-gl/gl"
	"time"
	"runtime"
	"sync/atomic"
)

var windows int32

// Make sure to not just drop the event channel if you don't want it, doing so will
// block the ui thread when that chan fills up. Instead start a new goroutine to ignore
// all values.
func NewWindow(w, h int) *Window {
	if atomic.AddInt32(&windows, 1) != 1 {
		panic("TODO: Look into how to do multiple windows")
	}
	
	ret := &Window{make(chan func(), 10), make(chan MouseEvent, 2), w, h}
	ch := make(chan struct{})
	go ret.guiLoop(ch)
	<-ch
	return ret
}

type Window struct {
	// Funcs that are received on this channel are run in the gl thread, allowing you to do drawing.
	// Do not receive from this channel.
	Paint chan func()
	// Events are sent to this channel. Make sure to not let this channel get full or else the ui thread
	// will block. If you want to ignore all events call the IgnoreEvents method.
	// Sending to this channel will simulate events.
	Events chan MouseEvent
	Width, Height int
}

// Call in order to ignore all events without blocking.
// There is no way to undo this operation.
func (w*Window) IgnoreEvents() {
	go func() {
		for _ = range w.Events {}
	}()
}

func SetColor(r, g, b float64) {
	gl.Color3d(r, g, b)
}

func Square(x, y, w, h int) {
	gl.Begin(gl.QUADS)
		gl.Vertex2i(x, y)
		gl.Vertex2i(x, y+h)
		gl.Vertex2i(x+w, y+h)
		gl.Vertex2i(x+w, y)
	gl.End()
}

type MouseEvent struct {
	X, Y float64
	Button MouseButton
	Down bool
}

type Mouser interface {
	MouseSend(MouseEvent)
}

type MouseButton int
const (
	LeftButton MouseButton = 1<<iota
	RightButton
)



func initGL(w, h int) {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	if err := glfw.OpenWindow(w, h, 8,8,8,8,0,0,glfw.Windowed); err != nil {
		panic(err)
	}

	errno := gl.Init()
	if errno != gl.NO_ERROR {
		str, err := glu.ErrorString(errno)
		if err != nil {
			panic(fmt.Sprint(err, errno))
		}
		panic(str)
	}

	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.LIGHTING)
	gl.ClearColor(0.2, 0.2, 0.23, 0.0)
}

func (w *Window) guiLoop(initted chan struct{}) {
	runtime.LockOSThread()

	initGL(w.Width, w.Height)
	close(initted)

	defer glfw.Terminate()
	defer glfw.CloseWindow()

	glfw.SetWindowTitle("Banana")
	glfw.SetSwapInterval(10)
	glfw.SetWindowSizeCallback(onResize)

	var event MouseEvent
	glfw.SetMouseButtonCallback(func(btn, state int) {
		// glfw is single threaded, so it's okay to alter event
		//  in this fashion
		switch btn {
			case 0:
				event.Button = LeftButton
			case 1:
				event.Button = RightButton
		}
		event.Down = state != 0
		//fmt.Println("mousebutton:", btn, state, event)
		w.Events<-event
	})
	glfw.SetMousePosCallback(func(x, y int) {
		event.X, event.Y = float64(x), float64(y)
		event.Button = 0
		w.Events<-event
	})

	tck := time.NewTicker(time.Second/30)
	_ = tck
mainLoop:
	for {
		good := false
		var d func()
pump:	for {
			select {
				case c, ok := <-w.Paint:
					if !ok {
						break mainLoop
					}
					good, d = true, c
					break pump
				case <-tck.C:
					break pump
				//default:
				//	break pump
			}
		}
		if good || resized {
			resized = false
			gl.Clear(gl.COLOR_BUFFER_BIT)
			if d != nil {
				d()
			}
		}
		glfw.SwapBuffers()
		glfw.Sleep(.006)
	}
}

// Only accessed on gl thread
var resized = false

func onResize(w, h int) {
	resized = true
	gl.DrawBuffer(gl.FRONT_AND_BACK)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, w, h)
	gl.Ortho(0, float64(w), float64(h), 0, -1, 1)
	gl.ClearColor(1,1,0,1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}


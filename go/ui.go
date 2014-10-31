
package main

import (
	"fmt"
	"github.com/go-gl/glfw"
	"github.com/go-gl/glu"
	"github.com/go-gl/gl"
	"time"
	"runtime"
)


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

// Make sure to not just drop the event channel if you don't want it, doing so will
// block the ui thread when that chan fills up. Instead start a new goroutine to ignore
// all values.
func StartGui() (chan<-func(),<-chan MouseEvent) {
	ch := make(chan struct{})
	ret := make(chan func(), 10)
	ret2 := make(chan MouseEvent)
	go guiLoop(ch, ret, ret2)
	<-ch
	return ret, ret2
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

func guiLoop(initted chan struct{}, refresh <-chan func(), evs chan<-MouseEvent) {
	runtime.LockOSThread()

	initGL(500, 500)
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
		fmt.Println("mousebutton:", btn, state, event)
		evs<-event
	})
	glfw.SetMousePosCallback(func(x, y int) {
		event.X, event.Y = float64(x), float64(y)
		event.Button = 0
		evs<-event
	})

	tck := time.NewTicker(time.Second/30)
	_ = tck
mainLoop:
	for {
		good := false
		var d func()
pump:	for {
			select {
				case c, ok := <-refresh:
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
		glfw.Sleep(.016)
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


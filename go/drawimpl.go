package main

import (
	"github.com/go-gl/gl"
	"fmt"
)


func (a actionCanvas) do() {
	gl.PushMatrix()
	gl.Translated(a.p.X + a.c.minOffset.Width, a.p.Y + a.c.minOffset.Height, 0)
	a.c.do()
	gl.PopMatrix()
}

func (a actionPoints) do() {
	var mode gl.GLenum
	switch a.mode {
		case triangles:
			mode = gl.TRIANGLES
		case polygon:
			mode = gl.POLYGON
		//case linestrip: <- Done in actionLineStrip
		//	mode = gl.LINE_STRIP
		default:
			panic(fmt.Sprint("Unknown drawMode:", a.mode))
	}
	gl.Begin(mode)

	for _,p := range a.ps {
		gl.Vertex2d(p.X, p.Y)
	}

	gl.End()
}

func (a actionLineStrip) do() {
	gl.LineWidth(float32(a.w))
	gl.Begin(gl.LINE_STRIP)

	for _,p := range a.ps {
		gl.Vertex2d(p.X, p.Y)
	}

	gl.End()
}


func (a actionColor) do() {
	gl.Color3d(a.c.R, a.c.G, a.c.B)
}

func (c*CanvasMaker) PushTranslate(x, y float64) {
	c.Func(func() {
			gl.PushMatrix()
			gl.Translated(x, y, 0)
		})
}
func (c*CanvasMaker) Translate(x, y float64) {
	c.Func(func() {
			gl.Translated(x, y, 0)
		})
}
func (c*CanvasMaker) PushScale(x, y float64) {
	c.Func(func() {
			gl.PushMatrix()
			gl.Scaled(x, y, 1)
		})
}
func (c*CanvasMaker) Scale(x, y float64) {
	c.Func(func() {
			gl.Scaled(x, y, 1)
		})
}

func (c*CanvasMaker) Pop() {
	c.Func(gl.PopMatrix)
}





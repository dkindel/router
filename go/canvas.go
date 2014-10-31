package main

import (
	"fmt"
	"sync/atomic"
)

type Dimensions struct {
	Width, Height float64
}

type Color struct {
	R, G, B float64
}

type TexCoord struct {
	X, Y float64
}

type Point struct {
	X, Y float64
}
func (p Point) Delta(dx, dy float64) Point {
	return Point{p.X + dx, p.Y + dy}
}

// Note: Due to the way I do this, painting could generate a bunch of garbage. This may
//    be negligable, but for painting that changes often it could add up.
// If so, there are a few things I suppose I could do, like having a "union"-esc thing since
//    most of the space would be floats.
type action interface {
	do()
}

type CanvasMaker struct {
	actions []action
	// Will be modified as this object gets used
	max, min Dimensions
	router MouseRouter
}

func (c*CanvasMaker) updateBoundsBox(b Box) {
	// TODO: Due to the ordering of points, this could be optimized
	c.updateBoundsPoint(b.TL)
	c.updateBoundsPoint(b.TR)
	c.updateBoundsPoint(b.BL)
	c.updateBoundsPoint(b.BR)
}

func (c*CanvasMaker) updateBoundsPoints(ps []Point) {
	for _,x := range ps {
		c.updateBoundsPoint(x)
	}
}
func (c*CanvasMaker) updateBoundsPoint(p Point) {
	if p.X < c.min.Width {
		c.min.Width = p.X
	} else if p.X > c.max.Width {
		c.max.Width = p.X
	}

	if p.Y < c.min.Height {
		c.min.Height = p.Y
	} else if p.Y > c.max.Height {
		c.max.Height = p.Y
	}
}

// An immutable drawing, made by calling CanvasMaker.Freeze
//   Note: there's not much of a difference right now between a
//   Canvas and a CanvasMaker, but they exist just in case I find
//   an optimization that requires some sort of freezing.
type Canvas struct {
	// Embedded so that users can get Width/Height easily
	Dimensions
	acts []action
	// The minimum of the canvas, so that we can translate to the right/down this
	//    so that when the user wrote to a negative point it doesn't come out to
	//    the left of where the drawer expects it to be.
	// To prevent an extra push/pop matrix, this is done in actionCanvas.do()
	minOffset Dimensions
	Router MouseRouter
	id uint32
}

func (c Canvas) Equals(o interface{}) bool {
	//fmt.Printf("Canvas.equals:\n\t%v\n\t%v\n", c, o)
	if o, ok := o.(Canvas); !ok {
		return false
	} else {
		return c.id == o.id
	}
}

func (c Canvas) do() {
	for _,x := range c.acts {
		x.do()
	}
}

type Box struct {
	TL, TR, BL, BR Point
}
func (b *Box) AddMargin(amt float64) {
	b.TL.X -= amt
	b.TL.Y -= amt
	b.TR.X += amt
	b.TR.Y -= amt
	b.BL.X -= amt
	b.BL.Y += amt
	b.BR.X += amt
	b.BR.Y += amt
}

func BoxFromDims(x, y, width, height float64) Box {
	ox, oy := x + width, y + height
	return Box{
		Point{ x,  y},
		Point{ox,  y},
		Point{ x, oy},
		Point{ox, oy},
	}
}

// Returns true if p is inside (including bounds) the Box
// Note: Pointer receiver to save on stack space; doesn't actually
//    alter Box
func (b *Box) Inside(p Point) bool {
	return b.TL.X <= p.X && b.TL.Y <= p.Y &&
	       b.TR.X >= p.X && b.TR.Y <= p.Y &&
	       b.BL.X <= p.X && b.BL.Y >= p.Y &&
	       b.BR.X >= p.X && b.BR.Y >= p.Y
}

type VarCanvas struct {
	f func(*CanvasMaker, Dimensions)
}

func (v VarCanvas) Draw(width, height float64) Canvas {
	ret := NewCanvasMaker()
	v.f(ret, Dimensions{width, height})
	return ret.Freeze()
}

func DrawWithDims(f func(*CanvasMaker, Dimensions)) VarCanvas {
	return VarCanvas{f}
}

func NewCanvasMaker() *CanvasMaker {
	return &CanvasMaker{router:MouseRouter{}}
}

var prevCanvasId uint32
func newCanvasId() uint32 {
	return atomic.AddUint32(&prevCanvasId, 1)
}

func (m CanvasMaker) Freeze() Canvas {
	mW, mH := -m.min.Width, -m.min.Height
	return Canvas{Dimensions{m.max.Width + mW, m.max.Height + mH}, m.actions, Dimensions{mW, mH},m.router, newCanvasId()}
}

type drawMode uint8
const (
	triangles drawMode = iota
	polygon

	//linestrip <- Not needed, has a special struct
)

// ACTIONS - The implementations (the do() funcs) will be in different files so that
//   I can eventually support more than just opengl

type actionFunc func()
func (a actionFunc) do() {
	(func())(a)()
}

type actionCanvas struct {
	p Point
	c Canvas
}

type actionPoints struct {
	ps []Point
	mode drawMode
}

type actionPointCols struct {
	ps []Point
	cs []Color
	mode drawMode
}

type actionPointTexs struct {
	ps []Point
	ts []TexCoord
	mode drawMode
}

type actionColor struct {
	c Color
}

type actionLineStrip struct {
	w float64
	ps []Point
}

// END ACTIONS

func (c*CanvasMaker) ChangeColor(col Color) {
	c.actions = append(c.actions, actionColor{col})
}

func (c*CanvasMaker) Polygon(ps []Point) {
	c.actions = append(c.actions, actionPoints{ps, polygon})
	c.updateBoundsPoints(ps)
}

func (c*CanvasMaker) LineStrip(w float64, ps []Point) {
	c.actions = append(c.actions, actionLineStrip{w, ps})
	c.updateBoundsPoints(ps)
}

func (c*CanvasMaker) Func(f func()) {
	c.actions = append(c.actions, actionFunc(f))
}

func (c*CanvasMaker) Canvas(p Point, d Canvas) Box {
	c.actions = append(c.actions, actionCanvas{p, d})

	ret := BoxFromDims(p.X, p.Y, d.Width, d.Height)
	c.updateBoundsBox(ret)
	return ret
}

func (c*CanvasMaker) CanvasHit(p Point, d Canvas, h Mouser) Box {
	ret := c.Canvas(p, d)
	
	c.router.Add(h, ret)
	
	return ret
}

// TODO: As well as left/top/right/bottom only margins
func (c*CanvasMaker) AddMargin(x, y float64) {
	panic("Unimplemented")
}



type hitbox struct {
	m Mouser
	b Box
}

type MouseRouter struct {
	boxes []hitbox
	// TODO: Be able to zoom and rotate things.
	//  I'll have to look into this
	//zoom float64
	//rotate float64
}

func (r*MouseRouter) Clear() {
	r.boxes = nil
}

func (r*MouseRouter) Add(m Mouser, b Box) {
	r.boxes = append(r.boxes, hitbox{m, b})
}

var debugMouse = false

// Will send the event to multiple Mousers if
//   it sees that it hits multiple boxes
func (r MouseRouter) Route(m MouseEvent) {
	p := m.Point

	sent := false
	if debugMouse { fmt.Println("Mouse Route:", m) }
	for _,x := range r.boxes {
		if x.b.Inside(p) {
			// We have to translate the point so that 
			newEvent, ref := m, x.b.TL
			newEvent.X -= ref.X
			newEvent.Y -= ref.Y
			if debugMouse { fmt.Printf("Mouse Send: %T %v\n", x.m, x.b) }
			x.m.MouseSend(newEvent)
			sent = true
		}
	}
	if !sent && debugMouse {
		fmt.Println("Mouse Fail:")
		for _,v := range r.boxes {
			fmt.Printf("\t|%g %g| %T\n",v.b.TL, v.b.TR, v.m)
			fmt.Printf("\t|%g %g|\n"   ,v.b.BL, v.b.BR)
		}
	}
}









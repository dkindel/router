
package main

import ("fmt")

func main() {
	CanvasDemo()
}

func CanvasDemo() {
	w := NewWindow(500, 500)
	
	mainWindow(w)
}

func ps(xy ...float64) []Point {
	l := len(xy)
	if l % 2 == 1 {
		panic("Must have even number of parameters!")
	}
	ret := make([]Point, l/2)
	
	for i := 0; i < l; i+=2 {
		p := &ret[i/2]
		p.X = xy[i]
		p.Y = xy[i+1]
	}
	
	return ret
}

func mainWindow(w *Window) {
	// Let's make a button and put it on the screen!
	
	btnpaintchan := make(chan VarCanvas)
	btn := NewButton(btnpaintchan)
	
	var router MouseRouter
	for {
		select {
		case x := <-btnpaintchan:
			cm := NewCanvasMaker()
			cm.CanvasHit(Point{10, 10}, x.Draw(50, 50), btn)
			c := cm.Freeze()
			router = c.Router
			w.Paint <- c.do
		case m := <-w.Events:
			router.Route(m)
		}
	}
}

func MissionDemo() {
	// The demo is this:
	//
	// The window that is popped up is simply an interface to receive mouse-clicks
	//
	// Each round consists of:
	//     - the terminal prints out available missions
	//     - Wait for mouse click
	//     - Complete the first mission in the list
	//
	// The mission structure is found in missions.go. It is relatively self documenting, but
	// an ascii drawing wouldn't be bad I guess.
	w := NewWindow(500, 500)
	defer close(w.Paint)
	w.Paint <- func(){}
	ev := w.Events
	
	ss := NewStoryState()
	
	av := ss.GetAvailable()
	fmt.Println(av)
	WaitClick(ev)
	ss.Complete(av[0], 0)
	av = ss.GetAvailable()
	fmt.Println(av)
	WaitClick(ev)
	ss.Complete(av[0], 0)
	av = ss.GetAvailable()
	fmt.Println(av)
	WaitClick(ev)
	ss.Complete(av[0], 10)
	av = ss.GetAvailable()
	fmt.Println(av)
	WaitClick(ev)
	ss.Complete(av[0], 90)
	av = ss.GetAvailable()
	fmt.Println(av)
	WaitClick(ev)
}

func WaitClick(ch <-chan MouseEvent) {
	for x := range ch {
		if x.Button == LeftButton && !x.Down {
			return
		}
	}
}

func ignore(ch <-chan MouseEvent) {
	// I can't wait for 1.4
	for _ = range ch {}
}
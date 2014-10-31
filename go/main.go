
package main

import ("fmt")

func main() {
	MissionDemo()
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
	w.Paint <- func() {
		SetColor(1, 0, 0)
		Square(4, 4, 50, 50)
		SetColor(0, 0, 1)
		Square(45, 45, 100, 50)
	}
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
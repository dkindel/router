package main


// Simple-stupid hand-rolled types for now

type Button struct {
	msgs chan MouseEvent
	paint chan<-VarCanvas
	
	// variables local to the instance
	down bool
}

func NewButton(paint chan<-VarCanvas) *Button {
	ret := &Button{
		make(chan MouseEvent, 1),
		paint,
		false,
	}
	go ret.loop()
	return ret
}

func (b *Button) MouseSend(m MouseEvent) {
	b.msgs <- m
}

func (b*Button) loop() {
	b.paint <- b.repaint()
	for x := range b.msgs {
		redraw := false
		if x.Button == LeftButton {
			b.down = x.Down
			redraw = true
		}
		
		if redraw {
			b.paint <- b.repaint()
		}
	}
}

func (b*Button) repaint() VarCanvas {
	// Capture relevant variables
	down := b.down
	return DrawWithDims(func(cm *CanvasMaker, dim Dimensions) {
		var col Color
		if down {
			col = Color{.5,.5,.5}
		} else {
			col = Color{.2,.2,.2}
		}
		cm.ChangeColor(col)
		cm.Polygon(ps(0, 0, dim.Width, 0, dim.Width, dim.Height, 0, dim.Height))
	})
}








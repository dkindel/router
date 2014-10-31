
package main

func main() {
	ui, ev := StartGui()
	go ignore(ev)
	
	ui <- func() {
		SetColor(1, 0, 0)
		Square(4, 4, 50, 50)
		SetColor(0, 0, 1)
		Square(45, 45, 100, 50)
	}
	
	select{}
}

func ignore(ch <-chan MouseEvent) {
	// I can't wait for 1.4
	for _ = range ch {}
}
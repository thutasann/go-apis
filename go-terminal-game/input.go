package main

import "os"

type input struct {
	pressedKey byte
}

func (i *input) update() {
	b := make([]byte, 1)
	os.Stdin.Read(b)
	i.pressedKey = b[0]
}

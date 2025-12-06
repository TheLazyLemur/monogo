package main

import "github.com/TheLazyLemur/monogo"

type RTSUI struct {
	monogo.BaseComponent
}

func NewRTSUI() *RTSUI {
	return &RTSUI{}
}

func (u *RTSUI) DrawUI() {
	monogo.DrawText("RTS Controls:", 10, 30, 16, monogo.White)
	monogo.DrawText("WASD/Edge Pan - Move Camera", 10, 50, 14, monogo.LightGray)
	monogo.DrawText("Mouse Wheel - Zoom", 10, 66, 14, monogo.LightGray)
	monogo.DrawText("Left Click/Drag - Select Units", 10, 82, 14, monogo.LightGray)
	monogo.DrawText("Right Click - Move Selected", 10, 98, 14, monogo.LightGray)
	monogo.DrawText("Shift+Click - Add to Selection", 10, 114, 14, monogo.LightGray)
	monogo.DrawText("ESC - Deselect All", 10, 130, 14, monogo.LightGray)
}

package main

import "github.com/TheLazyLemur/monogo"

// RTSCamera provides top-down RTS-style camera control
type RTSCamera struct {
	monogo.BaseComponent
	PanSpeed    float32
	ZoomSpeed   float32
	MinZoom     float32
	MaxZoom     float32
	EdgePanSize int // pixels from edge to trigger pan

	camX, camZ float32
	zoom       float32
}

func NewRTSCamera() *RTSCamera {
	return &RTSCamera{
		PanSpeed:    20,
		ZoomSpeed:   5,
		MinZoom:     5,
		MaxZoom:     30,
		EdgePanSize: 20,
		zoom:        15,
	}
}

func (c *RTSCamera) Start() {
	c.updateCamera()
}

func (c *RTSCamera) Update(dt float32) {
	// WASD panning
	c.camX += monogo.Horizontal() * c.PanSpeed * dt
	c.camZ -= monogo.Vertical() * c.PanSpeed * dt

	// Edge panning
	mx, my := monogo.MousePosition()
	sw, sh := monogo.ScreenWidth(), monogo.ScreenHeight()

	if mx < c.EdgePanSize {
		c.camX -= c.PanSpeed * dt
	} else if mx > sw-c.EdgePanSize {
		c.camX += c.PanSpeed * dt
	}
	if my < c.EdgePanSize {
		c.camZ -= c.PanSpeed * dt
	} else if my > sh-c.EdgePanSize {
		c.camZ += c.PanSpeed * dt
	}

	// Zoom with mouse wheel
	wheel := monogo.MouseWheel()
	c.zoom -= wheel * c.ZoomSpeed
	if c.zoom < c.MinZoom {
		c.zoom = c.MinZoom
	}
	if c.zoom > c.MaxZoom {
		c.zoom = c.MaxZoom
	}

	c.updateCamera()
}

func (c *RTSCamera) updateCamera() {
	cam := c.GameObject.Scene().Camera
	// Top-down angled view
	cam.SetPosition(c.camX, c.zoom, c.camZ+c.zoom*0.5)
	cam.SetTarget(c.camX, 0, c.camZ)
}

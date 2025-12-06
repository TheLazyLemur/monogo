package main

import "github.com/TheLazyLemur/monogo"

// SelectionManager handles box selection and click-to-move
type SelectionManager struct {
	monogo.BaseComponent
	selected []*Unit

	// Box selection state
	isDragging bool
	startX     int
	startY     int
	endX       int
	endY       int
}

func NewSelectionManager() *SelectionManager {
	return &SelectionManager{
		selected: make([]*Unit, 0),
	}
}

func (s *SelectionManager) Update(dt float32) {
	mx, my := monogo.MousePosition()

	// Left click - start box select
	if monogo.MousePressed(monogo.MouseLeft) {
		s.isDragging = true
		s.startX, s.startY = mx, my
		s.endX, s.endY = mx, my
	}

	// Dragging - update box
	if s.isDragging {
		s.endX, s.endY = mx, my
	}

	// Left release - finish selection
	if monogo.MouseReleased(monogo.MouseLeft) && s.isDragging {
		s.isDragging = false
		s.finishSelection()
	}

	// Right click - move selected units
	if monogo.MousePressed(monogo.MouseRight) && len(s.selected) > 0 {
		if pos, ok := s.GameObject.Scene().Camera.GetGroundPoint(); ok {
			s.moveSelectedTo(pos)
		}
	}

	// Escape - deselect all
	if monogo.KeyPressed(monogo.KeyEscape) {
		s.deselectAll()
	}
}

func (s *SelectionManager) finishSelection() {
	// Clear previous selection (unless shift held)
	if !monogo.KeyHeld(monogo.KeyShift) {
		s.deselectAll()
	}

	// Normalize box coordinates
	x1, x2 := s.startX, s.endX
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	y1, y2 := s.startY, s.endY
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Check if it's a click (small box) or drag
	isClick := (x2-x1) < 5 && (y2-y1) < 5

	// Find units in selection box
	scene := s.GameObject.Scene()
	cam := scene.Camera

	units := scene.FindGameObjectsWithTag("Unit")
	for _, obj := range units {
		unit := monogo.GetComponent[*Unit](obj)
		if unit == nil {
			continue
		}

		// Project unit position to screen
		sx, sy := cam.WorldToScreen(obj.Transform.Position)

		if isClick {
			// Click selection - check if clicked near unit
			dx, dy := sx-s.startX, sy-s.startY
			if dx*dx+dy*dy < 400 { // within 20 pixels
				s.selectUnit(unit)
				break // only select one on click
			}
		} else {
			// Box selection
			if sx >= x1 && sx <= x2 && sy >= y1 && sy <= y2 {
				s.selectUnit(unit)
			}
		}
	}
}

func (s *SelectionManager) selectUnit(unit *Unit) {
	unit.Selected = true
	// Check if already selected
	for _, u := range s.selected {
		if u == unit {
			return
		}
	}
	s.selected = append(s.selected, unit)
}

func (s *SelectionManager) deselectAll() {
	for _, u := range s.selected {
		u.Selected = false
	}
	s.selected = s.selected[:0]
}

func (s *SelectionManager) moveSelectedTo(pos monogo.Vector3) {
	// Simple formation: spread units around target
	count := len(s.selected)
	for i, unit := range s.selected {
		offset := monogo.Vec3(0, 0, 0)
		if count > 1 {
			// Spread in a circle
			angle := float32(i) / float32(count) * 6.28
			radius := float32(count) * 0.3
			offset = monogo.Vec3(
				monogo.Cos(angle*57.3)*radius,
				0,
				monogo.Sin(angle*57.3)*radius,
			)
		}
		unit.MoveTo(monogo.Vec3(
			pos.X+offset.X,
			pos.Y,
			pos.Z+offset.Z,
		))
	}
}

func (s *SelectionManager) DrawUI() {
	if !s.isDragging {
		return
	}

	// Draw selection box
	x1, x2 := s.startX, s.endX
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	y1, y2 := s.startY, s.endY
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Semi-transparent fill
	monogo.DrawRect(x1, y1, x2-x1, y2-y1, monogo.RGBA(0, 255, 0, 50))
	// Outline
	monogo.DrawRectLines(x1, y1, x2-x1, y2-y1, monogo.Green)
}

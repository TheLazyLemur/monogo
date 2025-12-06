package main

import "github.com/TheLazyLemur/monogo"

// Unit is a selectable, moveable RTS unit
type Unit struct {
	monogo.BaseComponent
	Selected  bool
	MoveSpeed float32

	targetPos   monogo.Vector3
	hasTarget   bool
	selectionID int // for selection manager
}

var unitIDCounter int

func NewUnit() *Unit {
	unitIDCounter++
	return &Unit{
		MoveSpeed:   5,
		selectionID: unitIDCounter,
	}
}

func (u *Unit) Update(dt float32) {
	if !u.hasTarget {
		return
	}

	pos := u.GameObject.Transform.Position
	dir := monogo.Vec3(
		u.targetPos.X-pos.X,
		0,
		u.targetPos.Z-pos.Z,
	)

	dist := monogo.Distance(pos, u.targetPos)
	if dist < 0.2 {
		u.hasTarget = false
		return
	}

	// Normalize and move
	dir = monogo.Normalize(dir)
	u.GameObject.Transform.Position.X += dir.X * u.MoveSpeed * dt
	u.GameObject.Transform.Position.Z += dir.Z * u.MoveSpeed * dt
}

func (u *Unit) MoveTo(pos monogo.Vector3) {
	u.targetPos = pos
	u.hasTarget = true
}

func (u *Unit) Stop() {
	u.hasTarget = false
}

func (u *Unit) Draw() {
	// Selection indicator (ring under unit)
	if u.Selected {
		pos := u.GameObject.Transform.Position
		monogo.DrawCircle3D(pos.X, 0.05, pos.Z, 0.7, monogo.Green)
	}
}

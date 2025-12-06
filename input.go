package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// Key constants
type Key = int32

const (
	KeyW      Key = rl.KeyW
	KeyA      Key = rl.KeyA
	KeyS      Key = rl.KeyS
	KeyD      Key = rl.KeyD
	KeyQ      Key = rl.KeyQ
	KeyE      Key = rl.KeyE
	KeyR      Key = rl.KeyR
	KeyF      Key = rl.KeyF
	KeySpace  Key = rl.KeySpace
	KeyShift  Key = rl.KeyLeftShift
	KeyCtrl   Key = rl.KeyLeftControl
	KeyAlt    Key = rl.KeyLeftAlt
	KeyTab    Key = rl.KeyTab
	KeyEscape Key = rl.KeyEscape
	KeyEnter  Key = rl.KeyEnter
	KeyUp     Key = rl.KeyUp
	KeyDown   Key = rl.KeyDown
	KeyLeft   Key = rl.KeyLeft
	KeyRight  Key = rl.KeyRight
	Key1      Key = rl.KeyOne
	Key2      Key = rl.KeyTwo
	Key3      Key = rl.KeyThree
	Key4      Key = rl.KeyFour
	Key5      Key = rl.KeyFive
)

// Mouse button constants
type MouseButton = rl.MouseButton

const (
	MouseLeft   MouseButton = rl.MouseLeftButton
	MouseRight  MouseButton = rl.MouseRightButton
	MouseMiddle MouseButton = rl.MouseMiddleButton
)

// KeyPressed returns true if key was just pressed this frame
func KeyPressed(key Key) bool {
	return rl.IsKeyPressed(key)
}

// KeyReleased returns true if key was just released this frame
func KeyReleased(key Key) bool {
	return rl.IsKeyReleased(key)
}

// KeyHeld returns true if key is currently held down
func KeyHeld(key Key) bool {
	return rl.IsKeyDown(key)
}

// MousePressed returns true if mouse button was just pressed
func MousePressed(button MouseButton) bool {
	return rl.IsMouseButtonPressed(button)
}

// MouseReleased returns true if mouse button was just released
func MouseReleased(button MouseButton) bool {
	return rl.IsMouseButtonReleased(button)
}

// MouseHeld returns true if mouse button is held down
func MouseHeld(button MouseButton) bool {
	return rl.IsMouseButtonDown(button)
}

// MousePosition returns current mouse position
func MousePosition() (x, y int) {
	pos := rl.GetMousePosition()
	return int(pos.X), int(pos.Y)
}

// MouseDelta returns mouse movement since last frame
func MouseDelta() (x, y float32) {
	delta := rl.GetMouseDelta()
	return delta.X, delta.Y
}

// MouseWheel returns mouse wheel movement
func MouseWheel() float32 {
	return rl.GetMouseWheelMove()
}

// Axis returns -1, 0, or 1 based on two keys (negative, positive)
func Axis(negative, positive Key) float32 {
	var value float32
	if KeyHeld(negative) {
		value -= 1
	}
	if KeyHeld(positive) {
		value += 1
	}
	return value
}

// Horizontal returns horizontal axis (A/D or Left/Right)
func Horizontal() float32 {
	h := Axis(KeyA, KeyD)
	if h == 0 {
		h = Axis(KeyLeft, KeyRight)
	}
	return h
}

// Vertical returns vertical axis (W/S or Up/Down)
func Vertical() float32 {
	v := Axis(KeyS, KeyW)
	if v == 0 {
		v = Axis(KeyDown, KeyUp)
	}
	return v
}

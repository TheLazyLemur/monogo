package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// Type aliases for clean API
type (
	Vector3 = rl.Vector3
	Vector2 = rl.Vector2
	Color   = rl.Color
	Model   = rl.Model
)

// Vec3 creates a 3D vector
func Vec3(x, y, z float32) Vector3 {
	return rl.NewVector3(x, y, z)
}

// Vec2 creates a 2D vector
func Vec2(x, y float32) Vector2 {
	return rl.NewVector2(x, y)
}

// RGB creates a color from RGB values (0-255)
func RGB(r, g, b uint8) Color {
	return rl.NewColor(r, g, b, 255)
}

// RGBA creates a color from RGBA values (0-255)
func RGBA(r, g, b, a uint8) Color {
	return rl.NewColor(r, g, b, a)
}

// Common colors
var (
	White       = rl.White
	Black       = rl.Black
	Gray        = rl.Gray
	LightGray   = rl.LightGray
	DarkGray    = rl.DarkGray
	Red         = rl.Red
	Green       = rl.Green
	Blue        = rl.Blue
	Yellow      = rl.Yellow
	Orange      = rl.Orange
	Purple      = rl.Purple
	Pink        = rl.Pink
	Brown       = rl.Brown
	Beige       = rl.Beige
	SkyBlue     = rl.SkyBlue
	Lime        = rl.Lime
	Gold        = rl.Gold
	Maroon      = rl.Maroon
	Magenta     = rl.Magenta
	Violet      = rl.Violet
	RayWhite    = rl.RayWhite
	Transparent = rl.Blank
)

// Delta returns frame time in seconds
func Delta() float32 {
	return rl.GetFrameTime()
}

// Time returns elapsed time since start
func Time() float64 {
	return rl.GetTime()
}

// RandomFloat returns a random float between min and max
func RandomFloat(min, max float32) float32 {
	return min + float32(rl.GetRandomValue(0, 10000))/10000.0*(max-min)
}

// RandomInt returns a random int between min and max (inclusive)
func RandomInt(min, max int) int {
	return int(rl.GetRandomValue(int32(min), int32(max)))
}

// Lerp linearly interpolates between a and b
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

// LerpVec3 linearly interpolates between two vectors
func LerpVec3(a, b Vector3, t float32) Vector3 {
	return Vec3(
		Lerp(a.X, b.X, t),
		Lerp(a.Y, b.Y, t),
		Lerp(a.Z, b.Z, t),
	)
}

// Clamp constrains a value between min and max
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Sin returns sine of angle in radians
func Sin(x float32) float32 {
	v := rl.Vector2Rotate(rl.NewVector2(1, 0), x)
	return v.Y
}

// Cos returns cosine of angle in radians
func Cos(x float32) float32 {
	v := rl.Vector2Rotate(rl.NewVector2(1, 0), x)
	return v.X
}

// Distance returns distance between two 3D points
func Distance(a, b Vector3) float32 {
	return rl.Vector3Distance(a, b)
}

// Normalize returns a normalized vector
func Normalize(v Vector3) Vector3 {
	return rl.Vector3Normalize(v)
}

// LoadModel loads a 3D model from file
func LoadModel(path string) Model {
	return rl.LoadModel(path)
}

// UI Drawing helpers

// DrawRect draws a filled rectangle
func DrawRect(x, y, width, height int, color Color) {
	rl.DrawRectangle(int32(x), int32(y), int32(width), int32(height), color)
}

// DrawRectLines draws a rectangle outline
func DrawRectLines(x, y, width, height int, color Color) {
	rl.DrawRectangleLines(int32(x), int32(y), int32(width), int32(height), color)
}

// DrawCircle3D draws a circle in 3D space (on XZ plane)
func DrawCircle3D(x, y, z, radius float32, color Color) {
	rl.DrawCircle3D(Vec3(x, y, z), radius, Vec3(1, 0, 0), 90, color)
}

// DrawLine3D draws a line in 3D space
func DrawLine3D(x1, y1, z1, x2, y2, z2 float32, color Color) {
	rl.DrawLine3D(Vec3(x1, y1, z1), Vec3(x2, y2, z2), color)
}

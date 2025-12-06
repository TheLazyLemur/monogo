package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// Transform holds position, rotation, and scale for a GameObject
// Position, Rotation, Scale are local (relative to parent)
type Transform struct {
	Position   Vector3 // Local position
	Rotation   Vector3 // Local rotation (euler angles in degrees)
	Scale      Vector3 // Local scale
	gameObject *GameObject
}

// NewTransform creates a transform with default values
func NewTransform() *Transform {
	return &Transform{
		Position: rl.NewVector3(0, 0, 0),
		Rotation: rl.NewVector3(0, 0, 0),
		Scale:    rl.NewVector3(1, 1, 1),
	}
}

// WorldPosition returns the absolute world position
func (t *Transform) WorldPosition() Vector3 {
	if t.gameObject == nil || t.gameObject.parent == nil {
		return t.Position
	}

	parentWorld := t.gameObject.parent.Transform.WorldPosition()
	parentScale := t.gameObject.parent.Transform.WorldScale()
	parentRot := t.gameObject.parent.Transform.WorldRotation()

	// Scale the local position
	scaledPos := Vec3(
		t.Position.X*parentScale.X,
		t.Position.Y*parentScale.Y,
		t.Position.Z*parentScale.Z,
	)

	// Rotate the offset by parent's Y rotation
	// Match raylib's MatrixRotateY convention
	yaw := parentRot.Y * rl.Deg2rad
	rotatedX := scaledPos.X*cosf(yaw) + scaledPos.Z*sinf(yaw)
	rotatedZ := -scaledPos.X*sinf(yaw) + scaledPos.Z*cosf(yaw)

	return Vec3(
		parentWorld.X+rotatedX,
		parentWorld.Y+scaledPos.Y,
		parentWorld.Z+rotatedZ,
	)
}

// WorldRotation returns the absolute world rotation
func (t *Transform) WorldRotation() Vector3 {
	if t.gameObject == nil || t.gameObject.parent == nil {
		return t.Rotation
	}

	parentRot := t.gameObject.parent.Transform.WorldRotation()
	return Vec3(
		parentRot.X+t.Rotation.X,
		parentRot.Y+t.Rotation.Y,
		parentRot.Z+t.Rotation.Z,
	)
}

// WorldScale returns the absolute world scale
func (t *Transform) WorldScale() Vector3 {
	if t.gameObject == nil || t.gameObject.parent == nil {
		return t.Scale
	}

	parentScale := t.gameObject.parent.Transform.WorldScale()
	return Vec3(
		parentScale.X*t.Scale.X,
		parentScale.Y*t.Scale.Y,
		parentScale.Z*t.Scale.Z,
	)
}

// SetWorldPosition sets position in world space
func (t *Transform) SetWorldPosition(pos Vector3) {
	if t.gameObject == nil || t.gameObject.parent == nil {
		t.Position = pos
		return
	}

	parentWorld := t.gameObject.parent.Transform.WorldPosition()
	parentScale := t.gameObject.parent.Transform.WorldScale()

	t.Position = Vec3(
		(pos.X-parentWorld.X)/parentScale.X,
		(pos.Y-parentWorld.Y)/parentScale.Y,
		(pos.Z-parentWorld.Z)/parentScale.Z,
	)
}

// Forward returns the forward direction vector
func (t *Transform) Forward() Vector3 {
	rot := t.WorldRotation()
	// Convert to radians
	yaw := rot.Y * rl.Deg2rad
	pitch := rot.X * rl.Deg2rad

	return Vec3(
		float32(-Sin(yaw)*Cos(pitch)),
		float32(Sin(pitch)),
		float32(-Cos(yaw)*Cos(pitch)),
	)
}

// Right returns the right direction vector
func (t *Transform) Right() Vector3 {
	rot := t.WorldRotation()
	yaw := rot.Y * rl.Deg2rad

	return Vec3(
		float32(Cos(yaw)),
		0,
		float32(-Sin(yaw)),
	)
}

func sinf(x float32) float32 {
	v := rl.Vector2Rotate(rl.NewVector2(1, 0), x)
	return v.Y
}

func cosf(x float32) float32 {
	v := rl.Vector2Rotate(rl.NewVector2(1, 0), x)
	return v.X
}

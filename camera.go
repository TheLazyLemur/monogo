package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// Camera wraps raylib camera with game-friendly features
type Camera struct {
	rl          rl.Camera3D
	target      *GameObject
	offset      rl.Vector3
	followSpeed float32
	mode        CameraMode
}

type CameraMode int

const (
	CameraModeFixed CameraMode = iota
	CameraModeFollow
	CameraModeFollowSmooth
)

// NewCamera creates a default camera
func NewCamera() *Camera {
	return &Camera{
		rl: rl.Camera3D{
			Position:   rl.NewVector3(10, 10, 10),
			Target:     rl.NewVector3(0, 0, 0),
			Up:         rl.NewVector3(0, 1, 0),
			Fovy:       45,
			Projection: rl.CameraPerspective,
		},
		offset:      rl.NewVector3(0, 5, 10),
		followSpeed: 5.0,
		mode:        CameraModeFixed,
	}
}

// Follow makes camera follow a GameObject
func (c *Camera) Follow(go_ *GameObject) {
	c.target = go_
	c.mode = CameraModeFollow
}

// FollowSmooth makes camera smoothly follow a GameObject
func (c *Camera) FollowSmooth(go_ *GameObject, speed float32) {
	c.target = go_
	c.followSpeed = speed
	c.mode = CameraModeFollowSmooth
}

// SetOffset sets the camera offset from target
func (c *Camera) SetOffset(x, y, z float32) {
	c.offset = rl.NewVector3(x, y, z)
}

// SetPosition sets fixed camera position
func (c *Camera) SetPosition(x, y, z float32) {
	c.rl.Position = rl.NewVector3(x, y, z)
}

// SetTarget sets fixed camera target
func (c *Camera) SetTarget(x, y, z float32) {
	c.rl.Target = rl.NewVector3(x, y, z)
}

// SetFOV sets field of view
func (c *Camera) SetFOV(fov float32) {
	c.rl.Fovy = fov
}

// Position returns current camera position
func (c *Camera) Position() rl.Vector3 {
	return c.rl.Position
}

// Target returns current camera target
func (c *Camera) Target() rl.Vector3 {
	return c.rl.Target
}

// LookAt makes camera look at a position
func (c *Camera) LookAt(x, y, z float32) {
	c.rl.Target = rl.NewVector3(x, y, z)
}

// update updates camera (called by engine)
func (c *Camera) update(dt float32) {
	if c.target == nil {
		return
	}

	targetPos := c.target.Transform.Position
	desiredTarget := targetPos
	desiredPosition := rl.NewVector3(
		targetPos.X+c.offset.X,
		targetPos.Y+c.offset.Y,
		targetPos.Z+c.offset.Z,
	)

	switch c.mode {
	case CameraModeFollow:
		c.rl.Target = desiredTarget
		c.rl.Position = desiredPosition
	case CameraModeFollowSmooth:
		t := c.followSpeed * dt
		if t > 1 {
			t = 1
		}
		c.rl.Target = lerpVec3(c.rl.Target, desiredTarget, t)
		c.rl.Position = lerpVec3(c.rl.Position, desiredPosition, t)
	}
}

func (c *Camera) getRl() rl.Camera3D {
	return c.rl
}

func lerpVec3(a, b rl.Vector3, t float32) rl.Vector3 {
	return rl.NewVector3(
		a.X+(b.X-a.X)*t,
		a.Y+(b.Y-a.Y)*t,
		a.Z+(b.Z-a.Z)*t,
	)
}

// GetMouseRay returns a ray from camera through mouse position
func (c *Camera) GetMouseRay() rl.Ray {
	return rl.GetScreenToWorldRay(rl.GetMousePosition(), c.rl)
}

// GetGroundPoint returns where mouse ray hits the ground plane (Y=0)
func (c *Camera) GetGroundPoint() (Vector3, bool) {
	ray := c.GetMouseRay()
	// Ray-plane intersection with Y=0
	if ray.Direction.Y == 0 {
		return Vector3{}, false
	}
	t := -ray.Position.Y / ray.Direction.Y
	if t < 0 {
		return Vector3{}, false
	}
	return Vec3(
		ray.Position.X+ray.Direction.X*t,
		0,
		ray.Position.Z+ray.Direction.Z*t,
	), true
}

// WorldToScreen converts 3D position to screen coordinates
func (c *Camera) WorldToScreen(pos Vector3) (x, y int) {
	screen := rl.GetWorldToScreen(pos, c.rl)
	return int(screen.X), int(screen.Y)
}

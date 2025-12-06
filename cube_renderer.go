package monogo

import rl "github.com/gen2brain/raylib-go/raylib"

// CubeRenderer renders a colored cube at the GameObject's transform
type CubeRenderer struct {
	BaseComponent
	Color         Color
	Size          Vector3
	model         rl.Model
	hasModel      bool
	pendingShader *Shader
}

// NewCubeRenderer creates a cube renderer with color and size
func NewCubeRenderer(color Color, size Vector3) *CubeRenderer {
	return &CubeRenderer{
		Color: color,
		Size:  size,
	}
}

// SetShader applies a custom shader to this renderer
func (c *CubeRenderer) SetShader(s *Shader) *CubeRenderer {
	if c.hasModel {
		mats := c.model.GetMaterials()
		if len(mats) > 0 {
			mats[0].Shader = s.GetRl()
		}
	} else {
		c.pendingShader = s
	}
	return c
}

func (c *CubeRenderer) Start() {
	// Create cube mesh and model
	mesh := rl.GenMeshCube(c.Size.X, c.Size.Y, c.Size.Z)
	c.model = rl.LoadModelFromMesh(mesh)
	c.hasModel = true

	// Apply pending shader if set before Start
	if c.pendingShader != nil {
		mats := c.model.GetMaterials()
		if len(mats) > 0 {
			mats[0].Shader = c.pendingShader.GetRl()
		}
		c.pendingShader = nil
	}
}

func (c *CubeRenderer) Draw() {
	if !c.hasModel {
		return
	}

	t := c.GameObject.Transform
	pos := t.WorldPosition()
	rot := t.WorldRotation()
	worldScale := t.WorldScale()

	// DrawModelEx takes position, rotation axis, rotation angle, scale, tint
	// For full euler rotation, we need to combine rotations
	// Simplified: just use Y rotation for now (most common case)
	rl.DrawModelEx(c.model, pos, Vec3(0, 1, 0), rot.Y, worldScale, c.Color)

	// Draw wireframe
	rl.DrawModelWiresEx(c.model, pos, Vec3(0, 1, 0), rot.Y, worldScale, Black)
}
